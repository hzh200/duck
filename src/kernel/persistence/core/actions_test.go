package core

import (
	"duck/kernel/persistence/core/dialects"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestExec(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Fatal(err)
	}

	_, err = engine.Exec("SELECT * FROM sqlite_master;")

	if err != nil {
		t.Fatal(err)
	}
}

type TestStruct struct {
	A int
	B string
	C []int64
	D [][]float64
	E map[string]int64
	F map[string][]bool
	G time.Time
}

func TestQuery(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Error(err)
	}

	_, err = engine.Exec("CREATE TABLE test_struct (a INTEGER, b TEXT, c TEXT, d TEXT, e TEXT, f TEXT, g TIMESTAMP);")

	if err != nil {
		t.Error(err)
	}

	_, err = engine.Exec("INSERT INTO test_struct VALUES(0, '0', '[1,2,3]', '[[1.0, 1.0], [2.0, 2.0], [3.0, 3.0]]', '{\"1\":1,\"2\":2}', '{\"1\":[true, false]}', DATE('2024-01-01 00:00:00'));")

	if err != nil {
		t.Error(err)
	}

	var models []TestStruct

	err = engine.Query("SELECT * FROM test_struct;", Parse(reflect.TypeOf(TestStruct{}), engine.dialect), &models)

	if err != nil {
		t.Error(err)
	}

	t.Log("models:", models)

	if len(models) != 1 {
		t.Error("query test failed")
	}

	model := models[0]
	if model.A != 0 || 
		model.B != "0" || 
		!reflect.DeepEqual(model.C, []int64{1, 2, 3}) || 
		!reflect.DeepEqual(model.D, [][]float64{{1.0, 1.0}, {2.0, 2.0}, {3.0, 3.0}}) || 
		!reflect.DeepEqual(model.E, map[string]int64{"1": 1, "2": 2}) ||
		!reflect.DeepEqual(model.F, map[string][]bool{"1": {true, false}}) || 
		!reflect.DeepEqual(model.G, time.Date(2024, 01, 01, 00, 00, 00, 000, time.UTC)) {
		t.Error("query test failed")
	}
}

func TestCreate(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Error(err)
	}

	_, err = engine.Create(TestStruct{})

	if err != nil {
		t.Error(err)
	}

	var models []dialects.SQLiteTableInfo
	err = engine.Query("SELECT * FROM sqlite_master WHERE name='test_struct';", Parse(reflect.TypeOf(dialects.SQLiteTableInfo{}), engine.dialect), &models)

	if err != nil {
		t.Error(err)
	}

	if len(models) != 1 || strings.Compare(models[0].TblName, "test_struct") != 0 {
		t.Error("create test failed")
	}
}

func TestSelect(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Error(err)
	}

	_, err = engine.Create(TestStruct{})

	if err != nil {
		t.Error(err)
	}

	engine.Exec("INSERT INTO test_struct VALUES(0, '0', '[1,2,3]', '[[1.0, 1.0], [2.0, 2.0], [3.0, 3.0]]', '{\"1\":1,\"2\":2}', '{\"1\":[true, false]}', DATE('2024-01-01 00:00:00'));")

	if err != nil {
		t.Error(err)
	}

	var models []TestStruct
	err = engine.Select(&models)

	if err != nil {
		t.Error(err)
	}

	if len(models) != 1 {
		t.Error("select test failed")
	}

	model := models[0]

	if model.A != 0 || 
		model.B != "0" || 
		!reflect.DeepEqual(model.C, []int64{1, 2, 3}) || 
		!reflect.DeepEqual(model.D, [][]float64{{1.0, 1.0}, {2.0, 2.0}, {3.0, 3.0}}) || 
		!reflect.DeepEqual(model.E, map[string]int64{"1": 1, "2": 2}) ||
		!reflect.DeepEqual(model.F, map[string][]bool{"1": {true, false}}) || 
		!reflect.DeepEqual(model.G, time.Date(2024, 01, 01, 00, 00, 00, 000, time.UTC)) {
		t.Error("select test failed")
	}
}

func TestInsert(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Error(err)
	}

	_, err = engine.Create(TestStruct{})

	if err != nil {
		t.Error(err)
	}

	err = engine.Insert(TestStruct{
		A: 0, 
		B: "0", 
		C: []int64{1, 2, 3}, 
		D: [][]float64{{1.0, 1.0}, {2.0, 2.0}, {3.0, 3.0}}, 
		E: map[string]int64{"1": 1, "2": 2}, 
		F: map[string][]bool{"1": {true, false}},
		G: time.Date(2024, 01, 01, 00, 00, 00, 000, time.UTC),
	})
	if err != nil {
		t.Error(err)
	}
	var models []TestStruct
	engine.Select(&models)
	if len(models) != 1 {
		t.Error("insert test failed")
	}

	model := models[0]

	if model.A != 0 || 
		model.B != "0" || 
		!reflect.DeepEqual(model.C, []int64{1, 2, 3}) || 
		!reflect.DeepEqual(model.D, [][]float64{{1.0, 1.0}, {2.0, 2.0}, {3.0, 3.0}}) || 
		!reflect.DeepEqual(model.E, map[string]int64{"1": 1, "2": 2}) || 
		!reflect.DeepEqual(model.F, map[string][]bool{"1": {true, false}}) || 
		!reflect.DeepEqual(model.G, time.Date(2024, 01, 01, 00, 00, 00, 000, time.UTC)) {
		t.Error("insert test failed")
	}
}
