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

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// RoomStatistic is an object representing the database table.
type RoomStatistic struct {
	RoomID int64 `boil:"room_id" json:"room_id" toml:"room_id" yaml:"room_id"`
	OnAir  int   `boil:"on_air" json:"on_air" toml:"on_air" yaml:"on_air"`

	R *roomStatisticR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L roomStatisticL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var RoomStatisticColumns = struct {
	RoomID string
	OnAir  string
}{
	RoomID: "room_id",
	OnAir:  "on_air",
}

// Generated where

var RoomStatisticWhere = struct {
	RoomID whereHelperint64
	OnAir  whereHelperint
}{
	RoomID: whereHelperint64{field: "\"room_statistics\".\"room_id\""},
	OnAir:  whereHelperint{field: "\"room_statistics\".\"on_air\""},
}

// RoomStatisticRels is where relationship names are stored.
var RoomStatisticRels = struct {
	Room string
}{
	Room: "Room",
}

// roomStatisticR is where relationships are stored.
type roomStatisticR struct {
	Room *Room
}

// NewStruct creates a new relationship struct
func (*roomStatisticR) NewStruct() *roomStatisticR {
	return &roomStatisticR{}
}

// roomStatisticL is where Load methods for each relationship are stored.
type roomStatisticL struct{}

var (
	roomStatisticAllColumns            = []string{"room_id", "on_air"}
	roomStatisticColumnsWithoutDefault = []string{"room_id"}
	roomStatisticColumnsWithDefault    = []string{"on_air"}
	roomStatisticPrimaryKeyColumns     = []string{"room_id"}
)

type (
	// RoomStatisticSlice is an alias for a slice of pointers to RoomStatistic.
	// This should generally be used opposed to []RoomStatistic.
	RoomStatisticSlice []*RoomStatistic

	roomStatisticQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	roomStatisticType                 = reflect.TypeOf(&RoomStatistic{})
	roomStatisticMapping              = queries.MakeStructMapping(roomStatisticType)
	roomStatisticPrimaryKeyMapping, _ = queries.BindMapping(roomStatisticType, roomStatisticMapping, roomStatisticPrimaryKeyColumns)
	roomStatisticInsertCacheMut       sync.RWMutex
	roomStatisticInsertCache          = make(map[string]insertCache)
	roomStatisticUpdateCacheMut       sync.RWMutex
	roomStatisticUpdateCache          = make(map[string]updateCache)
	roomStatisticUpsertCacheMut       sync.RWMutex
	roomStatisticUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single roomStatistic record from the query.
func (q roomStatisticQuery) One(exec boil.Executor) (*RoomStatistic, error) {
	o := &RoomStatistic{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for room_statistics")
	}

	return o, nil
}

// All returns all RoomStatistic records from the query.
func (q roomStatisticQuery) All(exec boil.Executor) (RoomStatisticSlice, error) {
	var o []*RoomStatistic

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to RoomStatistic slice")
	}

	return o, nil
}

// Count returns the count of all RoomStatistic records in the query.
func (q roomStatisticQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count room_statistics rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q roomStatisticQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if room_statistics exists")
	}

	return count > 0, nil
}

// Room pointed to by the foreign key.
func (o *RoomStatistic) Room(mods ...qm.QueryMod) roomQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.RoomID),
	}

	queryMods = append(queryMods, mods...)

	query := Rooms(queryMods...)
	queries.SetFrom(query.Query, "\"rooms\"")

	return query
}

// LoadRoom allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (roomStatisticL) LoadRoom(e boil.Executor, singular bool, maybeRoomStatistic interface{}, mods queries.Applicator) error {
	var slice []*RoomStatistic
	var object *RoomStatistic

	if singular {
		object = maybeRoomStatistic.(*RoomStatistic)
	} else {
		slice = *maybeRoomStatistic.(*[]*RoomStatistic)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &roomStatisticR{}
		}
		args = append(args, object.RoomID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &roomStatisticR{}
			}

			for _, a := range args {
				if a == obj.RoomID {
					continue Outer
				}
			}

			args = append(args, obj.RoomID)

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
		foreign.R.RoomStatistic = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.RoomID == foreign.ID {
				local.R.Room = foreign
				if foreign.R == nil {
					foreign.R = &roomR{}
				}
				foreign.R.RoomStatistic = local
				break
			}
		}
	}

	return nil
}

