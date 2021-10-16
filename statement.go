package xorm

import (
	"fmt"
)

type Statement struct {
	TableName  string
	Table      *Table
	Session    *Session
	LimitStr   int
	OffsetStr  int
	WhereStr   string
	ParamStr   []interface{}
	OrderStr   string
	JoinStr    string
	GroupByStr string
	HavingStr  string
}

// Limit sets LimitStr and OffsetStr
func (statement *Statement) Limit(start int, size ...int) *Statement {
	statement.LimitStr = start
	if len(size) > 0 {
		statement.OffsetStr = size[0]
	}
	return statement
}

// Offset sets OffsetStr
func (statement *Statement) Offset(offset int) *Statement {
	statement.OffsetStr = offset
	return statement
}

// OrderBy sets OrderStr
func (statement *Statement) OrderBy(order string) *Statement {
	statement.OrderStr = order
	return statement
}

func (statement Statement) genSelectSql(columnStr string) (a string) {
	session := statement.Session
	if session.Engine.Protocol == "mssql" {
		if statement.OffsetStr > 0 {
			a = fmt.Sprintf("select ROW_NUMBER() OVER(order by %v )as rownum,%v from %v",
				statement.Table.PKColumn().Name,
				columnStr,
				statement.Table.Name)
			if statement.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, statement.WhereStr)
			}
			a = fmt.Sprintf("select %v from (%v) "+
				"as a where rownum between %v and %v",
				columnStr,
				a,
				statement.OffsetStr,
				statement.LimitStr)
		} else if statement.LimitStr > 0 {
			a = fmt.Sprintf("SELECT top %v %v FROM %v", statement.LimitStr, columnStr, statement.Table.Name)
			if statement.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, statement.WhereStr)
			}
			if statement.GroupByStr != "" {
				a = fmt.Sprintf("%v %v", a, statement.GroupByStr)
			}
			if statement.HavingStr != "" {
				a = fmt.Sprintf("%v %v", a, statement.HavingStr)
			}
			if statement.OrderStr != "" {
				a = fmt.Sprintf("%v ORDER BY %v", a, statement.OrderStr)
			}
		} else {
			a = fmt.Sprintf("SELECT %v FROM %v", columnStr, statement.Table.Name)
			if statement.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, statement.WhereStr)
			}
			if statement.GroupByStr != "" {
				a = fmt.Sprintf("%v %v", a, statement.GroupByStr)
			}
			if statement.HavingStr != "" {
				a = fmt.Sprintf("%v %v", a, statement.HavingStr)
			}
			if statement.OrderStr != "" {
				a = fmt.Sprintf("%v ORDER BY %v", a, statement.OrderStr)
			}
		}
	}
	return
}
