package core

import (
	"fmt"
	"reflect"
	"strings"
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
}

func createClauseFunction(params []interface{}) string {
	schema := params[0].(*Schema)
	template := strings.Builder{}
	template.WriteString(fmt.Sprintf("CREATE TABLE %s (", schema.TableName))
	for i, field := range schema.Fields {
		template.WriteString(fmt.Sprintf("%s %s", field.FieldName, field.FieldType))
		if field.Constraints != "" {
			template.WriteString(fmt.Sprintf(" %s", field.Constraints))
		}
		if i != len(schema.Fields) - 1 {
			template.WriteString(", ")
		}
		template.WriteString("")
	}
	template.WriteString(");")
	return template.String()
}

func selectClauseFunction(params []interface{}) string {
	schema := params[0].(*Schema)
	template := strings.Builder{}
	template.WriteString(fmt.Sprintf("SELECT * FROM %s;", schema.TableName))
	return template.String()
}

func insertClauseFunction(params []interface{}) string {
	schema := params[0].(*Schema)
	valueMap := params[1].(map[reflect.Type]interface{})
	template := strings.Builder{}
	insertValues := make([]string, len(schema.Fields))
	for i, field := range schema.Fields {
		if reflect.TypeOf(valueMap[field.MemberType]).Kind() == reflect.String {
			insertValues[i] = fmt.Sprintf("'%s'", valueMap[field.MemberType])
		} else {
			insertValues[i] = fmt.Sprintf("%v", valueMap[field.MemberType])
		}
	}
	insertFields := make([]string, 0)
	for _, field := range schema.Fields {
		insertFields = append(insertFields, field.FieldName)
	}
	template.WriteString(fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", schema.TableName, strings.Join(insertFields, ","), strings.Join(insertValues, ", ")))
	return template.String()
}
