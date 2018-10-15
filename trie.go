package fstrie

import (
)


// Node for trie
type Node struct {
	key		string
	data	interface{}
	next	*Node
	down	*Node
}

// Tree for trie
type Trie struct {
	root	*Node
}

/* Create a new trie */
func New() (t Trie) {
	t.root = nil
	return
}

