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


func TestAdd(test *testing.T) {
	t := New()
	num := 5
	usr := t.Add("/usr", num)
	if t.Root.Down.Key != "usr" {
		test.Errorf("Expected usr as down of root, got: %v is %v", t.Root.Down, usr)
	}
}

