package core

import (
	"database/sql"
	"errors"
	"os"
)

type Engine struct {
	db *sql.DB
	dialect Dialect
}

func New(dbPath string, driver string, models []interface{}) (*Engine, error) {
	var db *sql.DB
	var file *os.File
	var err error

	// Create the databse file if not exists yet.
	_, existence := os.Stat(dbPath)
	if os.IsExist(existence) {
		file, err = os.Open(dbPath)
	} else {
		file, err = os.Create(dbPath)
	}

	if err != nil {
		return nil, err
	}

	err = file.Close()

	if err != nil {
		return nil, err
	}

	// Open the database.
	db, err = sql.Open(driver, dbPath)

	if err != nil {
		return nil, err
	}

	dialect := GetDialect(driver)
	
	if dialect == nil {
		return nil, errors.New("target dialect doesn't exist")
	}

	engine := Engine{}
	engine.db = db
	engine.dialect = dialect

	// Init database tables.
	if os.IsNotExist(existence) {
		for _, model := range models {
			err = engine.Create(model)
			if err != nil {
				return nil, err
			}
		}
	}
	
	return &engine, nil
}

func(engine *Engine) NewSession(model interface{}) *Session {
	session := Session{}
	session.engine = engine
	session.schema = Parse(model, engine.dialect)
	session.clauses = make(map[Clause][]interface{})
	return &session
}
