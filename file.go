package merkledag

const (
	FILE = iota
	DIR
)

type Node interface {
	Size() uint64
	Name() string
	Type() int
}

type File interface {
	Node

	Bytes() []byte
}

type Dir interface {
	Node

	It() DirIterator
}

type DirIterator interface {
	Next() bool

	Node() Node
}
