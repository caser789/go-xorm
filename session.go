package xorm

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Session struct {
	Db              *sql.DB
	Engine          *Engine
	Statements      []Statement
	Mapper          IMapper
	AutoCommit      bool
	ParamIteration  int
	CurStatementIdx int
}

func (session *Session) Init() {
	session.Statements = make([]Statement, 0)
	session.CurStatementIdx = -1

	session.ParamIteration = 1
}

func (session *Session) Begin() {}

func (session *Session) Rollback() {}

func (session *Session) Close() {
	defer session.Db.Close()
}

func (session *Session) CurrentStatement() *Statement {
	if session.CurStatementIdx > -1 {
		return &session.Statements[session.CurStatementIdx]
	}
	return nil
}

// newStatement creates a new statement, append it to the list and direct pointer to it
func (session *Session) newStatement() {
	state := Statement{}
	state.Session = session
	session.Statements = append(session.Statements, state)
	session.CurStatementIdx = len(session.Statements) - 1
}

func (session *Session) AutoStatement() *Statement {
	if session.CurStatementIdx == -1 {
		session.newStatement()
	}
	return session.CurrentStatement()
}

//Execute sql
func (session *Session) Exec(finalQueryString string, args ...interface{}) (sql.Result, error) {
	rs, err := session.Db.Prepare(finalQueryString)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	res, err := rs.Exec(args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (session *Session) Where(querystring interface{}, args ...interface{}) *Session {
	statement := session.AutoStatement()
	switch querystring := querystring.(type) {
	case string:
		statement.WhereStr = querystring
	case int:
		if session.Engine.Protocol == "pgsql" {
			statement.WhereStr = fmt.Sprintf("%v%v%v = $%v", session.Engine.QuoteIdentifier, statement.Table.PKColumn().Name, session.Engine.QuoteIdentifier, session.ParamIteration)
		} else {
			statement.WhereStr = fmt.Sprintf("%v%v%v = ?", session.Engine.QuoteIdentifier, statement.Table.PKColumn().Name, session.Engine.QuoteIdentifier)
			session.ParamIteration++
		}
		args = append(args, querystring)
	}
	statement.ParamStr = args
	return session
}

func (session *Session) Limit(start int, size ...int) *Session {
	session.AutoStatement().LimitStr = start
	if len(size) > 0 {
		session.CurrentStatement().OffsetStr = size[0]
	}
	return session
}

func (session *Session) Offset(offset int) *Session {
	session.AutoStatement().OffsetStr = offset
	return session
}

func (session *Session) OrderBy(order string) *Session {
	session.AutoStatement().OrderStr = order
	return session
}

func (session *Session) GroupBy(keys string) *Session {
	session.AutoStatement().GroupByStr = fmt.Sprintf("GROUP BY %v", keys)
	return session
}

func (session *Session) Having(conditions string) *Session {
	session.AutoStatement().HavingStr = fmt.Sprintf("HAVING %v", conditions)
	return session
}

//The join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (session *Session) Join(join_operator, tablename, condition string) *Session {
	if session.AutoStatement().JoinStr != "" {
		session.CurrentStatement().JoinStr = session.CurrentStatement().JoinStr + fmt.Sprintf("%v JOIN %v ON %v", join_operator, tablename, condition)
	} else {
		session.CurrentStatement().JoinStr = fmt.Sprintf("%v JOIN %v ON %v", join_operator, tablename, condition)
	}

	return session
}

func (session *Session) Commit() {
	for _, statement := range session.Statements {
		sql := statement.generateSql()
		session.Exec(sql)
	}
}

func StructName(s interface{}) string {
	v := reflect.TypeOf(s)
	return Type2StructName(v)
}

func Type2StructName(v reflect.Type) string {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Name()
}
