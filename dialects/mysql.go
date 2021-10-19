package dialects

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/caser789/go-xorm/core"
)

func init() {
	RegisterDialect("mysql", &mysql{})
}

type mysql struct {
	Base
	net               string
	addr              string
	params            map[string]string
	loc               *time.Location
	timeout           time.Duration
	tls               *tls.Config
	allowAllFiles     bool
	allowOldPasswords bool
	clientFoundRows   bool
}

func (db *mysql) Init(uri *Uri, drivername, dataSourceName string) error {
	return db.Base.Init(db, uri, drivername, dataSourceName)
}

func (db *mysql) SqlType(c *Column) string {
	var res string
	switch t := c.SQLType.Name; t {
	case Bool:
		res = TinyInt
	case Serial:
		c.IsAutoIncrement = true
		c.IsPrimaryKey = true
		c.Nullable = false
		res = Int
	case BigSerial:
		c.IsAutoIncrement = true
		c.IsPrimaryKey = true
		c.Nullable = false
		res = BigInt
	case Bytea:
		res = Blob
	case TimeStampz:
		res = Char
		c.Length = 64
	default:
		res = t
	}

	var hasLen1 bool = (c.Length > 0)
	var hasLen2 bool = (c.Length2 > 0)
	if hasLen1 {
		res += "(" + strconv.Itoa(c.Length) + ")"
	} else if hasLen2 {
		res += "(" + strconv.Itoa(c.Length) + "," + strconv.Itoa(c.Length2) + ")"
	}
	return res
}

func (db *mysql) SupportInsertMany() bool {
	return true
}

func (db *mysql) QuoteStr() string {
	return "`"
}

func (db *mysql) SupportEngine() bool {
	return true
}

func (db *mysql) AutoIncrStr() string {
	return "AUTO_INCREMENT"
}

func (db *mysql) SupportCharset() bool {
	return true
}

func (db *mysql) IndexOnTable() bool {
	return true
}

func (db *mysql) IndexCheckSql(tableName, idxName string) (string, []interface{}) {
	args := []interface{}{db.DbName, tableName, idxName}
	sql := "SELECT `INDEX_NAME` FROM `INFORMATION_SCHEMA`.`STATISTICS`"
	sql += " WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `INDEX_NAME`=?"
	return sql, args
}

func (db *mysql) ColumnCheckSql(tableName, colName string) (string, []interface{}) {
	args := []interface{}{db.DbName, tableName, colName}
	sql := "SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `COLUMN_NAME` = ?"
	return sql, args
}

func (db *mysql) TableCheckSql(tableName string) (string, []interface{}) {
	args := []interface{}{db.DbName, tableName}
	sql := "SELECT `TABLE_NAME` from `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=? and `TABLE_NAME`=?"
	return sql, args
}

func (db *mysql) GetColumns(tableName string) ([]string, map[string]*Column, error) {
	args := []interface{}{db.DbName, tableName}
	s := "SELECT `COLUMN_NAME`, `IS_NULLABLE`, `COLUMN_DEFAULT`, `COLUMN_TYPE`," +
		" `COLUMN_KEY`, `EXTRA` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"
	cnn, err := Open(db.DriverName(), db.DataSourceName())
	if err != nil {
		return nil, nil, err
	}
	defer cnn.Close()
	rows, err := cnn.Query(s, args...)
	if err != nil {
		return nil, nil, err
	}
	cols := make(map[string]*Column)
	colSeq := make([]string, 0)
	for rows.Next() {
		col := new(Column)
		col.Indexes = make(map[string]bool)

		var columnName, isNullable, colType, colKey, extra string
		var colDefault *string
		err = rows.Scan(&columnName, &isNullable, &colDefault, &colType, &colKey, &extra)
		if err != nil {
			return nil, nil, err
		}
		col.Name = strings.Trim(columnName, "` ")
		if "YES" == isNullable {
			col.Nullable = true
		}

		if colDefault != nil {
			col.Default = *colDefault
		}

		cts := strings.Split(colType, "(")
		var len1, len2 int
		if len(cts) == 2 {
			idx := strings.Index(cts[1], ")")
			lens := strings.Split(cts[1][0:idx], ",")
			len1, err = strconv.Atoi(strings.TrimSpace(lens[0]))
			if err != nil {
				return nil, nil, err
			}
			if len(lens) == 2 {
				len2, err = strconv.Atoi(lens[1])
				if err != nil {
					return nil, nil, err
				}
			}
		}
		colName := cts[0]
		colType = strings.ToUpper(colName)
		col.Length = len1
		col.Length2 = len2
		if _, ok := SqlTypes[colType]; ok {
			col.SQLType = SQLType{colType, len1, len2}
		} else {
			return nil, nil, errors.New(fmt.Sprintf("unkonw colType %v", colType))
		}

		if colKey == "PRI" {
			col.IsPrimaryKey = true
		}
		if colKey == "UNI" {
			//col.is
		}

		if extra == "auto_increment" {
			col.IsAutoIncrement = true
		}

		if col.SQLType.IsText() {
			if col.Default != "" {
				col.Default = "'" + col.Default + "'"
			}
		}
		cols[col.Name] = col
		colSeq = append(colSeq, col.Name)
	}
	return colSeq, cols, nil
}

func (db *mysql) GetTables() ([]*Table, error) {
	args := []interface{}{db.DbName}
	s := "SELECT `TABLE_NAME`, `ENGINE`, `TABLE_ROWS`, `AUTO_INCREMENT` from `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=?"
	cnn, err := Open(db.DriverName(), db.DataSourceName())
	if err != nil {
		return nil, err
	}
	defer cnn.Close()
	rows, err := cnn.Query(s, args...)
	if err != nil {
		return nil, err
	}

	tables := make([]*Table, 0)
	for rows.Next() {
		table := NewEmptyTable()
		var name, engine, tableRows string
		var autoIncr *string
		err = rows.Scan(&name, &engine, &tableRows, &autoIncr)
		if err != nil {
			return nil, err
		}

		table.Name = name
		tables = append(tables, table)
	}
	return tables, nil
}

func (db *mysql) GetIndexes(tableName string) (map[string]*Index, error) {
	args := []interface{}{db.DbName, tableName}
	s := "SELECT `INDEX_NAME`, `NON_UNIQUE`, `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"
	cnn, err := Open(db.DriverName(), db.DataSourceName())
	if err != nil {
		return nil, err
	}
	defer cnn.Close()
	rows, err := cnn.Query(s, args...)
	if err != nil {
		return nil, err
	}

	indexes := make(map[string]*Index, 0)
	for rows.Next() {
		var indexType int
		var indexName, colName, nonUnique string
		err = rows.Scan(&indexName, &nonUnique, &colName)
		if err != nil {
			return nil, err
		}

		if indexName == "PRIMARY" {
			continue
		}

		if "YES" == nonUnique || nonUnique == "1" {
			indexType = IndexType
		} else {
			indexType = UniqueType
		}

		colName = strings.Trim(colName, "` ")

		if strings.HasPrefix(indexName, "IDX_"+tableName) || strings.HasPrefix(indexName, "UQE_"+tableName) {
			indexName = indexName[5+len(tableName) : len(indexName)]
		}

		var index *Index
		var ok bool
		if index, ok = indexes[indexName]; !ok {
			index = new(Index)
			index.Type = indexType
			index.Name = indexName
			indexes[indexName] = index
		}
		index.AddColumn(colName)
	}
	return indexes, nil
}

func (db *mysql) Filters() []Filter {
	return []Filter{&IdFilter{}}
}
