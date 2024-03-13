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
	MAX_LISTLINE = 4096
	BLOB = "blob"
	LIST = "link"
	TREE = "tree"
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
	return processNode(node, store, h)
	
}
// 处理节点，返回默克尔树根
func processNode(node Node, store KVStore, h hash.Hash) []byte {
	obj := &Object{}
	switch node.Type() {
	case FILE:
		obj =  handleFile(node,store,h)
		break
	case DIR:
		obj =  handleDir(node,store,h)
		break	
	}
	jsonObj, _ := json.Marshal(obj)
	return computeHash(jsonObj,h)
	
}
// 处理文件，返回一个该文件对应的obj
func handleFile(node Node,store KVStore,h hash.Hash) *Object{
	obj := &Object{}
	FileNode,_:= node.(File)
		if FileNode.Size() > CHUNK_SIZE {

			// lowobj := Object{}
		    // 计算文件切片数量,向上取整
		numChunks := math.Ceil(float64(FileNode.Size()) / float64(CHUNK_SIZE))
		height := 0
		tmp := numChunks
		// 计算出要分几层
		for{
			height++
			tmp /= MAX_LISTLINE
			if tmp == 0{
			    break
			}
		}
		obj,_ = dfshandleFile(height,FileNode,store,0,h)
		}else{
			obj.Data = FileNode.Bytes()
			putObjInStore(obj,store,h)
			}
			return  obj
	
}

// 处理文件夹，返回对应的obj指针
func handleDir(node Node,store KVStore,h hash.Hash) *Object{
	dirNode , _ := node.(Dir)
	iter := dirNode.It()
	treeObject := &Object{}
	for iter.Next() {
		node := iter.Node()
		switch node.Type() {
		case FILE :
			file := node.(File)
			tmp := handleFile(node, store, h)
			jsonMarshal, _ := json.Marshal(tmp)
			treeObject.Links = append(treeObject.Links, Link{
				Hash: computeHash(jsonMarshal,h),
				Size: int(file.Size()),
				Name: file.Name(),
			})
			if tmp.Links == nil {
				treeObject.Data = append(treeObject.Data, []byte(BLOB)...)
			}else{
				treeObject.Data = append(treeObject.Data, []byte(LIST)...)
			}
			
			break
		case DIR :
			dir := node.(Dir)
			tmp := handleDir(node, store, h)
			jsonMarshal, _ := json.Marshal(tmp)
			treeObject.Links = append(treeObject.Links, Link{
				Hash: computeHash(jsonMarshal,h),
				Size: int(dir.Size()),
				Name: dir.Name(),
			})
			treeObject.Data = append(treeObject.Data, []byte(TREE)...)
			break
		}
	}
	putObjInStore(treeObject,store,h)
	return treeObject
}
// 处理大文件的方法 递归调用，返回当前生成的obj已经处理了多少数据
func dfshandleFile(height int, node File,store KVStore,start int,h hash.Hash) (*Object,int) {
	obj := &Object{}
	lendata := 0
	// 如果只分一层
 if height == 1{
	if(len(node.Bytes())- start < CHUNK_SIZE) {
		data := node.Bytes()[start:]
		obj.Data = append(obj.Data,data...)
		lendata = len(data)
		putObjInStore(obj,store,h)
		return obj,lendata
	}else{
	for i := 1; i<= MAX_LISTLINE; i++{
		end := start + CHUNK_SIZE
		// 确保不越界
		if end > len(node.Bytes()) {
			end = len(node.Bytes())
		}
		data :=  node.Bytes()[start:end]
		blobObj := Object{
			Links : nil,
			Data : data,
		}
		putObjInStore(&blobObj,store,h)
		jsonMarshal,_ := json.Marshal(blobObj)
		obj.Links = append(obj.Links,Link{
			Hash: computeHash(jsonMarshal,h),
			Size: int(len(data)),
		})
		obj.Data  = append(obj.Data,[]byte(BLOB)...)
		lendata += len(data)
		start += CHUNK_SIZE
		if start >= len(node.Bytes()){
		    break
		}
	}
	putObjInStore(obj,store,h)
		return obj,lendata
	}
 }else{
	for i := 1 ;i<=MAX_LISTLINE;i++{
		if start >= len(node.Bytes()){
			break
		}
		tmpObj,tmpLendata := dfshandleFile(height-1,node,store,start,h)
		lendata += tmpLendata
		jsonMarshal,_ := json.Marshal(tmpObj)
		obj.Links = append(obj.Links,Link{
			Hash: computeHash(jsonMarshal,h),
			Size: tmpLendata,
		})
		if(tmpObj.Links == nil){
		obj.Data  = append(obj.Data,[]byte(BLOB)...)
		}else{
			obj.Data = append(obj.Data,[]byte(LIST)...)
		}
		start += tmpLendata
	}
	putObjInStore(obj,store,h)
	return obj,lendata
 }
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
	has ,_ := store.Has(hash)
	if has { return }
	// fmt.Println("put obj in store:",hash)
	store.Put(hash,value)
	
}

