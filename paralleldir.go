package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ParallelDir struct {
	root *treeNode
}

type treeNode struct {
	name     string
	children []*treeNode
	mu       sync.Mutex
}

func (t *treeNode) SafeAppendChild(childNode *treeNode) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.children = append(t.children, childNode)
}

func (t *treeNode) ToString() string {
	return t.toStringWithIndent("")
}

func (t *treeNode) toStringWithIndent(indent string) string {
	var result strings.Builder

	result.WriteString(indent + t.name + "\n")
	for _, child := range t.children {
		result.WriteString(child.toStringWithIndent(indent + "\t"))
	}

	return result.String()
}

func NewParallelDir(baseDir string) *ParallelDir {
	return &ParallelDir{
		root: &treeNode{
			name:     baseDir,
			children: make([]*treeNode, 0),
		},
	}
}

func (pd *ParallelDir) baseDir() string {
	return pd.root.name
}

func (pd *ParallelDir) Run() {
	var wg sync.WaitGroup
	fmt.Println("Starting on", pd.baseDir())
	listDirsRecursively(pd.root, &wg)
	wg.Wait()
	fmt.Println(pd.root.ToString())
}

func listDirsRecursively(node *treeNode, wg *sync.WaitGroup) {
	files, err := os.ReadDir(node.name)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			childNode := &treeNode{
				name:     filepath.Join(node.name, file.Name()),
				children: make([]*treeNode, 0),
			}
			node.SafeAppendChild(childNode)
			wg.Go(func() {
				listDirsRecursively(childNode, wg)
			})
		}
	}
}
