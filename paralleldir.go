package main

import "fmt"

type ParallelDir struct {
	baseDir string
}

func NewParallelDir(baseDir string) *ParallelDir {
	return &ParallelDir{baseDir: baseDir}
}

func (pd *ParallelDir) Run() {
	fmt.Println("Starting on", pd.baseDir)
}
