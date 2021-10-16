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
