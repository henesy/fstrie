package fstrie

import (
	"testing"
)

func TestNew(test *testing.T) {
	t := New()
	if t.Root != nil {
		test.Errorf("Expected nil, got: %v", t.Root)
	}
}


