package merkledag

import (
	"encoding/json"
	"fmt"
	"hash"
	"math"
)

const (
	K = 1 << 10
	M = K << 10
	CHUNK_SIZE = 256 * K
	BLOB = "blob"
	LIST = "list"
	TREE = "tree"
)
const (
	
)
type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// TODO 将分片写入到KVStore中，并返回Merkle Root
	
	// 递归处理文件夹和文件，并将文件内容保存在 KVStore 中
	return processNode(node, store, h)
	
}
// 处理节点，返回默克尔树根
func processNode(node Node, store KVStore, h hash.Hash) []byte {
	obj := Object{}
	switch node.Type() {
	case FILE:
		hash :=  handleFile(node,store,h)
		obj.Links = append(obj.Links, Link{Name: node.Name(), Hash: hash, Size: int(node.Size())})
		obj.Data = append(obj.Data,)

		
	case DIR:
		hash :=  handleDir(node,store,h)
		obj.Links = append(obj.Links,Link{Name: node.Name(), Hash: hash, Size: int(node.Size())})
		obj.Data = append(obj.Data,LIST...)
	}
	return nil
	
	
}
 // 处理文件，返回文件的默克尔树根
func handleFile(node Node,store KVStore,h hash.Hash) []byte{
	obj := Object{}
	FileNode,ok := node.(File)
		if !ok {
			fmt.Println("error")
		    return nil
		}
		if FileNode.Size() > CHUNK_SIZE {
			lowobj := Object{}
		    // 计算文件切片数量,向上取整
		numChunks := uint64(math.Ceil(float64(FileNode.Size()) / float64(CHUNK_SIZE)))
		for i := uint64(0);i<numChunks;i++{
			// 当前分片大小
			size := uint64(CHUNK_SIZE)
			if i == numChunks -1 {
				size =  FileNode.Size() - (i*CHUNK_SIZE)
			}
			chunk := FileNode.Bytes()[i*CHUNK_SIZE : (i+1)*size]
			// 计算当前片的hash放入堆栈中，和放入存储器
   			hash := computeHash(chunk, h)
			lowobj.Links = append(obj.Links, Link{node.Name(), hash, int(size)})
			lowobj.Data = append(obj.Data,BLOB...)
			value,err := json.Marshal(lowobj)
			if err != nil{
				fmt.Println("json.Marshal err:",err)
				return nil
			}
			store.Put(hash,value)
		}
		hashes := getHashes(&lowobj)
		hash :=  computeMerkleRoot(hashes,h)
		obj.Links = append(obj.Links, Link{FileNode.Name(), hash, int(FileNode.Size())})
		obj.Data = append(obj.Data,LIST...)
		putObjInStore(&obj,store,h)
		}else{
			hash := computeHash(FileNode.Bytes(), h)
			obj.Links = append(obj.Links, Link{FileNode.Name(), hash, int(node.Size())})
			obj.Data = append(obj.Data,BLOB...)
			putObjInStore(&obj,store,h)
			return hash
		}
	 return nil
}
// 处理文件夹，返回默克尔树根
func handleDir(node Node,store KVStore,h hash.Hash) []byte{
	obj := Object{}
	dirNode,ok := node.(Dir)
	if !ok {
		return  nil
	} 
	hashes := [][]byte{{},{}}
	iter:= dirNode.It()
	for iter.Next() {
		hash :=  processNode(iter.Node(),store,h)
		if hash != nil{
		    hashes = append(hashes,hash)
		}
	}
	hash :=computeMerkleRoot(hashes,h)
	obj.Links = append(obj.Links, Link{node.Name(), hash, int(node.Size())})
	obj.Data = append(obj.Data,TREE...)
	putObjInStore(&obj,store,h)
	return hash
}

func computeHash(data []byte, h hash.Hash) []byte {
	h.Reset()
	h.Write(data)
	return h.Sum(nil)
}
func putObjInStore(obj *Object, store KVStore, h hash.Hash){
	value,err := json.Marshal(obj)
	if err != nil{
		fmt.Println("json.Marshal err:",err)
		return
	}
	hash := computeHash(value, h)
	store.Put(hash,value)
}
func getHashes(obj *Object) [][]byte {
 hashes := make([][]byte, len(obj.Links))
 for i, link := range obj.Links {
  hashes[i] = link.Hash
 }
 return hashes
}
func computeMerkleRoot(data [][]byte, h hash.Hash) []byte {
	if len(data) == 0{
		return nil
	}
	if len(data) == 1{
		return data[0]
	}
	var nextLevel [][]byte
	// 对于相邻节点计算hash
	for i := 0; i < len(data); i += 2 {
	    // 确保不出界
		end := i+2
		if(end >len(data)){
			end = len(data)
		}
		// 拼接两个叶子结点的hash
		hash := computeHash(append(data[i], data[i+1]...), h)
		nextLevel = append(nextLevel, hash[:])
	}
	// 递归计算下一层的默克尔树根
	return computeMerkleRoot(nextLevel, h)
}

