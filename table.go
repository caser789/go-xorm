package xorm

import (
	"reflect"
	"strings"
	"time"
)

type SQLType struct {
	Name           string
	DefaultLength  int
	DefaultLength2 int
}

func (s *SQLType) IsText() bool {
	return s.Name == Char || s.Name == Varchar || s.Name == TinyText ||
		s.Name == Text || s.Name == MediumText || s.Name == LongText
}

func (s *SQLType) IsBlob() bool {
	return (s.Name == TinyBlob) || (s.Name == Blob) ||
		s.Name == MediumBlob || s.Name == LongBlob ||
		s.Name == Binary || s.Name == VarBinary || s.Name == Bytea
}

var (
	Bit       = "BIT"
	TinyInt   = "TINYINT"
	SmallInt  = "SMALLINT"
	MediumInt = "MEDIUMINT"
	Int       = "INT"
	Integer   = "INTEGER"
	BigInt    = "BIGINT"

	Char       = "CHAR"
	Varchar    = "VARCHAR"
	TinyText   = "TINYTEXT"
	Text       = "TEXT"
	MediumText = "MEDIUMTEXT"
	LongText   = "LONGTEXT"
	Binary     = "BINARY"
	VarBinary  = "VARBINARY"

	Date      = "DATE"
	DateTime  = "DATETIME"
	Time      = "TIME"
	TimeStamp = "TIMESTAMP"

	Decimal = "DECIMAL"
	Numeric = "NUMERIC"

	Real   = "REAL"
	Float  = "FLOAT"
	Double = "DOUBLE"

	TinyBlob   = "TINYBLOB"
	Blob       = "BLOB"
	MediumBlob = "MEDIUMBLOB"
	LongBlob   = "LONGBLOB"
	Bytea      = "BYTEA"

	Bool = "BOOL"

	Serial    = "SERIAL"
	BigSerial = "BIGSERIAL"

	sqlTypes = map[string]bool{
		Bit:       true,
		TinyInt:   true,
		SmallInt:  true,
		MediumInt: true,
		Int:       true,
		Integer:   true,
		BigInt:    true,

		Char:       true,
		Varchar:    true,
		TinyText:   true,
		Text:       true,
		MediumText: true,
		LongText:   true,
		Binary:     true,
		VarBinary:  true,

		Date:      true,
		DateTime:  true,
		Time:      true,
		TimeStamp: true,

		Decimal: true,
		Numeric: true,

		Real:       true,
		Float:      true,
		Double:     true,
		TinyBlob:   true,
		Blob:       true,
		MediumBlob: true,
		LongBlob:   true,
		Bytea:      true,

		Bool: true,

		Serial:    true,
		BigSerial: true,
	}
)

var b byte
var tm time.Time

func Type2SQLType(t reflect.Type) (st SQLType) {
	switch k := t.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		st = SQLType{Int, 0, 0}
	case reflect.Int64, reflect.Uint64:
		st = SQLType{BigInt, 0, 0}
	case reflect.Float32:
		st = SQLType{Float, 0, 0}
	case reflect.Float64:
		st = SQLType{Double, 0, 0}
	case reflect.Complex64, reflect.Complex128:
		st = SQLType{Varchar, 64, 0}
	case reflect.Array, reflect.Slice:
		if t.Elem() == reflect.TypeOf(b) {
			st = SQLType{Blob, 0, 0}
		} else {
			st = SQLType{Text, 0, 0}
		}
	case reflect.Bool:
		st = SQLType{Bool, 0, 0}
	case reflect.String:
		st = SQLType{Varchar, 255, 0}
	case reflect.Struct:
		if t == reflect.TypeOf(tm) {
			st = SQLType{DateTime, 0, 0}
		} else {
			st = SQLType{Text, 0, 0}
		}
	default:
		st = SQLType{Text, 0, 0}
	}
	return
}

