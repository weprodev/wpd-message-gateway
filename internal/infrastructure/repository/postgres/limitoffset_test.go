package postgres

import "testing"

func TestClampLimitOffset(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
	}{
		{"defaults", 0, 0, defaultListLimit, 0},
		{"caps max", 1000, 0, maxListLimit, 0},
		{"negative offset", 10, -5, 10, 0},
		{"valid", 25, 50, 25, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := clampLimitOffset(tt.limit, tt.offset)
			if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset {
				t.Fatalf("got (%d,%d), want (%d,%d)", gotLimit, gotOffset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}

func TestLimitOffsetSQL(t *testing.T) {
	if got := limitOffsetSQL(25, 10); got != "LIMIT 25 OFFSET 10" {
		t.Fatalf("got %q", got)
	}
}
