package strings

import (
	"testing"
)

func TestInSlice(t *testing.T) {
	list := []string{"apple", "banana", "coconut"}
	if !InSlice("apple", list) {
		t.Error("Failed to find string in list.")
	}
	if InSlice("orange", list) {
		t.Error("Incorrectly found non-member string in list.")
	}
}
