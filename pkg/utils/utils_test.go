package utils

import "testing"

func TestRemoveZeroValues(t *testing.T) {
	t.Parallel()

	slice := []string{"", "a", "", "b", "c", "", "d", ""}
	expected := []string{"a", "b", "c", "d"}

	result := RemoveZeroValues(slice)

	if len(result) != len(expected) {
		t.Fatalf("Expected %v, but got %v", expected, result)
	}

	for i, value := range expected {
		if result[i] != value {
			t.Fatalf("Expected %v, but got %v", expected, result)
		}
	}
}
