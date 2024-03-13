package merkledag

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type TestFile struct {
	name string
	data []byte
}

func (file *TestFile) Size() uint64 {
	return uint64(len(file.data))
}

func (file *TestFile) Name() string {
	return file.name
}

func (file *TestFile) Type() int {
	return FILE
}

func (file *TestFile) Bytes() []byte {
	return file.data
}

type testDirIter struct {
	list []Node
	iter int
}

func (iter *testDirIter) Next() bool {
	if iter.iter+1 < len(iter.list) {
		iter.iter += 1
		return true
	}
	return false
}

func (iter *testDirIter) Node() Node {
	return iter.list[iter.iter]
}

type TestDir struct {
	list []Node
	name string
}

func (dir *TestDir) Size() uint64 {
	var len uint64 = 0
	for i := range dir.list {
		len += dir.list[i].Size()
	}
	return len
}

func (dir *TestDir) Name() string {
	return dir.name
}

func (dir *TestDir) Type() int {
	return DIR
}

func (dir *TestDir) It() DirIterator {
	it := &testDirIter{
		list: dir.list,
		iter: -1,
	}
	return it
}

type HashMap struct {
	mp map[string]([]byte)
}

func (hmp *HashMap) Has(key []byte) (bool, error) {
	return hmp.mp[string(key)] != nil, nil
}

func (hmp *HashMap) Put(key, value []byte) error {
	flag, _ := hmp.Has(key)
	if flag {
		panic("Key is same")
	}
	hmp.mp[string(key)] = value
	return nil
}

func (hmp *HashMap) Get(key []byte) ([]byte, error) {
	flag, _ := hmp.Has(key)
	if !flag {
		panic("Don't have the key")
	}
	return hmp.mp[string(key)], nil
}

func (hmp *HashMap) Delete(key []byte) error {
	return nil
}

func TestDag(t *testing.T) {
	kv := &HashMap{
		mp: make(map[string][]byte),
	}
	h := sha256.New()
	//a single small file
	file := &TestFile{
		name: "small",
		data: []byte("很多人在第一次看到这个东西的时侯是非常兴奋的。不过这个自动机叫作 Automaton，不是 Automation，这里的 AC 也不是 Accepted，而是 Aho–Corasick（Alfred V. Aho, Margaret J. Corasick. 1975），让萌新失望啦。切入正题。似乎在初学自动机相关的内容时，许多人难以建立对自动机的初步印象，尤其是在自学的时侯。而这篇文章就是为你们打造的。笔者在自学 AC 自动机后花费两天时间制作若干的 gif，呈现出一个相对直观的自动机形态。尽管这个图似乎不太可读，但这绝对是在作者自学的时侯，画得最认真的 gif 了。另外有些小伙伴问这个 gif 拿什么画的。笔者用 Windows 画图软件制作。"),
	}
	root := Add(kv, file, h)
	fmt.Printf("%x\n", root)
	// a single big file
	kv = &HashMap{
		mp: make(map[string][]byte),
	}
	h.Reset()
	context, err := os.ReadFile("/Users/apple/Desktop/code/UniswapV2/package.json")
	if err != nil {
		t.Error(err)
	}

	file = &TestFile{
		name: "big",
		data: context,
	}

	root = Add(kv, file, h)
	fmt.Printf("%x\n", root)
	// a folder
	kv = &HashMap{
		mp: make(map[string][]byte),
	}
	h.Reset()
	path :=  "/Users/apple/Desktop/code/warp-hack-thon"
	files, _ := ioutil.ReadDir(path)
	dir := &TestDir{
		list: make([]Node, len(files)),
		name: "Documents",
	}
	for i, fi := range files {
		newPath := path + "/" + fi.Name()
		if fi.IsDir() {
			context := search(newPath)
			context.name = fi.Name()
			dir.list[i] = context
		} else {
			context, err := os.ReadFile(newPath)
			if err != nil {
				t.Fatal(err)
			}
			file = &TestFile{
				name: fi.Name(),
				data: context,
			}
			dir.list[i] = file
		}
	}
	root = Add(kv, dir, h)
	fmt.Printf("%x\n", root)
}

func search(path string) *TestDir {
	files, _ := ioutil.ReadDir(path)
	dir := &TestDir{
		list: make([]Node, len(files)),
	}
	for i, fi := range files {
		newPath := path + "/" + fi.Name()
		if fi.IsDir() {
			context := search(newPath)
			context.name = fi.Name()
			dir.list[i] = context
		} else {
			context, err := os.ReadFile(newPath)
			if err != nil {
				context := search(newPath)
				context.name = fi.Name()
				dir.list[i] = context
				continue
			}
			file := &TestFile{
				name: fi.Name(),
				data: context,
			}
			dir.list[i] = file
		}
	}
	return dir
}
