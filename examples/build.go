package main

import (
	"fmt"
trie	"github.com/henesy/fstrie"
)


/* Build and print a trie */
func main() {
	t := trie.New()

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

	fmt.Print(t.String())
}

