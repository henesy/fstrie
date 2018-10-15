package fstrie

import (
	"strings"
	"fmt"
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

	var parentPath string
	if len(path) < 2 {
		parentPath = "/"
	} else {
		parentPath = mkPath(path[:len(path)-1])
	}
	p := t.Find(parentPath)
	if p == nil {
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
func (t *Trie) Find(keyPath string) (*Node) {
	if keyPath == "/" {
		return t.Root
	}
	path, key := getPath(keyPath)
	if t.Root == nil {
		return nil
	}
	cursor := t.Root

	for i := 0; i < len(path)+1; i++ {
		if cursor.Key == key && i == len(path) {
			// Last element is our file, if exists
			return cursor
		}
		cursor = cursor.getChild(path[i])
		if cursor == nil {
			// There are no more paths to walk
			break
		}
	}
	return nil
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

// Find the existent portion of a path
func (t *Trie) Existent(keyPath string) string {
	if keyPath == "/" {
		return keyPath
	}
	path, _ := getPath(keyPath)
	if t.Root == nil {
		return ""
	}
	cursor := t.Root
	realPath := ""

	for i := 0; i < len(path); i++ {
		cursor = cursor.getChild(path[i])
		if cursor == nil {
			// There are no more paths to walk
			break
		}
		realPath += "/" + cursor.Key
	}
	return realPath
}

// Remove then re-Add a node into a different location ;; will overwrite existing name `to`
func (t *Trie) Mv(from, to string) (*Node, error) {
	if from == "/" || to == "/" {
		// Not allowed to change "/"
		return nil, fmt.Errorf("Root is immutable.")
	}

	tPath, tKey := getPath(to)
	n := t.Find(from)
	if n == nil {
		return nil, fmt.Errorf("Node %v not found.", from)
	}

	if from == to {
		return n, nil
	}

	// Maybe this should be its own fn
	parentPath := mkPath(tPath[:len(tPath)-1])
	if t.Find(parentPath) == nil {
		// Parent does not exist
		return nil, fmt.Errorf("Parent node %v not found.", parentPath)
	}
	
	n.Key = tKey
	t.Remove(from)
	t.addNode(to, n)
	
	return n, nil
}

/* Unexported */

// Adds an existing node to the trie -- modifies Key and Next
func (t *Trie) addNode(keyPath string, n *Node) *Node {
	if keyPath == "/" {
		// Root is eternal
		return nil
	}
	path, key := getPath(keyPath)

	var parentPath string
	if len(path) < 2 {
		parentPath = "/"
	} else {
		parentPath = mkPath(path[:len(path)-1])
	}
	p := t.Find(parentPath)
	if p == nil {
		// No parent to add from
		return nil
	}
	
	// Insert after p
	if n != p.Down {
		// If we are not already the first down
		n.Next = p.Down
		p.Down = n
	}
	n.Key = key
	return n
}

// Walks until it finds the parent of an element ;; Returns {parent, child}
func (t *Trie) getParent(path *[]string, key *string) (*Node, *Node) {
	var p *Node
	if len(*path) == 1 {
		// Base node
		p = t.Root
	} else {
		// Spook
		parentPath := mkPath((*path)[:len(*path)-1])
		parent := t.Find(parentPath)
		if parent == nil {
			// If there's no parent, there's no parent
			return nil, nil
		}
		p = parent
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
