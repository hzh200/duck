package core

import (
	"database/sql"
	"duck/kernel/persistence/core/dialects"
	"errors"
	"os"
	"reflect"
)

type Engine struct {
	db *sql.DB
	dialect dialects.Dialect
}

func New(dbPath string, driver string, models []interface{}) (*Engine, error) {
	var db *sql.DB
	var err error

	// Check if the database file exists yet.
	_, existence := os.Stat(dbPath)

	// No need to create file here, creating file here would also cause some inconvinence in using in-memory databases
	// var file *os.File
	// if os.IsExist(existence) {
	// 	file, err = os.Open(dbPath)
	// } else {
	// 	file, err = os.Create(dbPath)
	// }

	// if err != nil {
	// 	return nil, err
	// }

	// err = file.Close()

	// if err != nil {
	// 	return nil, err
	// }

	// Open the database.
	db, err = sql.Open(driver, dbPath)

	if err != nil {
		return nil, err
	}

	dialect := dialects.GetDialect(driver)
	
	if dialect == nil {
		return nil, errors.New("target dialect doesn't exist")
	}

	engine := Engine{}
	engine.db = db
	engine.dialect = dialect

	// Init database tables.
	if os.IsNotExist(existence) {
		for _, model := range models {
			_, err = engine.Create(model)
			if err != nil {
				return nil, err
			}
		}
	}
	
	return &engine, nil
}

func(engine *Engine) NewSession(modelType reflect.Type) *Session {
	session := Session{}
	session.engine = engine
	session.schema = Parse(modelType, engine.dialect)
	session.clauses = make(map[Clause][]interface{})
	return &session
}
