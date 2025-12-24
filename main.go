package main

import (
	fstreev1 "github.com/rcdmrl/go-sandbox/fstree/v1"
	fstreev2 "github.com/rcdmrl/go-sandbox/fstree/v2"

	"log"

	tuiv1 "github.com/rcdmrl/go-sandbox/tui/v1"
)

func main() {
	pd1 := fstreev1.NewParallelDir("/Users/ricaamar/Documents/")
	pd2 := fstreev2.NewParallelDir("/Users/ricaamar/Documents/")

	form1 := tuiv1.NewMainForm(pd1, pd2)

	err := form1.Run()
	if err != nil {
		log.Fatal(err)
	}
	err = form1.Dispatch()
	if err != nil {
		log.Fatal(err)
	}
}
