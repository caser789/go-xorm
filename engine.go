package xorm

import (
	"database/sql"
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

func (table *Table) PKColumn() Column {
	return table.Columns[table.PrimaryKey]
}

type Engine struct {
	Mapper          IMapper
	Protocol        string
	UserName        string
	Password        string
	Host            string
	Port            int
	DBName          string
	Charset         string
	Others          string
	Tables          map[string]Table
	AutoIncrement   string
	ShowSQL         bool
	QuoteIdentifier string
}

func (e *Engine) OpenDB() (db *sql.DB, err error) {
	db = nil
	err = nil
	if e.Protocol == "sqlite" {
		// 'sqlite:///foo.db'
		db, err = sql.Open("sqlite3", e.Others)
		// 'sqlite:///:memory:'
	} else if e.Protocol == "mysql" {
		// 'mysql://<username>:<passwd>@<host>/<dbname>?charset=<encoding>'
		connstr := strings.Join([]string{e.UserName, ":",
			e.Password, "@tcp(", e.Host, ":3306)/", e.DBName, "?charset=", e.Charset}, "")
		db, err = sql.Open(e.Protocol, connstr)
	} else if e.Protocol == "mymysql" {
		//   DBNAME/USER/PASSWD
		connstr := strings.Join([]string{e.DBName, e.UserName, e.Password}, "/")
		db, err = sql.Open(e.Protocol, connstr)
		//   unix:SOCKPATH*DBNAME/USER/PASSWD
		//   unix:SOCKPATH,OPTIONS*DBNAME/USER/PASSWD
		//   tcp:ADDR*DBNAME/USER/PASSWD
		//   tcp:ADDR,OPTIONS*DBNAME/USER/PASSWD
	}

	return
}
