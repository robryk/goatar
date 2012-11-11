package index

import "fmt"
import "io"

type DebugIndexer struct {
	Writer io.Writer
}

func (di DebugIndexer) Index(ie *IndexEntry) error {
	_, err := fmt.Fprintf(di.Writer, "File %v begins at %v\n", ie.Path, ie.Offset)
	return err
}