package core

import "reflect"

func (engine *Engine) Exec(sql string) error {
	// To avoid SQL injections.
	statement, err := engine.db.Prepare(sql)
	
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) Query(sql string, schema *Schema) (interface{}, error) {
	res, err := engine.db.Query(sql)
	
	if err != nil {
		return nil, err
	}

	defer res.Close()

	models := make([]interface{}, 0)
	// Iterate and fetch the records from result cursor
	for res.Next() {
		fields := make([]interface{}, len(schema.Fields))
		for i, field := range schema.Fields {
			fields[i] = reflect.New(field.MemberType)
		}
		res.Scan(fields...)
		item := reflect.New(schema.ModelType)
		for i := 0; i < item.NumField(); i++ {
			item.Field(i).Set(reflect.ValueOf(fields[i]))
		}
		models = append(models, item)
	}

	return reflect.ValueOf(models).Interface(), nil
}

func (engine *Engine) Create(model interface{}) error {
	session := engine.NewSession(model)
	session.AddClause(CREATE, []interface{}{session.schema})
	session.Build()
	err := engine.Exec(session.sqlBuilder.String())
	return err
}

func (engine *Engine) Select(model interface{}) (interface{}, error) {
	session := engine.NewSession(model)
	session.AddClause(SELECT, []interface{}{session.schema})
	session.Build()
	models, err := engine.Query(session.sqlBuilder.String(), session.schema)

	if err != nil {
		return nil, err
	}
	
	return models, nil
}

func (engine *Engine) Insert(model interface{}) error {
	session := engine.NewSession(model)
	valueMap := make(map[reflect.Type]interface{})
	values := reflect.ValueOf(model)
	for i, field := range session.schema.Fields {
		valueMap[field.MemberType] = engine.dialect.DataFormatting(values.Field(i))
	}
	session.AddClause(INSERT, []interface{}{session.schema, valueMap})
	session.Build()
	err := engine.Exec(session.sqlBuilder.String())
	return err
}
