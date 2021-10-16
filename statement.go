package xorm

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

func (statement *Statement) Limit(start int, size ...int) *Statement {
	statement.LimitStr = start
	if len(size) > 0 {
		statement.OffsetStr = size[0]
	}
	return statement
}
