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
