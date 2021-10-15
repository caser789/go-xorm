package xorm

import (
	"reflect"
	"strings"
)

type SQLType struct {
	Name          string
	DefaultLength int
}

var (
	Int     = SQLType{"int", 11}
	Char    = SQLType{"char", 1}
	Varchar = SQLType{"varchar", 50}
	Date    = SQLType{"date", 24}
	Decimal = SQLType{"decimal", 26}
	Float   = SQLType{"float", 31}
	Double  = SQLType{"double", 31}
)

const (
	PQSQL   = "pqsql"
	MSSQL   = "mssql"
	SQLITE  = "sqlite"
	MYSQL   = "mysql"
	MYMYSQL = "mymysql"
)

type Column struct {
	Name          string
	FieldName     string
	SQLType       SQLType
	Length        int
	Nullable      bool
	Default       string
	IsUnique      bool
	IsPrimaryKey  bool
	AutoIncrement bool
}

type Table struct {
	Name       string
	Type       reflect.Type
	Columns    map[string]Column
	PrimaryKey string
}

func (table *Table) ColumnStr() string {
	colNames := make([]string, 0)
	for _, col := range table.Columns {
		if col.Name == "" {
			continue
		}
		colNames = append(colNames, col.Name)
	}
	return strings.Join(colNames, ", ")
}

func (table *Table) PlaceHolders() string {
	colNames := make([]string, 0)
	for _, col := range table.Columns {
		if col.Name == "" {
			continue
		}
		colNames = append(colNames, "?")
	}
	return strings.Join(colNames, ", ")
}
