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
