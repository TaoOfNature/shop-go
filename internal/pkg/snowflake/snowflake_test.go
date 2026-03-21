package snowflake

import "testing"

func TestNextIDReturnsUniqueIncreasingValues(t *testing.T) {
	gen := New(1)

	first := gen.NextID()
	second := gen.NextID()
	third := gen.NextID()

	if first <= 0 || second <= 0 || third <= 0 {
		t.Fatal("expected positive ids")
	}
	if !(first < second && second < third) {
		t.Fatalf("expected monotonic ids, got %d %d %d", first, second, third)
	}
}
