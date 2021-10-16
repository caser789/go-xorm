package xorm

import (
	"reflect"
	"testing"
)

func TestType2StructName(t *testing.T) {
	type User struct{}
	u := &User{}

	got := Type2StructName(reflect.TypeOf(u))
	want := "User"

	if got != want {
		t.Fatalf("test TestType2StructName, unexpected error: %v != %v", got, want)
	}

	a := StructName(u)
	b := "User"

	if a != b {
		t.Fatalf("test TestType2StructName, unexpected error: %v != %v", a, b)
	}
}

func TestSession_Where(t *testing.T) {
	var tests = []struct {
		desc               string
		protocol           string
		queryString        interface{}
		wantWhereStr       string
		wantParamIteration int
		wantParamStr       []interface{}
	}{
		{
			desc:               "test string querystring",
			protocol:           "pgsql",
			queryString:        "query_string",
			wantWhereStr:       "query_string",
			wantParamIteration: 1,
			wantParamStr:       []interface{}{1, 2, 3},
		},
		{
			desc:               "test int querystring pgsql",
			protocol:           "pgsql",
			queryString:        123,
			wantWhereStr:       "#name# = $1",
			wantParamIteration: 1,
			wantParamStr:       []interface{}{1, 2, 3, 123},
		},
		{
			desc:               "test int querystring non-pgsql",
			protocol:           "mysql",
			queryString:        123,
			wantWhereStr:       "#name# = ?",
			wantParamIteration: 2,
			wantParamStr:       []interface{}{1, 2, 3, 123},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			table := &Table{
				Columns: map[string]Column{
					"name": Column{Name: "name", IsPrimaryKey: true},
					"age":  Column{Name: "age"},
				},
				PrimaryKey: "name",
				Name:       "student",
			}

			session := &Session{
				Engine: &Engine{
					Protocol:        tt.protocol,
					QuoteIdentifier: "#",
				},
			}
			session.Init()
			session.Statements = append(
				session.Statements,
				Statement{
					TableName: "student",
					Table:     table,
				},
			)
			session.CurStatementIdx = 0

			session = session.Where(tt.queryString, 1, 2, 3)

			if session.AutoStatement().WhereStr != tt.wantWhereStr {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, session.AutoStatement().WhereStr, tt.wantWhereStr)
			}
			if session.ParamIteration != tt.wantParamIteration {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, session.ParamIteration, tt.wantParamIteration)
			}
			if !reflect.DeepEqual(session.AutoStatement().ParamStr, tt.wantParamStr) {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, session.AutoStatement().ParamStr, tt.wantParamStr)
			}
		})
	}
}
