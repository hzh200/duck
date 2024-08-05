package core

import (
	"fmt"
	"reflect"
	"strings"
	"duck/kernel/persistence/core/dialects"
)

type Clause int
const (
	CREATE Clause = iota
	SELECT
	INSERT
	UPDATE
	DELETE
	WHERE
	ORDERBY
	LIMIT
	ClauseBorder
)

type ClauseFunction func(params []interface{}) string

var ClauseFunctions map[Clause]ClauseFunction

func init() {
	ClauseFunctions = make(map[Clause]ClauseFunction)
	ClauseFunctions[CREATE] = createClauseFunction
	ClauseFunctions[SELECT] = selectClauseFunction
	ClauseFunctions[INSERT] = insertClauseFunction
	ClauseFunctions[UPDATE] = updateClauseFunction
	ClauseFunctions[DELETE] = deleteClauseFunction
	ClauseFunctions[WHERE] = whereClauseFunction
}

func createClauseFunction(params []interface{}) string {
	schema := params[0].(*Schema)
	dialect := params[1].(dialects.Dialect)
	template := strings.Builder{}
	template.WriteString(fmt.Sprintf("CREATE TABLE %s (", schema.TableName))
	for i, field := range schema.Fields {
		template.WriteString(fmt.Sprintf("%s %s", field.FieldName, field.FieldType))
		if reflect.DeepEqual(field, schema.PrimaryKeyField) {
			template.WriteString(fmt.Sprintf(" %s", dialect.PrimaryKey()))
		}

		if field.AutoIncrement {
			template.WriteString(fmt.Sprintf(" %s", dialect.AutoIncrement()))
		}

		if i != len(schema.Fields) - 1 {
			template.WriteString(", ")
		}

		template.WriteString("")
	}
	template.WriteString(")")
	return template.String()
}

func selectClauseFunction(params []interface{}) string {
	schema := params[0].(*Schema)
	template := strings.Builder{}
	template.WriteString(fmt.Sprintf("SELECT * FROM %s", schema.TableName))
	return template.String()
}

func insertClauseFunction(params []interface{}) string {
	schema := params[0].(*Schema)
	columns := params[1].([]string)
	values := params[2].([]string)
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", schema.TableName, strings.Join(columns, ","), strings.Join(values, ", "))
}

func updateClauseFunction(params []interface{}) string {
	schema := params[0].(*Schema)
	pairs := params[1].([]string)
	return fmt.Sprintf("UPDATE %s SET %s", schema.TableName, strings.Join(pairs, ", "))
}

func deleteClauseFunction(params []interface{}) string {
	schema := params[0].(*Schema)
	return fmt.Sprintf("DELETE FROM %s", schema.TableName)
}

func whereClauseFunction(params []interface{}) string {
	_ = params[0].(*Schema)
	conditions := params[1].([]string)
	if len(conditions) == 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " and "))
}
