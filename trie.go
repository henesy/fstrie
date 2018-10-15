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
	n := new(Node)
	n.Next = nil
	n.Down = nil
	n.Data = nil
	n.Key = "/"
	t.Root = n
	return
}

/* Add a new node */
func (t *Trie) Add(key string, data interface{}) {
}

/* Remove a node (by Find) */
func (t *Trie) Remove(key string) {
}

/* Find a node by string index */
func (t *Trie) Find(key string) {
}

/* Return a string version of the trie */
func (t *Trie) String() {
}

