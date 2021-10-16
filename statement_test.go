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
		t.Fatalf("test TestStatement_Limit, unexpected error: %v != %v", s.LimitStr, 1)
	}

	s = &Statement{}
	s = s.Limit(1)

	if s.LimitStr != 1 {
		t.Fatalf("test TestStatement_Limit, unexpected error: %v != %v", s.LimitStr, 1)
	}

	if s.OffsetStr != 0 {
		t.Fatalf("test TestStatement_Limit, unexpected error: %v != %v", s.LimitStr, 1)
	}
}
