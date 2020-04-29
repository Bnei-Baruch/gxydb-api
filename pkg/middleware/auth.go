package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	pkgerr "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
)

type Roles struct {
	Roles []string `json:"roles"`
}

type IDTokenClaims struct {
	Acr               string           `json:"acr"`
	AllowedOrigins    []string         `json:"allowed-origins"`
	Aud               string           `json:"aud"`
	AuthTime          int              `json:"auth_time"`
	Azp               string           `json:"azp"`
	Email             string           `json:"email"`
	Exp               int              `json:"exp"`
	FamilyName        string           `json:"family_name"`
	GivenName         string           `json:"given_name"`
	Iat               int              `json:"iat"`
	Iss               string           `json:"iss"`
	Jti               string           `json:"jti"`
	Name              string           `json:"name"`
	Nbf               int              `json:"nbf"`
	Nonce             string           `json:"nonce"`
	PreferredUsername string           `json:"preferred_username"`
	RealmAccess       Roles            `json:"realm_access"`
	ResourceAccess    map[string]Roles `json:"resource_access"`
	SessionState      string           `json:"session_state"`
	Sub               string           `json:"sub"`
	Typ               string           `json:"typ"`

	rolesMap map[string]struct{}
}

func (c *IDTokenClaims) HasRole(role string) bool {
	if c.rolesMap == nil {
		c.rolesMap = make(map[string]struct{})
		if c.RealmAccess.Roles != nil {
			for _, r := range c.RealmAccess.Roles {
				c.rolesMap[r] = struct{}{}
			}
		}
	}

	_, ok := c.rolesMap[role]
	return ok
}

func AuthenticationMiddleware(tokenVerifier *oidc.IDTokenVerifier, gwPwd func(string) (string, bool), skipApi, skipGateway bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/event" || r.URL.Path == "/protocol" {
				if skipGateway {
					next.ServeHTTP(w, r)
					return
				} else {
					username, password, ok := r.BasicAuth()
					if !ok {
						httputil.NewUnauthorizedError(pkgerr.Errorf("no `Authorization` header set")).Abort(w, r)
						return
					}

					hPwd, ok := gwPwd(username)
					if !ok {
						httputil.NewUnauthorizedError(pkgerr.Errorf("unknown gateway: %s", username)).Abort(w, r)
						return
					}

					if err := bcrypt.CompareHashAndPassword([]byte(hPwd), []byte(password)); err != nil {
						httputil.NewUnauthorizedError(pkgerr.Errorf("wrong password: %s", password)).Abort(w, r)
						return
					}

					next.ServeHTTP(w, r)
					return
				}
			}

			// APIs are using JWT
			if skipApi {
				next.ServeHTTP(w, r)
				return
			}

			auth := parseToken(r)
			if auth == "" {
				httputil.NewUnauthorizedError(pkgerr.Errorf("no `Authorization` header set")).Abort(w, r)
				return
			}

			token, err := tokenVerifier.Verify(context.TODO(), auth)
			if err != nil {
				httputil.NewUnauthorizedError(err).Abort(w, r)
				return
			}

			var claims IDTokenClaims
			if err := token.Claims(&claims); err != nil {
				httputil.NewBadRequestError(err, "malformed JWT claims").Abort(w, r)
				return
			}

			rCtx, ok := ContextFromRequest(r)
			if ok {
				rCtx.IDClaims = &claims
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseToken(r *http.Request) string {
	authHeader := strings.Split(strings.TrimSpace(r.Header.Get("Authorization")), " ")
	if len(authHeader) == 2 &&
		strings.ToLower(authHeader[0]) == "bearer" &&
		len(authHeader[1]) > 0 {
		return authHeader[1]
	}
	return ""
}

func isAllowedIP(ipAddr string) (bool, error) {
	ip := net.ParseIP(strings.TrimSpace(ipAddr))
	if ip == nil {
		return false, fmt.Errorf("invalid IP address %s", ipAddr)
	}

	_, lcl, _ := net.ParseCIDR("10.66.0.0/16")
	_, vpn, _ := net.ParseCIDR("172.16.102.0/24")
	return lcl.Contains(ip) || vpn.Contains(ip), nil
}
