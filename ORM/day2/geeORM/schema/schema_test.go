package schema

import (
	"testing"

	dialetc "github.com/ChenMiaoQiu/7days-golang-learn/ORM/day2/geeORM/dialect"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialetc.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse user struct")
	}

	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse name tag")
	}
}
