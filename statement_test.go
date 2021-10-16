package xorm

import (
	"testing"
)

func TestStatement_Limit(t *testing.T) {
	s := &Statement{}
	s = s.Limit(1, 2, 3)
	if s.LimitStr != 1 {
		t.Fatalf("test TestStatement_Limit, unexpected error: %v != %v", s.LimitStr, 1)
	}
	if s.OffsetStr != 2 {
		t.Fatalf("test TestStatement_Limit, unexpected error: %v != %v", s.OffsetStr, 2)
	}

	s = &Statement{}
	s = s.Limit(1)
	if s.LimitStr != 1 {
		t.Fatalf("test TestStatement_Limit, unexpected error: %v != %v", s.LimitStr, 1)
	}
	if s.OffsetStr != 0 {
		t.Fatalf("test TestStatement_Limit, unexpected error: %v != %v", s.OffsetStr, 0)
	}
}

func TestStatement_Offset(t *testing.T) {
	s := &Statement{}
	if s.OffsetStr != 0 {
		t.Fatalf("test TestStatement_Offset, unexpected error: %v != %v", s.OffsetStr, 0)
	}

	s = s.Offset(10)
	if s.OffsetStr != 10 {
		t.Fatalf("test TestStatement_Offset, unexpected error: %v != %v", s.LimitStr, 10)
	}
}

func TestStatement_OrderBy(t *testing.T) {
	s := &Statement{}
	if s.OrderStr != "" {
		t.Fatalf("test TestStatement_OrderBy, unexpected error: %v != %v", s.OrderStr, "")
	}

	s = s.OrderBy("ab")
	if s.OrderStr != "ab" {
		t.Fatalf("test TestStatement_OrderBy, unexpected error: %v != %v", s.OrderStr, "ab")
	}
}

func TestStatement_genSelectSql(t *testing.T) {
	table := &Table{
		Columns: map[string]Column{
			"name": Column{Name: "name", IsPrimaryKey: true},
			"age":  Column{Name: "age"},
		},
		PrimaryKey: "name",
		Name:       "student",
	}

	var tests = []struct {
		desc string
		s    *Statement
		want string
	}{
		{
			desc: "test mssql offset and limit",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:     table,
				OffsetStr: 10,
				LimitStr:  100,
			},
			want: "select col-a from (select ROW_NUMBER() OVER(order by name )as rownum,col-a from student) as a where rownum between 10 and 100",
		},
		{
			desc: "test mssql offset and limit and where",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:     table,
				OffsetStr: 10,
				LimitStr:  100,
				WhereStr:  "a == b",
			},
			want: "select col-a from (select ROW_NUMBER() OVER(order by name )as rownum,col-a from student WHERE a == b) as a where rownum between 10 and 100",
		},
		{
			desc: "test mssql no offset has limit",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:    table,
				LimitStr: 100,
			},
			want: "SELECT top 100 col-a FROM student",
		},
		{
			desc: "test mssql no offset has limit, where",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:    table,
				LimitStr: 100,
				WhereStr: "a == b",
			},
			want: "SELECT top 100 col-a FROM student WHERE a == b",
		},
		{
			desc: "test mssql no offset has limit, where, group by",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:      table,
				LimitStr:   100,
				WhereStr:   "a == b",
				GroupByStr: "group by c",
			},
			want: "SELECT top 100 col-a FROM student WHERE a == b group by c",
		},
		{
			desc: "test mssql no offset has limit, where, group by, having",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:      table,
				LimitStr:   100,
				WhereStr:   "a == b",
				GroupByStr: "group by c",
				HavingStr:  "having d = e",
			},
			want: "SELECT top 100 col-a FROM student WHERE a == b group by c having d = e",
		},
		{
			desc: "test mssql no offset has limit, where, group by, having, order by",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:      table,
				LimitStr:   100,
				WhereStr:   "a == b",
				GroupByStr: "group by c",
				HavingStr:  "having d = e",
				OrderStr:   "f",
			},
			want: "SELECT top 100 col-a FROM student WHERE a == b group by c having d = e ORDER BY f",
		},
		{
			desc: "test mssql no offset no limit",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table: table,
			},
			want: "SELECT col-a FROM student",
		},
		{
			desc: "test mssql no offset no limit, where",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:    table,
				WhereStr: "a == b",
			},
			want: "SELECT col-a FROM student WHERE a == b",
		},
		{
			desc: "test mssql no offset no limit, where, groupby",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:      table,
				WhereStr:   "a == b",
				GroupByStr: "group by c",
			},
			want: "SELECT col-a FROM student WHERE a == b group by c",
		},
		{
			desc: "test mssql no offset no limit, where, groupby, having",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:      table,
				WhereStr:   "a == b",
				GroupByStr: "group by c",
				HavingStr:  "having d = e",
			},
			want: "SELECT col-a FROM student WHERE a == b group by c having d = e",
		},
		{
			desc: "test mssql no offset no limit, where, groupby, having, order",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "mssql"},
				},
				Table:      table,
				WhereStr:   "a == b",
				GroupByStr: "group by c",
				HavingStr:  "having d = e",
				OrderStr:   "f",
			},
			want: "SELECT col-a FROM student WHERE a == b group by c having d = e ORDER BY f",
		},
		{
			desc: "test non-mssql",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "sql"},
				},
				Table: table,
			},
			want: "SELECT col-a FROM student",
		},
		{
			desc: "test non-mssql join",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "sql"},
				},
				Table:   table,
				JoinStr: "join b",
			},
			want: "SELECT col-a FROM student join b",
		},
		{
			desc: "test non-mssql join where groupby",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "sql"},
				},
				Table:      table,
				JoinStr:    "join b",
				WhereStr:   "c == d",
				GroupByStr: "group by e",
			},
			want: "SELECT col-a FROM student join b WHERE c == d group by e",
		},
		{
			desc: "test non-mssql join where groupby having",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "sql"},
				},
				Table:      table,
				JoinStr:    "join b",
				WhereStr:   "c == d",
				GroupByStr: "group by e",
				HavingStr:  "having f = g",
			},
			want: "SELECT col-a FROM student join b WHERE c == d group by e having f = g",
		},
		{
			desc: "test non-mssql join where groupby having order",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "sql"},
				},
				Table:      table,
				JoinStr:    "join b",
				WhereStr:   "c == d",
				GroupByStr: "group by e",
				HavingStr:  "having f = g",
				OrderStr:   "j",
			},
			want: "SELECT col-a FROM student join b WHERE c == d group by e having f = g ORDER BY j",
		},
		{
			desc: "test non-mssql join where groupby having order offset",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "sql"},
				},
				Table:      table,
				JoinStr:    "join b",
				WhereStr:   "c == d",
				GroupByStr: "group by e",
				HavingStr:  "having f = g",
				OrderStr:   "j",
				OffsetStr:  10,
				LimitStr:   100,
			},
			want: "SELECT col-a FROM student join b WHERE c == d group by e having f = g ORDER BY j LIMIT 10, 100",
		},
		{
			desc: "test non-mssql join where groupby having order limit without offset",
			s: &Statement{
				Session: &Session{
					Engine: &Engine{Protocol: "sql"},
				},
				Table:      table,
				JoinStr:    "join b",
				WhereStr:   "c == d",
				GroupByStr: "group by e",
				HavingStr:  "having f = g",
				OrderStr:   "j",
				LimitStr:   100,
			},
			want: "SELECT col-a FROM student join b WHERE c == d group by e having f = g ORDER BY j LIMIT 100",
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.s.genSelectSql("col-a")
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}

}
