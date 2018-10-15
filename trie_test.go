package fstrie

import (
	"testing"
)

func TestNew(test *testing.T) {
	t := New()
	if t.Root.Key != "/" || t.Root.Next != nil || t.Root.Down != nil || t.Root.Data != nil {
		test.Errorf("Expected root, got: %v", t.Root)
	}
}


