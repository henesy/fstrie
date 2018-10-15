package fstrie

import (
	"testing"
	"fmt"
)


// Test trie creation/initialization
func TestNew(test *testing.T) {
	t := New()
	if t.Root.Key != "/" || t.Root.Next != nil || t.Root.Down != nil || t.Root.Data != nil {
		test.Errorf("Expected root, got: %v", t.Root)
	}
}

// Test single addition
func TestAdd(test *testing.T) {
	t := New()
	num := 5
	usr := t.Add("/usr", num)
	if t.Root.Down.Key != "usr" {
		test.Errorf("Expected usr as down of root, got: %v is %v", t.Root.Down, usr)
	}
}

// Test several additions
func TestAddMany(test *testing.T) {
	t := New()
	a, b, c := 3, 2, 1
	num := t.Add("/num", 4)
	t.Add("/num/a", a)
	t.Add("/num/b", b)
	t.Add("/num/c", c)

	cs := num.Children()

	if len(cs) != 3 || cs[0].Key != "c" || cs[1].Key != "b" || cs[2].Key != "a" {
		test.Errorf("Expected 3 children as a, b, c, got: %v", cs)
	}
}

// Test single removal
func TestRemove(test *testing.T) {
	t := New()
	num := 7
	t.Add("/usr", num)
	numBack := t.Remove("/usr")
	if t.Root.Down != nil || numBack != num {
		test.Errorf("Expected nil root child and val back, got: %v and %v", t.Root, numBack)
	}
}

// Test several removals
func TestRemoveMany(test *testing.T) {
	t := New()
	a, b, c := 3, 2, 1
	num := t.Add("/num", 4)
	t.Add("/num/a", a)
	t.Add("/num/b", b)
	t.Add("/num/c", c)

	t.Remove("/num/a")
	t.Remove("/num/b")

	cs := num.Children()
	if len(cs) != 1 || cs[0].Key != "c" {
		test.Errorf("Expected only c remaining, got: %v", cs)
	}
}

// Test single Find
func TestFind(test *testing.T) {
	t := New()
	t.Add("/tmp", 4)
	t.Add("/tmp/dir", 8)
	a, b := 9, 21
	t.Add("/tmp/dir/a", a)
	t.Add("/tmp/dir/b", b)
	
	aN := t.Find("/tmp/dir/a")
	bN := t.Find("/tmp/dir/b")
	cN := t.Find("/tmp/dir/c")
	
	if cN != nil || aN.Data != a || bN.Data != b {
		test.Errorf("Expected to find a, b and not c, got: %v, %v, %v", aN, bN, cN)
	}
}

// Test existent path
func TestExistent(test *testing.T) {
	t := New()
	t.Add("/tmp", 4)
	t.Add("/tmp/dir", 8)
	a, b := 9, 21
	t.Add("/tmp/dir/a", a)
	t.Add("/tmp/dir/b", b)

	cPath := t.Existent("/tmp/dir/c")
	aPath := t.Existent("/tmp/dir/a")
	
	if cPath != "/tmp/dir" || aPath != "/tmp/dir/a" {
		test.Errorf("Expected c to not exist and a to exist, got: %v, %v", cPath, aPath)
	}
}

// Test moving
func TestMv(test *testing.T) {
	t := New()
	t.Add("/tmp", 5)
	t.Add("/tmp/dir", 22)
	a, b := 24, 11
	t.Add("/tmp/dir/a", a)
	t.Add("/tmp/dir/b", b)
	
	dir2, err := t.Mv("/tmp/dir", "/tmp/dir2")
	aN := t.Find("/tmp/dir2/a")
	bN := t.Find("/tmp/dir2/b")
	dir := t.Find("/tmp/dir") // Causes infinite loop
	
	if 	t.Find("/tmp/dir2").Key != "dir2" || aN.Data != a || bN.Data != b || err != nil || dir != nil {
		test.Errorf("Expected dir to be moved with children, got: %v with %v %v from %v", dir2, aN, bN, dir)
	}
}

// Test string of tree
func TestString(test *testing.T) {
	t := New()
	
	a, b, c := 3, 2, 1
	t.Add("/tmp", 3)
	t.Add("/num", 4)
	t.Add("/num/a", a)
	t.Add("/num/b", b)
	t.Add("/num/c", c)
	t.Add("/num/c/quote", "polar")
	t.Add("/num/c/curly", "machine")
	t.Add("/num/c/curly/puppy", "woof")
	t.Add("/num/c/curly/tow", "cable")
	t.Add("/num/b/king", "sword")
	t.Add("/abc", 0x321)

	fmt.Println(t.String())
}