const (
	TWOSIDES = iota + 1
	ONLYTODB
	ONLYFROMDB
)

const (
	NONEINDEX = iota
	SINGLEINDEX
	UNIONINDEX
)

const (
	NONEUNIQUE = iota
	SINGLEUNIQUE
	UNIONUNIQUE
)

type Index struct {
	Name     string
	IsUnique bool
	Cols     []*Column
}

func NewIndex(name string, isUnique bool) *Index {
	return &Index{name, isUnique, make([]*Column, 0)}
}

type Column struct {
	Name            string
	FieldName       string
	SQLType         SQLType
	Length          int
	Length2         int
	Nullable        bool
	Default         string
	UniqueType      int
	UniqueName      string
	IndexType       int
	IndexName       string
	IsPrimaryKey    bool
	IsAutoIncrement bool
	MapType         int
	IsCreated       bool
	IsUpdated       bool
}

func (col *Column) String(engine *Engine) string {
	sql := engine.Quote(col.Name) + " "

	sql += engine.SqlType(col) + " "

	if col.IsPrimaryKey {
		sql += "PRIMARY KEY "
	}

	if col.IsAutoIncrement {
		sql += engine.AutoIncrStr() + " "
	}

	if col.Nullable {
		sql += "NULL "
	} else {
		sql += "NOT NULL "
	}

	if col.Default != "" {
		sql += "DEFAULT " + col.Default + " "
	}
	return sql
}

func (col *Column) ValueOf(bean interface{}) reflect.Value {
	var fieldValue reflect.Value
	if strings.Contains(col.FieldName, ".") {
		fields := strings.Split(col.FieldName, ".")
		if len(fields) > 2 {
			return reflect.ValueOf(nil)
		}

		fieldValue = reflect.Indirect(reflect.ValueOf(bean)).FieldByName(fields[0])
		fieldValue = fieldValue.FieldByName(fields[1])
	} else {
		fieldValue = reflect.Indirect(reflect.ValueOf(bean)).FieldByName(col.FieldName)
	}
	return fieldValue
}

type Table struct {
	Name       string
	Type       reflect.Type
	ColumnsSeq []string
	Columns    map[string]*Column
	Indexes    map[string][]string
	Uniques    map[string][]string
	PrimaryKey string
	Created    string
	Updated    string
	Cacher     Cacher
}

func (table *Table) PKColumn() *Column {
	return table.Columns[table.PrimaryKey]
}

func (table *Table) AddColumn(col *Column) {
	table.ColumnsSeq = append(table.ColumnsSeq, col.Name)
	table.Columns[col.Name] = col
}

func (table *Table) GenCols(session *Session, bean interface{}, useCol bool, includeQuote bool) ([]string, []interface{}, error) {
	colNames := make([]string, 0)
	args := make([]interface{}, 0)

	for _, col := range table.Columns {
		if useCol {
			if _, ok := session.Statement.columnMap[col.Name]; !ok {
				continue
			}
		}
		if col.MapType == ONLYFROMDB {
			continue
		}

		fieldValue := col.ValueOf(bean)
		if col.IsAutoIncrement && fieldValue.Int() == 0 {
			continue
		}

		if session.Statement.ColumnStr != "" {
			if _, ok := session.Statement.columnMap[col.Name]; !ok {
				continue
			}
		}

		if (col.IsCreated || col.IsUpdated) && session.Statement.UseAutoTime {
			args = append(args, time.Now())
		} else {
			arg, err := session.value2Interface(col, fieldValue)
			if err != nil {
				return colNames, args, err
			}
			args = append(args, arg)
		}

		if includeQuote {
			colNames = append(colNames, session.Engine.Quote(col.Name)+" = ?")
		} else {
			colNames = append(colNames, col.Name)
		}
	}
	return colNames, args, nil
}

type Conversion interface {
	FromDB([]byte) error
	ToDB() ([]byte, error)
}
