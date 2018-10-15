# fstrie
Filesystem style trie data structure.

fstrie is a string-indexed trie to represent a file system.

## Dependencies

None.

## Install

	go get github.com/henesy/fstrie

## Usage

Create a new trie:

	t := trie.New()

Add a node:

	node := t.Add("/tmp/data", mydata)

Remove a node:

	mydata := t.Remove("/tmp/data")

Find a node:

	node := t.Find("/tmp/data")

Get all children of a node:

	nodes := node.Children()

## References

- Ryan

