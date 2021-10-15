package xorm

import (
	"testing"
)

func TestEngine_ColumnStr(t *testing.T) {
	var tests = []struct {
		desc string
		s    *Table
		want string
	}{
		{
			desc: "empty",
			s:    &Table{Columns: map[string]Column{}},
			want: "",
		},
		{
			desc: "table column in string",
			s: &Table{Columns: map[string]Column{
				"name": Column{Name: "name"},
				"age":  Column{Name: "age"},
			}},
			want: "name, age",
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.s.ColumnStr()
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}
}
