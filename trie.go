// Package fstrie provides a trie data structure that is string-addressed and unsorted.
// Nodes are singly-linked and do not reference their parents and data storage is an interface{}.
// fstrie has a fairly strict syntax which is remniscient of file paths in a traditional unix-like operating system. 
// In general, fstrie expects correctly formatted paths, though it handles non-existent nodes.
package fstrie

import (
	"strings"
	"fmt"
)


// Node for trie tree
type Node struct {
	Key		string
	Data	interface{}
	Next	*Node
	Down	*Node
}

// Trie tree container
type Trie struct {
	Root	*Node
}


// New creates an initialized Trie
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
	path, key := GetPath(keyPath)

	var parentPath string
	if len(path) < 2 {
		parentPath = "/"
	} else {
		parentPath = MkPath(path[:len(path)-1])
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

// Remove a node (by Find) ;; Returns removed node's Data
func (t *Trie) Remove(keyPath string) interface{} {
	if keyPath == "/" {
		// Not allowed to delete root
		return nil
	}
	path, key := GetPath(keyPath)

	p, _ := t.GetParent(&path, &key)
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

// Find a node by string index via trie-traversal ;; Returns the found node or nil
func (t *Trie) Find(keyPath string) (*Node) {
	if keyPath == "/" {
		return t.Root
	}
	path, key := GetPath(keyPath)
	if t.Root == nil {
		return nil
	}
	cursor := t.Root

	for i := 0; i < len(path)+1; i++ {
		if cursor.Key == key && i == len(path) {
			// Last element is our file, if exists
			return cursor
		}
		cursor = cursor.GetChild(path[i])
		if cursor == nil {
			// There are no more paths to walk
			break
		}
	}
	return nil
}

// Find a node's value by asking parent nodes for the child and return the "found" value (virtual get)
func (t *Trie) Get(keyPath string) (interface{}) {
	path, _ := GetPath(keyPath)
	// Prepend / to make recursion clean
	path = append([]string{"/"}, path...)
	return Walk(t.Root, path)
}

// Perform a walking operation to recursively descend further into the trie (if able)
func Walk(n *Node, path []string) (interface{}) {
	/* The desired functionality is to replace this method with a handler from the wrapper program */
	children := n.Children()
	
	if len(path) < 1 {
		// Should never happen
		return nil
	} else if len(children) < 1 && len(path) > 1 {
		// This is a data node, has no children, but looking for a child -- best case to replace Walk() for
		return nil
	} else if len(path) == 1 {
		// The node is us
		return n.Data
	} else {
		// We must recurse further, captain -- strip ourselves out
		path = path[1:]
		for _, c := range children {
			if c.Key == path[0] {
				return Walk(c, path)
			}
		}
	}
	return nil
}

// Return a string version of the trie in the style of tree(1)
func (t *Trie) String() (out string) {
	out += ""
	indent := 0
	t.Root.String(&out, indent)
	return
}

// Children returns the set of top-level children for a node
func (n *Node) Children() [](*Node) {
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

// Existent portion of a path is extracted and returned
func (t *Trie) Existent(keyPath string) string {
	if keyPath == "/" {
		return keyPath
	}
	path, _ := GetPath(keyPath)
	if t.Root == nil {
		return ""
	}
	cursor := t.Root
	realPath := ""

	for i := 0; i < len(path); i++ {
		cursor = cursor.GetChild(path[i])
		if cursor == nil {
			// There are no more paths to walk
			break
		}
		realPath += "/" + cursor.Key
	}
	return realPath
}

// Mv Removes then re-Adds a node into a different location ;; will overwrite existing name 'to'
func (t *Trie) Mv(from, to string) (*Node, error) {
	if from == "/" || to == "/" {
		// Not allowed to change "/"
		return nil, fmt.Errorf("Root is immutable.")
	}

	tPath, tKey := GetPath(to)
	n := t.Find(from)
	if n == nil {
		return nil, fmt.Errorf("Node %v not found.", from)
	}

	if from == to {
		return n, nil
	}

	// Maybe this should be its own fn
	parentPath := MkPath(tPath[:len(tPath)-1])
	if t.Find(parentPath) == nil {
		// Parent does not exist
		return nil, fmt.Errorf("Parent node %v not found.", parentPath)
	}
	
	n.Key = tKey
	t.Remove(from)
	t.AddNode(to, n)
	
	return n, nil
}

// Build the string representation of a Node
func (n *Node) String(out *string, indent int) {
	// Set our key
	*out += n.Key + "\n"
	
	// Set children's keys and call their string()
	cs := n.Children()
	for p, v := range cs {
		lead := ""
		if p == len(cs)-1 {
			lead = "└"
		} else {
			lead = "├"
		}
		
		lead += "─"
		
		for i := 0; i < indent; i++ {
			*out += "│ "
		}
		*out += lead
		v.String(out, indent+1)
	}
}

// Adds an existing node to the trie -- modifies Key and Next
func (t *Trie) AddNode(keyPath string, n *Node) *Node {
	if keyPath == "/" {
		// Root is eternal
		return nil
	}
	path, key := GetPath(keyPath)

	var parentPath string
	if len(path) < 2 {
		parentPath = "/"
	} else {
		parentPath = MkPath(path[:len(path)-1])
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
func (t *Trie) GetParent(path *[]string, key *string) (*Node, *Node) {
	var p *Node
	if len(*path) == 1 {
		// Base node
		p = t.Root
	} else {
		// Spook
		parentPath := MkPath((*path)[:len(*path)-1])
		parent := t.Find(parentPath)
		if parent == nil {
			// If there's no parent, there's no parent
			return nil, nil
		}
		p = parent
	}

	c := p.GetChild(*key)
	if c == nil {
		// If there's no child named key, there's no parent
		return nil, nil
	}

	return p, c
}

// Does a down → (next…) while able and matches a child with Key key
func (n *Node) GetChild(key string) *Node {
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
func GetPath(keyPath string) (path []string, key string) {
	// Due to how split works, if we make a path with a leading /, strip 0th
	path = strings.Split(keyPath, "/")[1:]
	key = path[len(path)-1]
	return
}

// Utility to generate path strings -- slow
func MkPath(path []string) string {
	keyPath := ""
	for _, v := range path {
		keyPath += "/" + v
	}
	return keyPath
}
