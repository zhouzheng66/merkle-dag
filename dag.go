package merkledag

import "hash"

type IPFSLink struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []IPFSLink
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// TODO 将分片写入到KVStore中，并返回Merkle Root
	return nil
}
