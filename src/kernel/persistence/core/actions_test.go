package core

import (
	"duck/kernel/persistence/core/dialects"
	"fmt"
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
	A int `constraints:"PrimaryKey AutoIncrement"`
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

	fields := []string{
		fmt.Sprintf("a %s %s %s", 
			engine.dialect.DataTypeMapping(reflect.TypeFor[int]()), 
			engine.dialect.PrimaryKey(), 
			engine.dialect.AutoIncrement()),
		fmt.Sprintf("b %s", engine.dialect.DataTypeMapping(reflect.TypeFor[string]())),
		fmt.Sprintf("c %s", engine.dialect.DataTypeMapping(reflect.TypeFor[[]int64]())),
		fmt.Sprintf("d %s", engine.dialect.DataTypeMapping(reflect.TypeFor[[][]float64]())),
		fmt.Sprintf("e %s", engine.dialect.DataTypeMapping(reflect.TypeFor[map[string]int64]())),
		fmt.Sprintf("f %s", engine.dialect.DataTypeMapping(reflect.TypeFor[map[string][]bool]())),
		fmt.Sprintf("g %s", engine.dialect.DataTypeMapping(reflect.TypeFor[time.Time]())),
	}

	createSQLBuilder := strings.Builder{}
	createSQLBuilder.WriteString("CREATE TABLE test_struct(")
	for i, field := range fields {
		createSQLBuilder.WriteString(field)
		if i != len(fields) - 1 {
			createSQLBuilder.WriteString(", ")
		}
	}
	createSQLBuilder.WriteString(");")

	_, err = engine.Exec(createSQLBuilder.String())

	if err != nil {
		t.Error(err)
	}

	values := []string{
		"'0'",
		"'[1,2,3]'",
		"'[[1.0, 1.0], [2.0, 2.0], [3.0, 3.0]]'",
		"'{\"1\":1,\"2\":2}'",
		"'{\"1\":[true, false]}'",
		"DATE('2024-01-01 00:00:00')",
	}

	insertSQLBuilder := strings.Builder{}
	insertSQLBuilder.WriteString("INSERT INTO test_struct(b, c, d, e, f, g) ")
	insertSQLBuilder.WriteString("VALUES(")
	for i, value := range values {
		insertSQLBuilder.WriteString(value)
		if i != len(values) - 1 {
			insertSQLBuilder.WriteString(", ")
		}
	}
	insertSQLBuilder.WriteString(");")

	_, err = engine.Exec(insertSQLBuilder.String())

	if err != nil {
		t.Error(err)
	}

	var models []TestStruct

	err = engine.Query("SELECT * FROM test_struct;", Parse(reflect.TypeOf(TestStruct{}), engine.dialect, false), &models)

	if err != nil {
		t.Error(err)
	}

	if len(models) != 1 {
		t.Error("query test failed")
	}

	model := models[0]
	if model.A != 1 || 
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
	err = engine.Query("SELECT * FROM sqlite_master WHERE name='test_struct';", Parse(reflect.TypeOf(dialects.SQLiteTableInfo{}), engine.dialect, true), &models)

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

	values := []string{
		"'0'",
		"'[1,2,3]'",
		"'[[1.0, 1.0], [2.0, 2.0], [3.0, 3.0]]'",
		"'{\"1\":1,\"2\":2}'",
		"'{\"1\":[true, false]}'",
		"DATE('2024-01-01 00:00:00')",
	}

	insertSQLBuilder := strings.Builder{}
	insertSQLBuilder.WriteString("INSERT INTO test_struct(b, c, d, e, f, g) ")
	insertSQLBuilder.WriteString("VALUES(")
	for i, value := range values {
		insertSQLBuilder.WriteString(value)
		if i != len(values) - 1 {
			insertSQLBuilder.WriteString(", ")
		}
	}
	insertSQLBuilder.WriteString(");")

	_, err = engine.Exec(insertSQLBuilder.String())
	
	if err != nil {
		t.Error(err)
	}

	var models []TestStruct
	err = engine.Select(&models, []string{})

	if err != nil {
		t.Error(err)
	}

	if len(models) != 1 {
		t.Error("select test failed")
	}

	model := models[0]

	if model.A != 1 || 
		model.B != "0" || 
		!reflect.DeepEqual(model.C, []int64{1, 2, 3}) || 
		!reflect.DeepEqual(model.D, [][]float64{{1.0, 1.0}, {2.0, 2.0}, {3.0, 3.0}}) || 
		!reflect.DeepEqual(model.E, map[string]int64{"1": 1, "2": 2}) ||
		!reflect.DeepEqual(model.F, map[string][]bool{"1": {true, false}}) || 
		!reflect.DeepEqual(model.G, time.Date(2024, 01, 01, 00, 00, 00, 000, time.UTC)) {
		t.Error("select test failed")
	}
}

