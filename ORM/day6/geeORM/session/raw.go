package session

import (
	"database/sql"
	"strings"

	"github.com/ChenMiaoQiu/7days-golang-learn/ORM/day6/geeORM/clause"
	"github.com/ChenMiaoQiu/7days-golang-learn/ORM/day6/geeORM/dialect"
	"github.com/ChenMiaoQiu/7days-golang-learn/ORM/day6/geeORM/log"
	"github.com/ChenMiaoQiu/7days-golang-learn/ORM/day6/geeORM/schema"
)

type Session struct {
	db       *sql.DB
	dialetc  dialect.Dialect
	tx       *sql.Tx
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []interface{}
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialetc dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialetc: dialetc,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) Db() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()

	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.Db().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.Db().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.Db().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
