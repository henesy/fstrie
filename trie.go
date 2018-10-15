package fstrie

import (
)


// Node for trie
type Node struct {
	Key		string
	Data	interface{}
	Next	*Node
	Down	*Node
}

// Tree for trie
type Trie struct {
	Root	*Node
}


/* Create a new trie */
func New() (t Trie) {
	t.Root = nil
	return
}

