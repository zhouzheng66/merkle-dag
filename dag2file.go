package merkledag

import (
	"encoding/json"
	"strings"
)

// Hash to file

// example a path : /doc/tmp/temp.txt
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 根据hash和path， 返回对应的文件, hash对应的类型是tree
	flag, _ := store.Has(hash)
	if flag {
		objBinary, _ := store.Get(hash)
		var obj Object
		json.Unmarshal(objBinary, &obj)
		pathArr := strings.Split(path, "/")
		cur := 1
		return getFileByDir(obj, pathArr, cur, store)
	}
	return nil
}

func getFileByDir(obj Object, pathArr []string, cur int, store KVStore) []byte {
	if cur >= len(pathArr) {
		return nil
	}
	index := 0
	for i := range obj.Links {
		objType := string(obj.Data[index : index+4])
		index += 4
		objInfo := obj.Links[i]
		if objInfo.Name != pathArr[cur] {
			continue
		}
		switch objType {
		case TREE:
			objDirBinary, _ := store.Get(objInfo.Hash)
			var objDir Object
			json.Unmarshal(objDirBinary, &objDir)
			ans := getFileByDir(objDir, pathArr, cur+1, store)
			if ans != nil {
				return ans
			}
		case BLOB:
			ans, _ := store.Get(objInfo.Hash)
			return ans
		case LIST:
			objLinkBinary, _ := store.Get(objInfo.Hash)
			var objLink Object
			json.Unmarshal(objLinkBinary, &objLink)
			ans := getFileByList(objLink, store)
			return ans
		}
	}
	return nil
}

func getFileByList(obj Object, store KVStore) []byte {
	ans := make([]byte, 0)
	index := 0
	for i := range obj.Links {
		curObjType := string(obj.Data[index : index+4])
		index += 4
		curObjLink := obj.Links[i]
		curObjBinary, _ := store.Get(curObjLink.Hash)
		var curObj Object
		json.Unmarshal(curObjBinary, &curObj)
		if curObjType == BLOB {
			ans = append(ans, curObjBinary...)
		} else { //List
			tmp := getFileByList(curObj, store)
			ans = append(ans, tmp...)
		}
	}
	return ans
}
