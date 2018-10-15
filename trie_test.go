package fstrie

import (
	"testing"
	"fmt"
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

func TestAddMany(test *testing.T) {
	t := New()
	a, b, c := 3, 2, 1
	num := t.Add("/num", 4)
	t.Add("/num/a", a)
	t.Add("/num/b", b)
	t.Add("/num/c", c)
	
	cs := num.Children()
	
	fmt.Println(num)
	
	if len(cs) != 3 || cs[0].Key != "c" || cs[1].Key != "b" || cs[2].Key != "a" {
		test.Errorf("Expected 3 children as a, b, c, got: %v", cs)
	}
}