// SetRoom of the roomStatistic to the related item.
// Sets o.R.Room to related.
// Adds o to related.R.RoomStatistic.
func (o *RoomStatistic) SetRoom(exec boil.Executor, insert bool, related *Room) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"room_statistics\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"room_id"}),
		strmangle.WhereClause("\"", "\"", 2, roomStatisticPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.RoomID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.RoomID = related.ID
	if o.R == nil {
		o.R = &roomStatisticR{
			Room: related,
		}
	} else {
		o.R.Room = related
	}

	if related.R == nil {
		related.R = &roomR{
			RoomStatistic: o,
		}
	} else {
		related.R.RoomStatistic = o
	}

	return nil
}

// RoomStatistics retrieves all the records using an executor.
func RoomStatistics(mods ...qm.QueryMod) roomStatisticQuery {
	mods = append(mods, qm.From("\"room_statistics\""))
	return roomStatisticQuery{NewQuery(mods...)}
}

// FindRoomStatistic retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRoomStatistic(exec boil.Executor, roomID int64, selectCols ...string) (*RoomStatistic, error) {
	roomStatisticObj := &RoomStatistic{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"room_statistics\" where \"room_id\"=$1", sel,
	)

	q := queries.Raw(query, roomID)

	err := q.Bind(nil, exec, roomStatisticObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from room_statistics")
	}

	return roomStatisticObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *RoomStatistic) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no room_statistics provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(roomStatisticColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	roomStatisticInsertCacheMut.RLock()
	cache, cached := roomStatisticInsertCache[key]
	roomStatisticInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			roomStatisticAllColumns,
			roomStatisticColumnsWithDefault,
			roomStatisticColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(roomStatisticType, roomStatisticMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(roomStatisticType, roomStatisticMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"room_statistics\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"room_statistics\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into room_statistics")
	}

	if !cached {
		roomStatisticInsertCacheMut.Lock()
		roomStatisticInsertCache[key] = cache
		roomStatisticInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the RoomStatistic.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *RoomStatistic) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	roomStatisticUpdateCacheMut.RLock()
	cache, cached := roomStatisticUpdateCache[key]
	roomStatisticUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			roomStatisticAllColumns,
			roomStatisticPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update room_statistics, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"room_statistics\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, roomStatisticPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(roomStatisticType, roomStatisticMapping, append(wl, roomStatisticPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update room_statistics row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for room_statistics")
	}

	if !cached {
		roomStatisticUpdateCacheMut.Lock()
		roomStatisticUpdateCache[key] = cache
		roomStatisticUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q roomStatisticQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for room_statistics")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for room_statistics")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RoomStatisticSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roomStatisticPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"room_statistics\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, roomStatisticPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in roomStatistic slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all roomStatistic")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *RoomStatistic) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no room_statistics provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(roomStatisticColumnsWithDefault, o)

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

	roomStatisticUpsertCacheMut.RLock()
	cache, cached := roomStatisticUpsertCache[key]
	roomStatisticUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			roomStatisticAllColumns,
			roomStatisticColumnsWithDefault,
			roomStatisticColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			roomStatisticAllColumns,
			roomStatisticPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert room_statistics, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(roomStatisticPrimaryKeyColumns))
			copy(conflict, roomStatisticPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"room_statistics\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(roomStatisticType, roomStatisticMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(roomStatisticType, roomStatisticMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert room_statistics")
	}

	if !cached {
		roomStatisticUpsertCacheMut.Lock()
		roomStatisticUpsertCache[key] = cache
		roomStatisticUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single RoomStatistic record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *RoomStatistic) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no RoomStatistic provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), roomStatisticPrimaryKeyMapping)
	sql := "DELETE FROM \"room_statistics\" WHERE \"room_id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from room_statistics")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for room_statistics")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q roomStatisticQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no roomStatisticQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from room_statistics")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for room_statistics")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RoomStatisticSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roomStatisticPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"room_statistics\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, roomStatisticPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from roomStatistic slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for room_statistics")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *RoomStatistic) Reload(exec boil.Executor) error {
	ret, err := FindRoomStatistic(exec, o.RoomID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RoomStatisticSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := RoomStatisticSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roomStatisticPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"room_statistics\".* FROM \"room_statistics\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, roomStatisticPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RoomStatisticSlice")
	}

	*o = slice

	return nil
}

// RoomStatisticExists checks if the RoomStatistic row exists.
func RoomStatisticExists(exec boil.Executor, roomID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"room_statistics\" where \"room_id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, roomID)
	}
	row := exec.QueryRow(sql, roomID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if room_statistics exists")
	}

	return exists, nil
}
