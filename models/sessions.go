// Code generated by SQLBoiler 3.6.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Session is an object representing the database table.
type Session struct {
	ID             int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID         int64       `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	RoomID         null.Int64  `boil:"room_id" json:"room_id,omitempty" toml:"room_id" yaml:"room_id,omitempty"`
	GatewayID      null.Int64  `boil:"gateway_id" json:"gateway_id,omitempty" toml:"gateway_id" yaml:"gateway_id,omitempty"`
	GatewaySession null.Int64  `boil:"gateway_session" json:"gateway_session,omitempty" toml:"gateway_session" yaml:"gateway_session,omitempty"`
	GatewayHandle  null.Int64  `boil:"gateway_handle" json:"gateway_handle,omitempty" toml:"gateway_handle" yaml:"gateway_handle,omitempty"`
	GatewayFeed    null.Int64  `boil:"gateway_feed" json:"gateway_feed,omitempty" toml:"gateway_feed" yaml:"gateway_feed,omitempty"`
	Display        null.String `boil:"display" json:"display,omitempty" toml:"display" yaml:"display,omitempty"`
	Camera         bool        `boil:"camera" json:"camera" toml:"camera" yaml:"camera"`
	Question       bool        `boil:"question" json:"question" toml:"question" yaml:"question"`
	SelfTest       bool        `boil:"self_test" json:"self_test" toml:"self_test" yaml:"self_test"`
	SoundTest      bool        `boil:"sound_test" json:"sound_test" toml:"sound_test" yaml:"sound_test"`
	UserAgent      null.String `boil:"user_agent" json:"user_agent,omitempty" toml:"user_agent" yaml:"user_agent,omitempty"`
	IPAddress      null.String `boil:"ip_address" json:"ip_address,omitempty" toml:"ip_address" yaml:"ip_address,omitempty"`
	Properties     null.JSON   `boil:"properties" json:"properties,omitempty" toml:"properties" yaml:"properties,omitempty"`
	CreatedAt      time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      null.Time   `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	RemovedAt      null.Time   `boil:"removed_at" json:"removed_at,omitempty" toml:"removed_at" yaml:"removed_at,omitempty"`

	R *sessionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sessionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SessionColumns = struct {
	ID             string
	UserID         string
	RoomID         string
	GatewayID      string
	GatewaySession string
	GatewayHandle  string
	GatewayFeed    string
	Display        string
	Camera         string
	Question       string
	SelfTest       string
	SoundTest      string
	UserAgent      string
	IPAddress      string
	Properties     string
	CreatedAt      string
	UpdatedAt      string
	RemovedAt      string
}{
	ID:             "id",
	UserID:         "user_id",
	RoomID:         "room_id",
	GatewayID:      "gateway_id",
	GatewaySession: "gateway_session",
	GatewayHandle:  "gateway_handle",
	GatewayFeed:    "gateway_feed",
	Display:        "display",
	Camera:         "camera",
	Question:       "question",
	SelfTest:       "self_test",
	SoundTest:      "sound_test",
	UserAgent:      "user_agent",
	IPAddress:      "ip_address",
	Properties:     "properties",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	RemovedAt:      "removed_at",
}

// Generated where

type whereHelpernull_Int64 struct{ field string }

func (w whereHelpernull_Int64) EQ(x null.Int64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int64) NEQ(x null.Int64) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int64) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int64) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_Int64) LT(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int64) LTE(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int64) GT(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int64) GTE(x null.Int64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var SessionWhere = struct {
	ID             whereHelperint64
	UserID         whereHelperint64
	RoomID         whereHelpernull_Int64
	GatewayID      whereHelpernull_Int64
	GatewaySession whereHelpernull_Int64
	GatewayHandle  whereHelpernull_Int64
	GatewayFeed    whereHelpernull_Int64
	Display        whereHelpernull_String
	Camera         whereHelperbool
	Question       whereHelperbool
	SelfTest       whereHelperbool
	SoundTest      whereHelperbool
	UserAgent      whereHelpernull_String
	IPAddress      whereHelpernull_String
	Properties     whereHelpernull_JSON
	CreatedAt      whereHelpertime_Time
	UpdatedAt      whereHelpernull_Time
	RemovedAt      whereHelpernull_Time
}{
	ID:             whereHelperint64{field: "\"sessions\".\"id\""},
	UserID:         whereHelperint64{field: "\"sessions\".\"user_id\""},
	RoomID:         whereHelpernull_Int64{field: "\"sessions\".\"room_id\""},
	GatewayID:      whereHelpernull_Int64{field: "\"sessions\".\"gateway_id\""},
	GatewaySession: whereHelpernull_Int64{field: "\"sessions\".\"gateway_session\""},
	GatewayHandle:  whereHelpernull_Int64{field: "\"sessions\".\"gateway_handle\""},
	GatewayFeed:    whereHelpernull_Int64{field: "\"sessions\".\"gateway_feed\""},
	Display:        whereHelpernull_String{field: "\"sessions\".\"display\""},
	Camera:         whereHelperbool{field: "\"sessions\".\"camera\""},
	Question:       whereHelperbool{field: "\"sessions\".\"question\""},
	SelfTest:       whereHelperbool{field: "\"sessions\".\"self_test\""},
	SoundTest:      whereHelperbool{field: "\"sessions\".\"sound_test\""},
	UserAgent:      whereHelpernull_String{field: "\"sessions\".\"user_agent\""},
	IPAddress:      whereHelpernull_String{field: "\"sessions\".\"ip_address\""},
	Properties:     whereHelpernull_JSON{field: "\"sessions\".\"properties\""},
	CreatedAt:      whereHelpertime_Time{field: "\"sessions\".\"created_at\""},
	UpdatedAt:      whereHelpernull_Time{field: "\"sessions\".\"updated_at\""},
	RemovedAt:      whereHelpernull_Time{field: "\"sessions\".\"removed_at\""},
}

// SessionRels is where relationship names are stored.
var SessionRels = struct {
	Gateway string
	Room    string
	User    string
}{
	Gateway: "Gateway",
	Room:    "Room",
	User:    "User",
}

// sessionR is where relationships are stored.
type sessionR struct {
	Gateway *Gateway
	Room    *Room
	User    *User
}

// NewStruct creates a new relationship struct
func (*sessionR) NewStruct() *sessionR {
	return &sessionR{}
}

// sessionL is where Load methods for each relationship are stored.
type sessionL struct{}

var (
	sessionAllColumns            = []string{"id", "user_id", "room_id", "gateway_id", "gateway_session", "gateway_handle", "gateway_feed", "display", "camera", "question", "self_test", "sound_test", "user_agent", "ip_address", "properties", "created_at", "updated_at", "removed_at"}
	sessionColumnsWithoutDefault = []string{"user_id", "room_id", "gateway_id", "gateway_session", "gateway_handle", "gateway_feed", "display", "user_agent", "ip_address", "properties", "updated_at", "removed_at"}
	sessionColumnsWithDefault    = []string{"id", "camera", "question", "self_test", "sound_test", "created_at"}
	sessionPrimaryKeyColumns     = []string{"id"}
)

type (
	// SessionSlice is an alias for a slice of pointers to Session.
	// This should generally be used opposed to []Session.
	SessionSlice []*Session

	sessionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sessionType                 = reflect.TypeOf(&Session{})
	sessionMapping              = queries.MakeStructMapping(sessionType)
	sessionPrimaryKeyMapping, _ = queries.BindMapping(sessionType, sessionMapping, sessionPrimaryKeyColumns)
	sessionInsertCacheMut       sync.RWMutex
	sessionInsertCache          = make(map[string]insertCache)
	sessionUpdateCacheMut       sync.RWMutex
	sessionUpdateCache          = make(map[string]updateCache)
	sessionUpsertCacheMut       sync.RWMutex
	sessionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single session record from the query.
func (q sessionQuery) One(exec boil.Executor) (*Session, error) {
	o := &Session{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for sessions")
	}

	return o, nil
}

// All returns all Session records from the query.
func (q sessionQuery) All(exec boil.Executor) (SessionSlice, error) {
	var o []*Session

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Session slice")
	}

	return o, nil
}

// Count returns the count of all Session records in the query.
func (q sessionQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count sessions rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q sessionQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if sessions exists")
	}

	return count > 0, nil
}

