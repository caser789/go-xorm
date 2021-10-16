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
			}},
			want: "name",
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

func TestEngine_PlaceHolders(t *testing.T) {
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
			want: "?, ?",
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.s.PlaceHolders()
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}
}

func TestEngine_PKColumn(t *testing.T) {
	var tests = []struct {
		desc string
		s    *Table
		want Column
	}{
		{
			desc: "test table primary key column",
			s: &Table{
				Columns: map[string]Column{
					"name": Column{Name: "name", IsPrimaryKey: true},
					"age":  Column{Name: "age"},
				},
				PrimaryKey: "name",
			},
			want: Column{Name: "name", IsPrimaryKey: true},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.s.PKColumn()
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}
}

func TestSQLType_genSQL(t *testing.T) {
	var tests = []struct {
		desc string
		s    SQLType
		want string
	}{
		{
			desc: "test Int",
			s:    Int,
			want: "int(111)",
		},
		{
			desc: "test Char",
			s:    Char,
			want: "char(111)",
		},
		{
			desc: "test Varchar",
			s:    Varchar,
			want: "varchar(111)",
		},
		{
			desc: "test Date",
			s:    Date,
			want: " datetime ",
		},
		{
			desc: "test Decimal",
			s:    Decimal,
			want: "decimal(111)",
		},
		{
			desc: "test Float",
			s:    Float,
			want: "float(111)",
		},
		{
			desc: "test Double",
			s:    Double,
			want: "double(111)",
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.s.genSQL(111)
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}
}

func TestEngine_genCreateSQL(t *testing.T) {
	e := Engine{
		AutoIncrement: "engine-autoinc",
	}
	var tests = []struct {
		desc string
		s    *Table
		want string
	}{
		{
			desc: "test genCreateSQL",
			s: &Table{
				Columns: map[string]Column{
					"age": Column{
						SQLType:       Int,
						Name:          "age",
						Length:        123,
						Nullable:      false,
						Default:       "345",
						IsUnique:      true,
						AutoIncrement: true,
						IsPrimaryKey:  true,
					},
				},
				PrimaryKey: "name",
				Name:       "student",
			},
			want: "CREATE TABLE IF NOT EXISTS `student` (`age` int(123)  NOT NULL PRIMARY KEY engine-autoinc Unique);",
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := e.genCreateSQL(tt.s)
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}
}

func TestEngine_genDropSQL(t *testing.T) {
	e := Engine{
		AutoIncrement: "engine-autoinc",
	}
	var tests = []struct {
		desc string
		s    *Table
		want string
	}{
		{
			desc: "test genDropSQL",
			s: &Table{
				Columns: map[string]Column{
					"name": Column{
						Name:         "name",
						IsPrimaryKey: true,
						Length:       22,
					},
					"age": Column{
						SQLType:       Int,
						Name:          "age",
						Length:        123,
						Nullable:      false,
						Default:       "345",
						IsUnique:      true,
						AutoIncrement: true,
					},
				},
				PrimaryKey: "name",
				Name:       "student",
			},
			want: "DROP TABLE IF EXISTS `student`;",
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := e.genDropSQL(tt.s)
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}
}
