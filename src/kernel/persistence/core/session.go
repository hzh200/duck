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
			clauseStr := ClauseFunctions[clause](data)
			if len(clauseStr) > 0 && session.sqlBuilder.Len() > 0 && clauseStr[0] != 32 && session.sqlBuilder.String()[session.sqlBuilder.Len() - 1] != 32 {
				clauseStr = string(append([]byte{' '}, []byte(clauseStr)...))
			}
			session.sqlBuilder.WriteString(clauseStr)
		}
	}
	session.sqlBuilder.WriteString(";")
}
