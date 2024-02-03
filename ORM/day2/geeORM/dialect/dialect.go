package dialect

import (
	"reflect"
)

var dialectMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialetc Dialect) {
	dialectMap[name] = dialetc
}

func GetDialect(name string) (dialetc Dialect, ok bool) {
	dialetc, ok = dialectMap[name]
	return
}
