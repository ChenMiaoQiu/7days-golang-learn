package session

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ChenMiaoQiu/7days-golang-learn/ORM/day2/geeORM/log"
	"github.com/ChenMiaoQiu/7days-golang-learn/ORM/day2/geeORM/schema"
)

func (s *Session) Model(values interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(values) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(values, s.dialetc)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.refTable
	var cloumns []string
	for _, field := range table.Fields {
		cloumns = append(cloumns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(cloumns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE %s;", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialetc.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
