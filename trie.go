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

// Add a new node ;; Returns added node
func (t *Trie) Add(keyPath string, data interface{}) *Node {
	if keyPath == "/" {
		// Root is eternal
		return nil
	}
	path, key := getPath(keyPath)

	// Returns parent if !ok
	var parentPath string
	if len(path) < 2 {
		parentPath = "/"
	} else {
		parentPath = mkPath(path[:len(path)-1])
	}
	p, ok := t.Find(parentPath)
	if !ok {
		// No parent to add from
		return nil
	}
	
	// Insert after p
	n := new(Node)
	n.Next = p.Down
	n.Down = nil
	n.Key = key
	n.Data = data
	p.Down = n
	return n
}

// Remove a node (by Find) ;; Returns true if removed
func (t *Trie) Remove(keyPath string) interface{} {
	if keyPath == "/" {
		// Not allowed to delete root
		return nil
	}
	path, key := getPath(keyPath)

	p, _ := t.getParent(&path, &key)
	if p == nil {
		// Parent does not exist, no such file
		return nil
	}
	
	// When we delete a file, if it's a dir, its children are lost; preserve continuity	
	cursor := p.Down.Next
	last := p.Down
	
	if last.Key == key {
		// First child node
		tmp := p.Down
		p.Down = p.Down.Next
		return tmp.Data
	}
	
	for {
		if cursor == nil {
			// This should never happen
			return nil
		}
		if cursor.Key == key {
			// Remove this node
			tmp := cursor
			last.Next = cursor.Next
			return tmp.Data
		}
		last = cursor
		cursor = cursor.Next
	}

	return nil
}

// Find a node by string index ;; Returns {found node, true} or {last found node, false}
func (t *Trie) Find(keyPath string) (*Node, bool) {
	if keyPath == "/" {
		return t.Root, true
	}
	path, key := getPath(keyPath)
	if t.Root == nil {
		return nil, false
	}
	cursor := t.Root
	last := cursor

	for i := 0; i < len(path)+1; i++ {
		if cursor.Key == key && i == len(path) {
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

// Return the set of top-level children for a Node
func (n *Node) Children() []*Node {
	var children []*Node
	cursor := n.Down
	for {
		if cursor == nil {
			break
		}
		
		children = append(children, cursor)

		cursor = cursor.Next
	}
	return children
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

// Utility to extract path and key from keyPath
func getPath(keyPath string) (path []string, key string) {
	// Due to how split works, if we make a path with a leading /, strip 0th
	path = strings.Split(keyPath, "/")[1:]
	key = path[len(path)-1]
	return
}

// Utility to generate path strings -- slow
func mkPath(path []string) string {
	keyPath := ""
	for _, v := range path {
		keyPath += "/" + v
	}
	return keyPath
}
