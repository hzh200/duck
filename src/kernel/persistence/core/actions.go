package core

import (
	"database/sql"
	"duck/kernel/log"
	"fmt"
	"reflect"
)

func (engine *Engine) Exec(sql string) (sql.Result, error) {
	log.Info(sql)
	// To avoid SQL injections.
	statement, err := engine.db.Prepare(sql)

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	res, err := statement.Exec()

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (engine *Engine) Query(sql string, schema *Schema, models interface{}) error {
	log.Info(sql)
	modelSlice := reflect.Indirect(reflect.ValueOf(models))

	res, err := engine.db.Query(sql)

	if err != nil {
		return err
	}

	defer res.Close()

	// Iterate and fetch the records from result cursor
	for res.Next() {
		item := reflect.New(schema.ModelType).Elem()

		cols, err := res.ColumnTypes()

		if err != nil {
			return err
		}

		fieldNum := len(schema.Fields)
		scanFields := make([]interface{}, fieldNum)

		for i, col := range cols {
			if col.DatabaseTypeName() == "TEXT" {
				scanFields[i] = reflect.New(reflect.TypeFor[string]()).Elem().Addr().Interface()
			} else if col.DatabaseTypeName() == "BLOB" {
				scanFields[i] = reflect.New(reflect.TypeFor[[]byte]()).Elem().Addr().Interface()
			} else {
				scanFields[i] = reflect.New(item.FieldByName(schema.Fields[i].MemberName).Type()).Elem().Addr().Interface()
			}
		}

		err = res.Scan(scanFields...)

		if err != nil {
			return err
		}

		for i := 0; i < fieldNum; i++ {
			val := reflect.Indirect(reflect.ValueOf(scanFields[i])).Interface()
			field := item.FieldByName(schema.Fields[i].MemberName)
			field.Set(reflect.ValueOf(engine.dialect.DataDeformatting(val, field.Type())))
		}

		modelSlice.Set(reflect.Append(modelSlice, item))
	}

	return nil
}

func (engine *Engine) Create(model interface{}) (sql.Result, error) {
	session := engine.NewSession(reflect.TypeOf(model))
	session.AddClause(CREATE, []interface{}{session.schema})
	session.Build()
	res, err := engine.Exec(session.sqlBuilder.String())
	return res, err
}

func (engine *Engine) Select(models interface{}) error {
	// models argument should be address of a struct slice
	modelSlice := reflect.Indirect(reflect.ValueOf(models))
	modelType := modelSlice.Type().Elem()
	session := engine.NewSession(modelType)
	session.AddClause(SELECT, []interface{}{session.schema})
	session.Build()
	err := engine.Query(session.sqlBuilder.String(), session.schema, models)
	return err
}

func (engine *Engine) Insert(model interface{}) error {
	session := engine.NewSession(reflect.TypeOf(model))
	values := make([]string, 0)
	for i, _ := range session.schema.Fields {
		values = append(values, fmt.Sprintf("%v", engine.dialect.DataFormatting(reflect.ValueOf(model).Field(i).Interface())))
	}
	session.AddClause(INSERT, []interface{}{session.schema, values})
	session.Build()
	_, err := engine.Exec(session.sqlBuilder.String())
	return err
}
