package merkledag

import (
	"hash"
)

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// TODO 将分片写入到KVStore中，并返回Merkle Root
	// 递归处理文件夹和文件，并将文件内容保存在 KVStore 中
	// 用一个二维切片来表示堆栈结构，存储上一个堆栈和当前的数据
	var hashes [][]byte
	processNode(node, store, h, &hashes)

	// 构建 Merkle 树并计算 Merkle Root
	merkleRoot := computeMerkleRoot(hashes, h)
	return merkleRoot
	
}

func processNode(node Node, store KVStore, h hash.Hash, hashes *[][]byte) {
	switch node.Type() {
	case 0:
		FileNode,ok := node.(File)
		if !ok {
		    return 
		}
		// 判断是否大于256kb 进行切片
		chunkSize := uint64(256)
		// 计算文件切片数量
		numChunks := (FileNode.Size() / chunkSize)
		for i := uint64(0);i<numChunks;i++{
			// 当前分片大小
			size := chunkSize
			if 1 == numChunks -1 {
				size =  FileNode.Size() - (i*chunkSize)
			}
			chunk := FileNode.Bytes()[i*chunkSize : (i+1)*size]
			// 计算当前片的hash放入堆栈中，和放入存储器
   			hash := computeHash(chunk, h)
			store.Put(hash, chunk)
			*hashes = append(*hashes,hash)
		}

		
	case 1:
		// 如果是文件夹，则将其转换为文件夹接口类型
		dirNode, ok := node.(Dir)
		if !ok {

			return 
		}
		// 循环递归处理这个文件夹
		iter := dirNode.It()
		for iter.Next() {
			processNode(iter.Node(), store, h, hashes)
		}
	}
	
	
}

func computeHash(data []byte, h hash.Hash) []byte {
	h.Reset()
	h.Write(data)
	return h.Sum(nil)
}

func computeMerkleRoot(hashes [][]byte, h hash.Hash) []byte {
	// 逐层计算 Merkle 树的哈希值
	for len(hashes) > 1 {
		var newHashes [][]byte
		for i := 0; i < len(hashes); i += 2 {
			// 如果有奇数个哈希，则复制最后一个哈希
			if i+1 == len(hashes) {
				newHashes = append(newHashes, hashes[i])
			} else {
				// 计算两个哈希的父哈希
				h.Reset()
				h.Write(hashes[i])
				h.Write(hashes[i+1])
				parentHash := h.Sum(nil)
				newHashes = append(newHashes, parentHash)
			}
		}
		hashes = newHashes
	}
	return hashes[0]
}

