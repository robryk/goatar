package main

import "flag"
import "fmt"
import "os"

import "github.com/robryk/goatar/tar"
import "github.com/robryk/goatar/index"

var indexerName = flag.String("indexer", "simple", "Indexer to use for outputting")

func main() {
	flag.Parse()

	var indexer index.Indexer
	if *indexerName == "simple" {
		indexer = index.NewSimpleIndexer(os.Stdout)
	} else if *indexerName == "debug" {
		indexer = index.NewDebugIndexer(os.Stdout)
	} else {
		panic("invalid indexer name")
	}
	err := tar.Index(os.Stdin, indexer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
