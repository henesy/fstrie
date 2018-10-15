# fstrie

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

Find the most complete path to a file:

	path := t.Existent("/tmp/data/newt")

Move (and/or rename) a node:

	node, err := t.Mv("/tmp/a", "/tmp/b")

Print a string representation of the trie:

	fmt.Println(t.String())

## Testing

	go test

## References

- @RyanRadomski7

