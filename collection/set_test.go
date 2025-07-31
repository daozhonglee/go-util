package collection

import "testing"

func TestNewSet(t *testing.T) {
	set := NewSet(1, 2, 3)
	if set.Len() != 3 {
		t.Errorf("Expected length 3, got %d", set.Len())
	}
}

func TestSetAdd(t *testing.T) {
	set := NewSet()
	set.Add(1, 2, 3)
	if set.Len() != 3 {
		t.Errorf("Expected length 3, got %d", set.Len())
	}
}

func TestSetRemove(t *testing.T) {
	set := NewSet(1, 2, 3)
	set.Remove(2)
	if set.Len() != 2 {
		t.Errorf("Expected length 2, got %d", set.Len())
	}
	if set.Exists(2) {
		t.Error("Element 2 should not exist after removal")
	}
}

func TestSetUnion(t *testing.T) {
	set1 := NewSet(1, 2)
	set2 := NewSet(2, 3)
	union := set1.Union(set2)
	if union.Len() != 3 {
		t.Errorf("Expected union length 3, got %d", union.Len())
	}
}

func TestSetIntersect(t *testing.T) {
	set1 := NewSet(1, 2, 3)
	set2 := NewSet(2, 3, 4)
	intersection := set1.Intersect(set2)
	if intersection.Len() != 2 {
		t.Errorf("Expected intersection length 2, got %d", intersection.Len())
	}
}

func TestSetClear(t *testing.T) {
	set := NewSet(1, 2, 3)
	set.Clear()
	if set.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", set.Len())
	}
}
