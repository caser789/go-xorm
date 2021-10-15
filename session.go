package xorm

import (
	"database/sql"
)

type Session struct {
	Db             *sql.DB
	Engine         *Engine
	Mapper         IMapper
	ParamIteration int
}

func (session *Session) Init() {}
