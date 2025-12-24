package fstree

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type treeNode struct {
	Name     string      `json:"name"`
	Children []*treeNode `json:"children"`
	mu       sync.Mutex
}

type ParallelDir struct {
	Root *treeNode `json:"root"`
}

func NewParallelDir(baseDir string) *ParallelDir {
	return &ParallelDir{
		Root: &treeNode{
			Name:     baseDir,
			Children: make([]*treeNode, 0),
		},
	}
}

func (pd *ParallelDir) baseDir() string {
	return pd.Root.Name
}

func (pd *ParallelDir) Run() {
	var wg sync.WaitGroup
	fmt.Println("Starting on", pd.baseDir())
	pd.Root.listDirsRecursively(&wg)
	wg.Wait()
	fmt.Println(pd.Root.toString())
	fmt.Println(strings.Repeat("-", 100))
	fmt.Println(pd.toJson())
}

func (pd *ParallelDir) toJson() string {
	bytes, err := json.Marshal(pd.Root)
	if err != nil {
		log.Fatal("Error converting final tree to JSON:", err.Error())
	}
	return string(bytes)
}

func (t *treeNode) safeAppendChild(childNode *treeNode) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Children = append(t.Children, childNode)
}

func (t *treeNode) toString() string {
	return t.toStringWithIndent("")
}

func (t *treeNode) toStringWithIndent(indent string) string {
	var result strings.Builder

	result.WriteString(indent + t.Name + "\n")
	for _, child := range t.Children {
		result.WriteString(child.toStringWithIndent(indent + "\t"))
	}

	return result.String()
}

func (t *treeNode) listDirsRecursively(wg *sync.WaitGroup) {
	files, err := os.ReadDir(t.Name)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			childNode := &treeNode{
				Name:     filepath.Join(t.Name, file.Name()),
				Children: make([]*treeNode, 0),
			}
			t.safeAppendChild(childNode)
			wg.Go(func() {
				childNode.listDirsRecursively(wg)
			})
		}
	}
}
