package v2

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
}

type nodeUpdate struct {
	Parent *treeNode
	Child  *treeNode
}

type ParallelDir struct {
	Root           *treeNode `json:"root"`
	nodeUpdateChan chan *nodeUpdate
}

func NewParallelDir(baseDir string) *ParallelDir {
	return &ParallelDir{
		Root: &treeNode{
			Name:     baseDir,
			Children: make([]*treeNode, 0),
		},
		nodeUpdateChan: make(chan *nodeUpdate),
	}
}

func (pd *ParallelDir) baseDir() string {
	return pd.Root.Name
}

// appendChildren consumes update requests sent to the ParallelDir.nodeUpdateChan and performs the update
func (pd *ParallelDir) appendChildren() {
	for updateRequest := range pd.nodeUpdateChan {
		updateRequest.Parent.Children = append(updateRequest.Parent.Children, updateRequest.Child)
	}
}

func (pd *ParallelDir) Run() {
	// WG for the parallel traversal
	var wg sync.WaitGroup
	// WG because the listDirsRecursively might finish, the wg.Wait() "unlocks" and the app might finish
	// before the updater manages to read all messages.
	var updaterWG sync.WaitGroup
	fmt.Println("Starting on", pd.baseDir())

	// pd.appendChildren() needs to be running before listDirsRecursively writes to chan
	updaterWG.Add(1)
	go func() {
		updaterWG.Done()
		pd.appendChildren()
	}()

	pd.Root.listDirsRecursively(&wg, pd.nodeUpdateChan)
	wg.Wait()
	close(pd.nodeUpdateChan) // closing chan after the parallel dir traversal finishes
	updaterWG.Wait()         // Waiting after the parallel dir traversal finishes

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

func (t *treeNode) listDirsRecursively(wg *sync.WaitGroup, ch chan<- *nodeUpdate) {
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
			//t.safeAppendChild(childNode)
			update := nodeUpdate{
				Parent: t,
				Child:  childNode,
			}
			ch <- &update

			wg.Go(func() {
				childNode.listDirsRecursively(wg, ch)
			})
		}
	}
}
