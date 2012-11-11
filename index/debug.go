package index

import "fmt"
import "io"

type DebugIndexer struct {
	Writer io.Writer
}

func NewDebugIndexer(w io.Writer) DebugIndexer {
	return DebugIndexer{w}
}

func (di DebugIndexer) Index(ie *IndexEntry) error {
	if ie.Offset != nil {
		_, err := fmt.Fprintf(di.Writer, "File %v begins at %v\n", *ie.Path, *ie.Offset)
		return err
	}
	return nil
}
