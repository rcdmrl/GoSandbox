package main

import (
	fstreev1 "github.com/rcdmrl/go-sandbox/fstree/v1"
)

func main() {
	pd := fstreev1.NewParallelDir("/Users/ricaamar/Documents/")
	pd.Run()
}
