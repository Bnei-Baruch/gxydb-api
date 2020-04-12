// Code generated by SQLBoiler 3.6.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testCompositesRooms(t *testing.T) {
	t.Parallel()

	query := CompositesRooms()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testCompositesRoomsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCompositesRoomsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := CompositesRooms().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCompositesRoomsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := CompositesRoomSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCompositesRoomsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := CompositesRoomExists(tx, o.CompositeID, o.RoomID, o.GatewayID)
	if err != nil {
		t.Errorf("Unable to check if CompositesRoom exists: %s", err)
	}
	if !e {
		t.Errorf("Expected CompositesRoomExists to return true, but got false.")
	}
}

func testCompositesRoomsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	compositesRoomFound, err := FindCompositesRoom(tx, o.CompositeID, o.RoomID, o.GatewayID)
	if err != nil {
		t.Error(err)
	}

	if compositesRoomFound == nil {
		t.Error("want a record, got nil")
	}
}

func testCompositesRoomsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = CompositesRooms().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testCompositesRoomsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := CompositesRooms().One(tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testCompositesRoomsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	compositesRoomOne := &CompositesRoom{}
	compositesRoomTwo := &CompositesRoom{}
	if err = randomize.Struct(seed, compositesRoomOne, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}
	if err = randomize.Struct(seed, compositesRoomTwo, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = compositesRoomOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = compositesRoomTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := CompositesRooms().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testCompositesRoomsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	compositesRoomOne := &CompositesRoom{}
	compositesRoomTwo := &CompositesRoom{}
	if err = randomize.Struct(seed, compositesRoomOne, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}
	if err = randomize.Struct(seed, compositesRoomTwo, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = compositesRoomOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = compositesRoomTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testCompositesRoomsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCompositesRoomsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(compositesRoomColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCompositesRoomToOneCompositeUsingComposite(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local CompositesRoom
	var foreign Composite

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, compositeDBTypes, false, compositeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Composite struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.CompositeID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Composite().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CompositesRoomSlice{&local}
	if err = local.L.LoadComposite(tx, false, (*[]*CompositesRoom)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Composite == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Composite = nil
	if err = local.L.LoadComposite(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Composite == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCompositesRoomToOneGatewayUsingGateway(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local CompositesRoom
	var foreign Gateway

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, gatewayDBTypes, false, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.GatewayID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Gateway().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CompositesRoomSlice{&local}
	if err = local.L.LoadGateway(tx, false, (*[]*CompositesRoom)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Gateway == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Gateway = nil
	if err = local.L.LoadGateway(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Gateway == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCompositesRoomToOneRoomUsingRoom(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local CompositesRoom
	var foreign Room

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, roomDBTypes, false, roomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Room struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.RoomID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Room().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CompositesRoomSlice{&local}
	if err = local.L.LoadRoom(tx, false, (*[]*CompositesRoom)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Room == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Room = nil
	if err = local.L.LoadRoom(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Room == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCompositesRoomToOneSetOpCompositeUsingComposite(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a CompositesRoom
	var b, c Composite

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, compositesRoomDBTypes, false, strmangle.SetComplement(compositesRoomPrimaryKeyColumns, compositesRoomColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, compositeDBTypes, false, strmangle.SetComplement(compositePrimaryKeyColumns, compositeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, compositeDBTypes, false, strmangle.SetComplement(compositePrimaryKeyColumns, compositeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Composite{&b, &c} {
		err = a.SetComposite(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Composite != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.CompositesRooms[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.CompositeID != x.ID {
			t.Error("foreign key was wrong value", a.CompositeID)
		}

		if exists, err := CompositesRoomExists(tx, a.CompositeID, a.RoomID, a.GatewayID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testCompositesRoomToOneSetOpGatewayUsingGateway(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a CompositesRoom
	var b, c Gateway

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, compositesRoomDBTypes, false, strmangle.SetComplement(compositesRoomPrimaryKeyColumns, compositesRoomColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, gatewayDBTypes, false, strmangle.SetComplement(gatewayPrimaryKeyColumns, gatewayColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, gatewayDBTypes, false, strmangle.SetComplement(gatewayPrimaryKeyColumns, gatewayColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Gateway{&b, &c} {
		err = a.SetGateway(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Gateway != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.CompositesRooms[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.GatewayID != x.ID {
			t.Error("foreign key was wrong value", a.GatewayID)
		}

		if exists, err := CompositesRoomExists(tx, a.CompositeID, a.RoomID, a.GatewayID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testCompositesRoomToOneSetOpRoomUsingRoom(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a CompositesRoom
	var b, c Room

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, compositesRoomDBTypes, false, strmangle.SetComplement(compositesRoomPrimaryKeyColumns, compositesRoomColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, roomDBTypes, false, strmangle.SetComplement(roomPrimaryKeyColumns, roomColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, roomDBTypes, false, strmangle.SetComplement(roomPrimaryKeyColumns, roomColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Room{&b, &c} {
		err = a.SetRoom(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Room != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.CompositesRooms[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.RoomID != x.ID {
			t.Error("foreign key was wrong value", a.RoomID)
		}

		if exists, err := CompositesRoomExists(tx, a.CompositeID, a.RoomID, a.GatewayID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}

func testCompositesRoomsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testCompositesRoomsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := CompositesRoomSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testCompositesRoomsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := CompositesRooms().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	compositesRoomDBTypes = map[string]string{`CompositeID`: `bigint`, `RoomID`: `bigint`, `GatewayID`: `bigint`, `Position`: `integer`}
	_                     = bytes.MinRead
)

func testCompositesRoomsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(compositesRoomPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(compositesRoomAllColumns) == len(compositesRoomPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testCompositesRoomsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(compositesRoomAllColumns) == len(compositesRoomPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &CompositesRoom{}
	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, compositesRoomDBTypes, true, compositesRoomPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(compositesRoomAllColumns, compositesRoomPrimaryKeyColumns) {
		fields = compositesRoomAllColumns
	} else {
		fields = strmangle.SetComplement(
			compositesRoomAllColumns,
			compositesRoomPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := CompositesRoomSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testCompositesRoomsUpsert(t *testing.T) {
	t.Parallel()

	if len(compositesRoomAllColumns) == len(compositesRoomPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := CompositesRoom{}
	if err = randomize.Struct(seed, &o, compositesRoomDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert CompositesRoom: %s", err)
	}

	count, err := CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, compositesRoomDBTypes, false, compositesRoomPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CompositesRoom struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert CompositesRoom: %s", err)
	}

	count, err = CompositesRooms().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
