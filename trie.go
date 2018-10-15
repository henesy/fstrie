package fstrie

import (
	"strings"
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


// Create a new trie
func New() (t Trie) {
	n := new(Node)
	n.Next = nil
	n.Down = nil
	n.Data = nil
	n.Key = "/"
	t.Root = n
	return
}

// Add a new node, all controls whether `mkdir -p` behaviour ;; Returns added node
func (t *Trie) Add(keyPath string, data interface{}, all bool) *Node {
	if keyPath == "/" {
		// Root is eternal
		return nil
	}
	//path, key := mkPath(keyPath)
	_, key := mkPath(keyPath)
	p, ok := t.Find(keyPath)
	if !ok {
		if all {
			// Make all files in path if don't exist
			
		} else {
			// No parent to add from
			return nil
		}
	}
	
	// Insert after p
	// TODO -- this does not account for Down… this is wrong.
	tmp := p.Next
	n := new(Node)
	n.Next = tmp
	n.Down = nil
	n.Key = key
	n.Data = data
	p.Next = n
	return n
}

// Remove a node (by Find) ;; Returns true if removed
func (t *Trie) Remove(keyPath string) bool {
	if keyPath == "/" {
		// Not allowed to delete root
		return false
	}
	path, key := mkPath(keyPath)
	// Get the parent
	p, c := t.getParent(&path, &key)
	if p == nil {
		return false
	}
	
	// When we delete a file, if it's a dir, its children are lost; preserve continuity
	// TODO -- this doesn't account for Down... ; this is wrong
	tmp := p.Next
	p.Next = c.Next
	c.Next = tmp

	return true
}

// Find a node by string index ;; Returns {found node, true} or {last found node, false}
func (t *Trie) Find(keyPath string) (*Node, bool) {
	if keyPath == "/" {
		return t.Root, true
	}
	path, key := mkPath(keyPath)
	if t.Root == nil {
		return nil, false
	}
	cursor := t.Root
	last := cursor

	for i := 0; i < len(path); i++ {
		if cursor.Key == key && i == len(path)-1 {
			// Last element is our file, if exists
			return cursor, true
		}
		last = cursor
		cursor = cursor.getChild(path[i])
		if cursor == nil {
			// There are no more paths to walk
			break
		}
	}
	return last, false
}

// Return a string version of the trie -- du -a
func (t *Trie) String() {
	
}

// Unexported

// Walks until it finds the parent of an element ;; Returns {parent, child}
func (t *Trie) getParent(path *[]string, key *string) (*Node, *Node) {
	p, ok := t.Find((*path)[len(*path)-2])
	if !ok {
		// If there's no parent, there's no parent
		return nil, nil
	}

	c := p.getChild(*key)
	if c == nil {
		// If there's no child named key, there's no parent
		return nil, nil
	}

	return p, c
}

// Does a down → (next…) while able and matches a child with Key key
func (n *Node) getChild(key string) *Node {
	if n.Down == nil {
		return nil
	}
	cursor := n.Down
	for {
		if cursor.Key == key {
			return cursor
		}
		if cursor.Next == nil {
			return nil
		}
		cursor = cursor.Next
	}
}

// Utility for code reuse
func mkPath(keyPath string) (path []string, key string) {
	// Due to how split works, if we make a path with a leading /, strip 0th
	path = strings.Split(keyPath, "/")[1:]
	key = path[len(path)-1]
	return
}
