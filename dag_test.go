package merkledag

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

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
	context, err := os.ReadFile("/Users/apple/Desktop/code/jsnixiang/1/js.js")
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
	// path := "/home/xiuuix/go"
	path := "/Users/apple/Desktop/code/jsnixiang"
	files, _ := ioutil.ReadDir(path)
	dir := &TestDir{
		list: make([]Node, len(files)),
		name: "/",
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

func TestDag2file(t *testing.T) {
	kv := &HashMap{
		mp: make(map[string][]byte),
	}
	h := sha256.New()
	// a folder
	kv = &HashMap{
		mp: make(map[string][]byte),
	}
	// path := "/home/xiuuix/go"
	path := "/Users/apple/Desktop/code/jsnixiang"
	files, _ := ioutil.ReadDir(path)
	dir := &TestDir{
		list: make([]Node, len(files)),
		name: "/",
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
			file := &TestFile{
				name: fi.Name(),
				data: context,
			}
			dir.list[i] = file
		}
	}
	root := Add(kv, dir, h)
	fmt.Printf("%xroot\n", root)
	// buffer_go := Hash2File(kv, root, "/pkg/mod/bazil.org/fuse@v0.0.0-20200117225306-7b5117fecadc/buffer.go", nil)
	// buffer_go := Hash2File(kv, root, "/Users/apple/Desktop/code/jsnixiang/1/js.js", nil)
	// fmt.Println("test",string(buffer_go),"go")
	dlv := Hash2File(kv, root, "/1/js.js", nil)
	context, _ := os.ReadFile("/Users/apple/Desktop/code/jsnixiang/1/js.js")
	h.Reset()
	h.Write(context)
	tmp1 := h.Sum(nil)
	h.Reset()
	h.Write(dlv)
	tmp2 := h.Sum(nil)
	fmt.Println(tmp1)
	fmt.Println(tmp2)
	fmt.Println(string(tmp1) == string(tmp2))
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