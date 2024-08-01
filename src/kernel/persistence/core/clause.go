package core

import (
	"fmt"
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
	values := params[1].([]string)
	template := strings.Builder{}
	columns := make([]string, 0)
	for _, field := range schema.Fields {
		columns = append(columns, field.FieldName)
	}
	template.WriteString(fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", schema.TableName, strings.Join(columns, ","), strings.Join(values, ", ")))
	return template.String()
}
