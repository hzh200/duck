package core

import (
	"strings"
)

type Session struct {
	engine *Engine
	schema *Schema
	clauses map[Clause][]interface{}
	sqlBuilder strings.Builder
}

func (session *Session) AddClause(clause Clause, data []interface{}) {
	session.clauses[clause] = data
}

func (session *Session) Build() {
	for clause := CREATE; clause < ClauseBorder; clause++ {
		if data, ok := session.clauses[clause]; ok {
			session.sqlBuilder.WriteString(ClauseFunctions[clause](data))
		}
	}
}
