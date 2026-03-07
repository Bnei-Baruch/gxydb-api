package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/hlog"

	"github.com/Bnei-Baruch/gxydb-api/pkg/httputil"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				err, ok := p.(error)
				if !ok {
					err = fmt.Errorf("panic: %+v", p)
				}

				hlog.FromRequest(r).Error().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("stack", string(debug.Stack())).
					Msgf("panic recovered: %v", err)

				httputil.NewInternalError(err).Abort(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