// Gateway pointed to by the foreign key.
func (o *Session) Gateway(mods ...qm.QueryMod) gatewayQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.GatewayID),
	}

	queryMods = append(queryMods, mods...)

	query := Gateways(queryMods...)
	queries.SetFrom(query.Query, "\"gateways\"")

	return query
}

// Room pointed to by the foreign key.
func (o *Session) Room(mods ...qm.QueryMod) roomQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.RoomID),
	}

	queryMods = append(queryMods, mods...)

	query := Rooms(queryMods...)
	queries.SetFrom(query.Query, "\"rooms\"")

	return query
}

// User pointed to by the foreign key.
func (o *Session) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadGateway allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (sessionL) LoadGateway(e boil.Executor, singular bool, maybeSession interface{}, mods queries.Applicator) error {
	var slice []*Session
	var object *Session

	if singular {
		object = maybeSession.(*Session)
	} else {
		slice = *maybeSession.(*[]*Session)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &sessionR{}
		}
		if !queries.IsNil(object.GatewayID) {
			args = append(args, object.GatewayID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &sessionR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.GatewayID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.GatewayID) {
				args = append(args, obj.GatewayID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`gateways`), qm.WhereIn(`gateways.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Gateway")
	}

	var resultSlice []*Gateway
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Gateway")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for gateways")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for gateways")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Gateway = foreign
		if foreign.R == nil {
			foreign.R = &gatewayR{}
		}
		foreign.R.Sessions = append(foreign.R.Sessions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.GatewayID, foreign.ID) {
				local.R.Gateway = foreign
				if foreign.R == nil {
					foreign.R = &gatewayR{}
				}
				foreign.R.Sessions = append(foreign.R.Sessions, local)
				break
			}
		}
	}

	return nil
}

// LoadRoom allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (sessionL) LoadRoom(e boil.Executor, singular bool, maybeSession interface{}, mods queries.Applicator) error {
	var slice []*Session
	var object *Session

	if singular {
		object = maybeSession.(*Session)
	} else {
		slice = *maybeSession.(*[]*Session)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &sessionR{}
		}
		if !queries.IsNil(object.RoomID) {
			args = append(args, object.RoomID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &sessionR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.RoomID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.RoomID) {
				args = append(args, obj.RoomID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`rooms`), qm.WhereIn(`rooms.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Room")
	}

	var resultSlice []*Room
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Room")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for rooms")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for rooms")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Room = foreign
		if foreign.R == nil {
			foreign.R = &roomR{}
		}
		foreign.R.Sessions = append(foreign.R.Sessions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.RoomID, foreign.ID) {
				local.R.Room = foreign
				if foreign.R == nil {
					foreign.R = &roomR{}
				}
				foreign.R.Sessions = append(foreign.R.Sessions, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (sessionL) LoadUser(e boil.Executor, singular bool, maybeSession interface{}, mods queries.Applicator) error {
	var slice []*Session
	var object *Session

	if singular {
		object = maybeSession.(*Session)
	} else {
		slice = *maybeSession.(*[]*Session)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &sessionR{}
		}
		args = append(args, object.UserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &sessionR{}
			}

			for _, a := range args {
				if a == obj.UserID {
					continue Outer
				}
			}

			args = append(args, obj.UserID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`users`), qm.WhereIn(`users.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.Sessions = append(foreign.R.Sessions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.Sessions = append(foreign.R.Sessions, local)
				break
			}
		}
	}

	return nil
}

// SetGateway of the session to the related item.
// Sets o.R.Gateway to related.
// Adds o to related.R.Sessions.
func (o *Session) SetGateway(exec boil.Executor, insert bool, related *Gateway) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sessions\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"gateway_id"}),
		strmangle.WhereClause("\"", "\"", 2, sessionPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.GatewayID, related.ID)
	if o.R == nil {
		o.R = &sessionR{
			Gateway: related,
		}
	} else {
		o.R.Gateway = related
	}

	if related.R == nil {
		related.R = &gatewayR{
			Sessions: SessionSlice{o},
		}
	} else {
		related.R.Sessions = append(related.R.Sessions, o)
	}

	return nil
}

// RemoveGateway relationship.
// Sets o.R.Gateway to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Session) RemoveGateway(exec boil.Executor, related *Gateway) error {
	var err error

	queries.SetScanner(&o.GatewayID, nil)
	if _, err = o.Update(exec, boil.Whitelist("gateway_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.Gateway = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Sessions {
		if queries.Equal(o.GatewayID, ri.GatewayID) {
			continue
		}

		ln := len(related.R.Sessions)
		if ln > 1 && i < ln-1 {
			related.R.Sessions[i] = related.R.Sessions[ln-1]
		}
		related.R.Sessions = related.R.Sessions[:ln-1]
		break
	}
	return nil
}

// SetRoom of the session to the related item.
// Sets o.R.Room to related.
// Adds o to related.R.Sessions.
func (o *Session) SetRoom(exec boil.Executor, insert bool, related *Room) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sessions\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"room_id"}),
		strmangle.WhereClause("\"", "\"", 2, sessionPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.RoomID, related.ID)
	if o.R == nil {
		o.R = &sessionR{
			Room: related,
		}
	} else {
		o.R.Room = related
	}

	if related.R == nil {
		related.R = &roomR{
			Sessions: SessionSlice{o},
		}
	} else {
		related.R.Sessions = append(related.R.Sessions, o)
	}

	return nil
}

// RemoveRoom relationship.
// Sets o.R.Room to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Session) RemoveRoom(exec boil.Executor, related *Room) error {
	var err error

	queries.SetScanner(&o.RoomID, nil)
	if _, err = o.Update(exec, boil.Whitelist("room_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.Room = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Sessions {
		if queries.Equal(o.RoomID, ri.RoomID) {
			continue
		}

		ln := len(related.R.Sessions)
		if ln > 1 && i < ln-1 {
			related.R.Sessions[i] = related.R.Sessions[ln-1]
		}
		related.R.Sessions = related.R.Sessions[:ln-1]
		break
	}
	return nil
}

// SetUser of the session to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Sessions.
func (o *Session) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sessions\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, sessionPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.ID
	if o.R == nil {
		o.R = &sessionR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			Sessions: SessionSlice{o},
		}
	} else {
		related.R.Sessions = append(related.R.Sessions, o)
	}

	return nil
}

// Sessions retrieves all the records using an executor.
func Sessions(mods ...qm.QueryMod) sessionQuery {
	mods = append(mods, qm.From("\"sessions\""))
	return sessionQuery{NewQuery(mods...)}
}

// FindSession retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSession(exec boil.Executor, iD int64, selectCols ...string) (*Session, error) {
	sessionObj := &Session{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"sessions\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, sessionObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from sessions")
	}

	return sessionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Session) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no sessions provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(sessionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	sessionInsertCacheMut.RLock()
	cache, cached := sessionInsertCache[key]
	sessionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			sessionAllColumns,
			sessionColumnsWithDefault,
			sessionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(sessionType, sessionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sessionType, sessionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"sessions\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"sessions\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into sessions")
	}

	if !cached {
		sessionInsertCacheMut.Lock()
		sessionInsertCache[key] = cache
		sessionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Session.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Session) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	sessionUpdateCacheMut.RLock()
	cache, cached := sessionUpdateCache[key]
	sessionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			sessionAllColumns,
			sessionPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update sessions, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"sessions\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sessionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sessionType, sessionMapping, append(wl, sessionPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update sessions row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for sessions")
	}

	if !cached {
		sessionUpdateCacheMut.Lock()
		sessionUpdateCache[key] = cache
		sessionUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q sessionQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for sessions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for sessions")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SessionSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"sessions\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, sessionPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in session slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all session")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Session) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no sessions provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(sessionColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	sessionUpsertCacheMut.RLock()
	cache, cached := sessionUpsertCache[key]
	sessionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			sessionAllColumns,
			sessionColumnsWithDefault,
			sessionColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			sessionAllColumns,
			sessionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert sessions, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(sessionPrimaryKeyColumns))
			copy(conflict, sessionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"sessions\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(sessionType, sessionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sessionType, sessionMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert sessions")
	}

	if !cached {
		sessionUpsertCacheMut.Lock()
		sessionUpsertCache[key] = cache
		sessionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Session record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Session) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Session provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sessionPrimaryKeyMapping)
	sql := "DELETE FROM \"sessions\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from sessions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for sessions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q sessionQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no sessionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from sessions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for sessions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SessionSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"sessions\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sessionPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from session slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for sessions")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Session) Reload(exec boil.Executor) error {
	ret, err := FindSession(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SessionSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SessionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sessionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"sessions\".* FROM \"sessions\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sessionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SessionSlice")
	}

	*o = slice

	return nil
}

// SessionExists checks if the Session row exists.
func SessionExists(exec boil.Executor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"sessions\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if sessions exists")
	}

	return exists, nil
}
