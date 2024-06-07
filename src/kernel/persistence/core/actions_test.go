package core

import (
	"duck/kernel/persistence/core/dialects"
	"reflect"
	"strings"
	"testing"

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
}

func TestQuery(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Error(err)
	}

	engine.Exec("CREATE TABLE test_struct (a INTEGER, b TEXT);")
	engine.Exec("INSERT INTO test_struct VALUES(0, '0');")

	var models []TestStruct

	err = engine.Query("SELECT * FROM test_struct;", Parse(reflect.TypeOf(TestStruct{}), engine.dialect), &models)

	if err != nil {
		t.Error(err)
	}

	if len(models) != 1 || models[0].A != 0 || models[0].B != "0" {
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

	engine.Exec("INSERT INTO test_struct VALUES(0, '0');")
	
	if err != nil {
		t.Error(err)
	}

	var tests []TestStruct
	err = engine.Select(&tests)

	if err != nil {
		t.Error(err)
	}

	if len(tests) != 1 || tests[0].A != 0 || tests[0].B != "0" {
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

	err = engine.Insert(TestStruct{A: 0, B: "0"})
	if err != nil {
		t.Error(err)
	}
	var models []TestStruct
	engine.Select(&models)
	if len(models) != 1 || models[0].A != 0 || models[0].B != "0" {
		t.Error("insert test failed")
	}
}

