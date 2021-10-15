package xorm

import (
	"testing"
)

func TestSnakeMapper_Table2Obj(t *testing.T) {
	var tests = []struct {
		desc  string
		s     IMapper
		input string
		want  string
	}{
		{
			desc:  "a_b -> AB",
			s:     &SnakeMapper{},
			input: "a_b",
			want:  "AB",
		},
		{
			desc:  "ab_cd_ -> AbCd",
			s:     &SnakeMapper{},
			input: "ab_cd_",
			want:  "AbCd",
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.s.Table2Obj(tt.input)
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}
}

func TestSnakeMapper_Obj2Table(t *testing.T) {
	var tests = []struct {
		desc  string
		s     IMapper
		input string
		want  string
	}{
		{
			desc:  "AB -> a_b",
			s:     &SnakeMapper{},
			input: "AB",
			want:  "a_b",
		},
		{
			desc:  "AbCd -> ab_cd",
			s:     &SnakeMapper{},
			input: "AbCd",
			want:  "ab_cd",
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.s.Obj2Table(tt.input)
			if got != tt.want {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v", i, tt.desc, tt.want, got)
			}
		})
	}
}