func TestSelectConditional(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Error(err)
	}

	_, err = engine.Create(TestStruct{})

	if err != nil {
		t.Error(err)
	}

	values := []string{
		"'0'",
		"'[1,2,3]'",
		"'[[1.0, 1.0], [2.0, 2.0], [3.0, 3.0]]'",
		"'{\"1\":1,\"2\":2}'",
		"'{\"1\":[true, false]}'",
		"DATE('2024-01-01 00:00:00')",
	}

	insertSQLBuilder := strings.Builder{}
	insertSQLBuilder.WriteString("INSERT INTO test_struct(b, c, d, e, f, g) ")
	insertSQLBuilder.WriteString("VALUES(")
	for i, value := range values {
		insertSQLBuilder.WriteString(value)
		if i != len(values) - 1 {
			insertSQLBuilder.WriteString(", ")
		}
	}
	insertSQLBuilder.WriteString(");")

	_, err = engine.Exec(insertSQLBuilder.String())
	
	if err != nil {
		t.Error(err)
	}

	var models []TestStruct

	err = engine.Select(&models, []string{"g > DATE('2024-01-01 00:00:00')"})

	if err != nil {
		t.Error(err)
	}

	if len(models) != 0 {
		t.Error("select test failed")
	}

	err = engine.Select(&models, []string{"g = DATE('2024-01-01 00:00:00')"})

	if err != nil {
		t.Error(err)
	}

	if len(models) != 1 {
		t.Error("select test failed")
	}

	models = make([]TestStruct, 0)

	err = engine.Select(&models, []string{"b = '0'", "g = DATE('2024-01-01 00:00:00')"})

	if err != nil {
		t.Error(err)
	}

	if len(models) != 1 {
		t.Error("select test failed")
	}

	model := models[0]
	if model.A != 1 || 
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
	engine.Select(&models, []string{})

	if len(models) != 1 {
		t.Error("insert test failed")
	}

	model := models[0]

	if model.A != 1 || 
		model.B != "0" || 
		!reflect.DeepEqual(model.C, []int64{1, 2, 3}) || 
		!reflect.DeepEqual(model.D, [][]float64{{1.0, 1.0}, {2.0, 2.0}, {3.0, 3.0}}) || 
		!reflect.DeepEqual(model.E, map[string]int64{"1": 1, "2": 2}) || 
		!reflect.DeepEqual(model.F, map[string][]bool{"1": {true, false}}) || 
		!reflect.DeepEqual(model.G, time.Date(2024, 01, 01, 00, 00, 00, 000, time.UTC)) {
		t.Error("insert test failed")
	}
}

func TestUpdate(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Error(err)
	}

	_, err = engine.Create(TestStruct{})

	if err != nil {
		t.Error(err)
	}

	err = engine.Insert(TestStruct{
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
	engine.Select(&models, []string{})

	if len(models) != 1 {
		t.Error("update test failed")
	}

	model := models[0]

	model.B = "1"
	model.C = []int64{2, 3, 4}
	model.D = [][]float64{{2.0, 2.0}, {3.0, 3.0}, {4.0, 4.0}}
	model.E["3"] = 3
	model.F["2"] = []bool{false, true}
	model.G = time.Date(2024, 12, 31, 00, 00, 00, 000, time.UTC)

	err = engine.Update(model)

	if err != nil {
		t.Error(err)
	}

	if model.A != 1 || 
		model.B != "1" || 
		!reflect.DeepEqual(model.C, []int64{2, 3, 4}) || 
		!reflect.DeepEqual(model.D, [][]float64{{2.0, 2.0}, {3.0, 3.0}, {4.0, 4.0}}) || 
		!reflect.DeepEqual(model.E, map[string]int64{"1": 1, "2": 2, "3": 3}) || 
		!reflect.DeepEqual(model.F, map[string][]bool{"1": {true, false}, "2": {false, true}}) || 
		!reflect.DeepEqual(model.G, time.Date(2024, 12, 31, 00, 00, 00, 000, time.UTC)) {
		t.Error("update test failed")
	}
}

func TestDelete(t *testing.T) {
	engine, err := New("file::memory:", "sqlite3", []interface{}{})

	if err != nil {
		t.Error(err)
	}

	_, err = engine.Create(TestStruct{})

	if err != nil {
		t.Error(err)
	}

	err = engine.Insert(TestStruct{
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
	engine.Select(&models, []string{})

	if len(models) != 1 {
		t.Error("delete test failed")
	}

	model := models[0]

	err = engine.Delete(model)

	if err != nil {
		t.Error(err)
	}

	models = make([]TestStruct, 0)
	engine.Select(&models, []string{})

	if len(models) != 0 {
		t.Error("delete test failed")
	}
}
