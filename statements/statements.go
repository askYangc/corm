package statements

import (
	"github.com/askYangc/corm/parse"
	"log"
	"reflect"
	"strings"
)

const ACTION_SELECT = 0
const ACTION_INSERT = 1
const ACTION_UPDATE = 2
const ACTION_INSERTORUPDATE = 3
const ACTION_DELETE = 4
const ACTION_GET = 5

type Conditions struct {

}

type Suffix struct {

}

type Statements struct {
	DoAction int
	Table *parse.SqlTable

	/*
	//如果是普通函数，Value即传入的具有地址的对象，如果是select，则value是一个新创建的局部变量reflect.Value。
	*/
	Value reflect.Value

	Columns []string		//字段
	Builder strings.Builder

	//Omit
	Omit []string			//column

	//args
	FuncArgs []interface{}

	//limit
	LimitOffset	uint32
	LimitNum	uint32
}

func (s *Statements) Reset() {
	s.DoAction = 0
	s.Table = nil
	s.LimitOffset = 0
	s.LimitNum = 0

	s.Builder.Reset()
	if len(s.Columns) != 0 {
		s.Columns = make([]string, 0)
	}

	if len(s.Omit) != 0 {
		s.Omit = make([]string, 0)
	}

	if len(s.FuncArgs) != 0 {
		s.FuncArgs = make([]interface{}, 0)
	}
}

func (s *Statements) SetLimit(offset, num uint32) {
	s.LimitOffset = offset
	s.LimitNum = num
}

func (s *Statements) isZero(val reflect.Value) bool{
	switch val.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return false
	}

	if val.IsZero() {
		return true
	}

	return false
}

func (s *Statements) PrimaryKeyIsZero() bool {
	f := s.Table.GetSqlField(s.Table.PrimaryTag)
	val := s.Value.FieldByName(f.FiledName)
	if val.IsZero() {
		return true
	}
	return false
}

//加载表名
func (s *Statements) tableNameJoin() {
	s.Builder.WriteString(s.Table.TableName)
}


//构建?形式的语句
func (s *Statements) GetSql() string {
	switch s.DoAction {
	case ACTION_INSERT:
		return s.GenerateInsertSql()
	case ACTION_UPDATE:
		return s.GenerateUpdateSql()
	case ACTION_INSERTORUPDATE:
		return s.GenerateInsertorUpdateSql()
	case ACTION_DELETE:
		return s.GenerateDeleteSql()
	case ACTION_GET:
		return s.GenerateGetSql()
	case ACTION_SELECT:
		return s.GenerateSelectSql()
	}
	return ""
}

func (s *Statements) GetArgs() (args []interface{}){
	switch s.DoAction {
	case ACTION_INSERT:
		return s.GenerateInsertArgs()
	case ACTION_UPDATE:
		return s.GenerateUpdateArgs()
	case ACTION_INSERTORUPDATE:
		return s.GenerateInsertorUpdateArgs()
	case ACTION_DELETE:
		return s.GenerateDeleteArgs()
	case ACTION_GET:
		return s.GenerateGetArgs()
	case ACTION_SELECT:
		return s.GenerateSelectArgs()
	}
	return nil
}

func (s *Statements) GetColumnsArgs(column string) interface{}{
	field := s.Table.GetSqlField(column)
	if field == nil {
		log.Panicf("column %s not found in table %s", column, s.Table.TableName)
	}
	val := s.Value.FieldByName(field.FiledName)
	return val.Interface()
}



func (s *Statements) Join() (string, []interface{}) {
	return s.GetSql(), s.GetArgs()
}

// InSlice checks given string in string slice or not.
func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}