package v2

import (
	"encoding/json"
	"testing"
)

func TestTreeNode_ToString(t *testing.T) {
	root := &treeNode{
		Name:     "root",
		Children: []*treeNode{},
	}
	child := &treeNode{
		Name:     "child",
		Children: []*treeNode{},
	}
	root.Children = append(root.Children, child)

	result := root.toString()

	expected := "root\n\tchild\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// func TestTreeNode_ListDirectoriesRecursivelly(t *testing.T) {
// 	// TODO: How?
// }

func TestParallelDir_NewAndBaseDir(t *testing.T) {
	rootDir := "/some/path"
	pd := NewParallelDir(rootDir)

	if pd.baseDir() != rootDir {
		t.Errorf("expected %s, got %s", rootDir, pd.baseDir())
	}

	if pd.Root == nil {
		t.Error("expected Root to be initialized")
	}
}

func TestParallelDir_ToJson(t *testing.T) {
	// roundtrip test
	pd := NewParallelDir("rootdir")
	pd.Root.Children = append(pd.Root.Children, &treeNode{
		Name:     "childdir",
		Children: []*treeNode{},
	})

	// tp json
	jsonStr := pd.toJson()

	// back to a struct
	var parsed treeNode
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if parsed.Name != pd.Root.Name {
		t.Errorf("expected name '%s', got %s", pd.Root.Name, parsed.Name)
	}
	if len(parsed.Children) != len(pd.Root.Children) {
		t.Errorf("expected %d children nodes, got %d", len(pd.Root.Children), len(parsed.Children))
	}
}
