package core

import (
	"database/sql"
	"duck/kernel/log"
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

func (engine *Engine) Query(sql string, schema *Schema, models interface{}, ) error {
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
		itemFields := make([]interface{}, len(schema.Fields))
		for i, field := range schema.Fields {
			itemFields[i] = item.FieldByName(field.MemberName).Addr().Interface()
		}
		res.Scan(itemFields...)
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
	valueMap := make(map[reflect.Type]interface{})
	values := reflect.ValueOf(model)
	for i, field := range session.schema.Fields {
		valueMap[field.MemberType] = engine.dialect.DataFormatting(values.Field(i).Interface())
	}
	session.AddClause(INSERT, []interface{}{session.schema, valueMap})
	session.Build()
	_, err := engine.Exec(session.sqlBuilder.String())
	return err
}
