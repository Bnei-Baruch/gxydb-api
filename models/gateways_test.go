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

func testGateways(t *testing.T) {
	t.Parallel()

	query := Gateways()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testGatewaysDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
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

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testGatewaysQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Gateways().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testGatewaysSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := GatewaySlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testGatewaysExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := GatewayExists(tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Gateway exists: %s", err)
	}
	if !e {
		t.Errorf("Expected GatewayExists to return true, but got false.")
	}
}

func testGatewaysFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	gatewayFound, err := FindGateway(tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if gatewayFound == nil {
		t.Error("want a record, got nil")
	}
}

func testGatewaysBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Gateways().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testGatewaysOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Gateways().One(tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testGatewaysAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	gatewayOne := &Gateway{}
	gatewayTwo := &Gateway{}
	if err = randomize.Struct(seed, gatewayOne, gatewayDBTypes, false, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}
	if err = randomize.Struct(seed, gatewayTwo, gatewayDBTypes, false, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = gatewayOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = gatewayTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Gateways().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testGatewaysCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	gatewayOne := &Gateway{}
	gatewayTwo := &Gateway{}
	if err = randomize.Struct(seed, gatewayOne, gatewayDBTypes, false, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}
	if err = randomize.Struct(seed, gatewayTwo, gatewayDBTypes, false, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = gatewayOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = gatewayTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testGatewaysInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testGatewaysInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(gatewayColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testGatewayToManyCompositesRooms(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Gateway
	var b, c CompositesRoom

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, compositesRoomDBTypes, false, compositesRoomColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.GatewayID = a.ID
	c.GatewayID = a.ID

	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.CompositesRooms().All(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.GatewayID == b.GatewayID {
			bFound = true
		}
		if v.GatewayID == c.GatewayID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := GatewaySlice{&a}
	if err = a.L.LoadCompositesRooms(tx, false, (*[]*Gateway)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CompositesRooms); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.CompositesRooms = nil
	if err = a.L.LoadCompositesRooms(tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CompositesRooms); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testGatewayToManyDefaultGatewayRooms(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Gateway
	var b, c Room

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, roomDBTypes, false, roomColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, roomDBTypes, false, roomColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.DefaultGatewayID = a.ID
	c.DefaultGatewayID = a.ID

	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.DefaultGatewayRooms().All(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.DefaultGatewayID == b.DefaultGatewayID {
			bFound = true
		}
		if v.DefaultGatewayID == c.DefaultGatewayID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := GatewaySlice{&a}
	if err = a.L.LoadDefaultGatewayRooms(tx, false, (*[]*Gateway)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DefaultGatewayRooms); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.DefaultGatewayRooms = nil
	if err = a.L.LoadDefaultGatewayRooms(tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DefaultGatewayRooms); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testGatewayToManySessions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Gateway
	var b, c Session

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, sessionDBTypes, false, sessionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, sessionDBTypes, false, sessionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&b.GatewayID, a.ID)
	queries.Assign(&c.GatewayID, a.ID)
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Sessions().All(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if queries.Equal(v.GatewayID, b.GatewayID) {
			bFound = true
		}
		if queries.Equal(v.GatewayID, c.GatewayID) {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := GatewaySlice{&a}
	if err = a.L.LoadSessions(tx, false, (*[]*Gateway)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Sessions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Sessions = nil
	if err = a.L.LoadSessions(tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Sessions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testGatewayToManyAddOpCompositesRooms(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Gateway
	var b, c, d, e CompositesRoom

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, gatewayDBTypes, false, strmangle.SetComplement(gatewayPrimaryKeyColumns, gatewayColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*CompositesRoom{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, compositesRoomDBTypes, false, strmangle.SetComplement(compositesRoomPrimaryKeyColumns, compositesRoomColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*CompositesRoom{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddCompositesRooms(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.GatewayID {
			t.Error("foreign key was wrong value", a.ID, first.GatewayID)
		}
		if a.ID != second.GatewayID {
			t.Error("foreign key was wrong value", a.ID, second.GatewayID)
		}

		if first.R.Gateway != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Gateway != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.CompositesRooms[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.CompositesRooms[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.CompositesRooms().Count(tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testGatewayToManyAddOpDefaultGatewayRooms(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Gateway
	var b, c, d, e Room

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, gatewayDBTypes, false, strmangle.SetComplement(gatewayPrimaryKeyColumns, gatewayColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Room{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, roomDBTypes, false, strmangle.SetComplement(roomPrimaryKeyColumns, roomColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Room{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDefaultGatewayRooms(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.DefaultGatewayID {
			t.Error("foreign key was wrong value", a.ID, first.DefaultGatewayID)
		}
		if a.ID != second.DefaultGatewayID {
			t.Error("foreign key was wrong value", a.ID, second.DefaultGatewayID)
		}

		if first.R.DefaultGateway != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.DefaultGateway != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.DefaultGatewayRooms[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.DefaultGatewayRooms[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.DefaultGatewayRooms().Count(tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testGatewayToManyAddOpSessions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Gateway
	var b, c, d, e Session

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, gatewayDBTypes, false, strmangle.SetComplement(gatewayPrimaryKeyColumns, gatewayColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Session{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sessionDBTypes, false, strmangle.SetComplement(sessionPrimaryKeyColumns, sessionColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Session{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddSessions(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if !queries.Equal(a.ID, first.GatewayID) {
			t.Error("foreign key was wrong value", a.ID, first.GatewayID)
		}
		if !queries.Equal(a.ID, second.GatewayID) {
			t.Error("foreign key was wrong value", a.ID, second.GatewayID)
		}

		if first.R.Gateway != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Gateway != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Sessions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Sessions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Sessions().Count(tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testGatewayToManySetOpSessions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Gateway
	var b, c, d, e Session

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, gatewayDBTypes, false, strmangle.SetComplement(gatewayPrimaryKeyColumns, gatewayColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Session{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sessionDBTypes, false, strmangle.SetComplement(sessionPrimaryKeyColumns, sessionColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetSessions(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Sessions().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetSessions(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Sessions().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.GatewayID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.GatewayID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.GatewayID) {
		t.Error("foreign key was wrong value", a.ID, d.GatewayID)
	}
	if !queries.Equal(a.ID, e.GatewayID) {
		t.Error("foreign key was wrong value", a.ID, e.GatewayID)
	}

	if b.R.Gateway != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Gateway != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Gateway != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Gateway != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.Sessions[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Sessions[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testGatewayToManyRemoveOpSessions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Gateway
	var b, c, d, e Session

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, gatewayDBTypes, false, strmangle.SetComplement(gatewayPrimaryKeyColumns, gatewayColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Session{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sessionDBTypes, false, strmangle.SetComplement(sessionPrimaryKeyColumns, sessionColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddSessions(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Sessions().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveSessions(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Sessions().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.GatewayID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.GatewayID) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Gateway != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Gateway != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Gateway != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Gateway != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.Sessions) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Sessions[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Sessions[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testGatewaysReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
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

func testGatewaysReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := GatewaySlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testGatewaysSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Gateways().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	gatewayDBTypes = map[string]string{`ID`: `bigint`, `Name`: `character varying`, `Description`: `character varying`, `URL`: `character varying`, `AdminURL`: `character varying`, `AdminPassword`: `text`, `Disabled`: `boolean`, `Properties`: `jsonb`, `CreatedAt`: `timestamp with time zone`, `UpdatedAt`: `timestamp with time zone`, `RemovedAt`: `timestamp with time zone`}
	_              = bytes.MinRead
)

func testGatewaysUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(gatewayPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(gatewayAllColumns) == len(gatewayPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testGatewaysSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(gatewayAllColumns) == len(gatewayPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Gateway{}
	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, gatewayDBTypes, true, gatewayPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(gatewayAllColumns, gatewayPrimaryKeyColumns) {
		fields = gatewayAllColumns
	} else {
		fields = strmangle.SetComplement(
			gatewayAllColumns,
			gatewayPrimaryKeyColumns,
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

	slice := GatewaySlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testGatewaysUpsert(t *testing.T) {
	t.Parallel()

	if len(gatewayAllColumns) == len(gatewayPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Gateway{}
	if err = randomize.Struct(seed, &o, gatewayDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Gateway: %s", err)
	}

	count, err := Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, gatewayDBTypes, false, gatewayPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Gateway struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Gateway: %s", err)
	}

	count, err = Gateways().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}