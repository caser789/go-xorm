package xorm

import (
	"database/sql"
	"fmt"
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
