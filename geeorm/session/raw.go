package session

import (
	"database/sql"
	"geeorm/clause"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/schema"
	"strings"
)

type Session struct {
	db *sql.DB
	sql strings.Builder
	sqlVals []interface{}
	dialect dialect.Dialect
	refTable *schema.Schema
	clause clause.Clause
	tx *sql.Tx
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db: db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVals = nil
	s.clause = clause.Clause{}
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVals = append(s.sqlVals, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVals...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	return s.DB().QueryRow(s.sql.String(), s.sqlVals...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVals...); err != nil {
		log.Error(err)
	}
	return
}