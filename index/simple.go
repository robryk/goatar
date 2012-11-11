package index

import "io"

import "code.google.com/p/goprotobuf/proto"
import "code.google.com/p/leveldb-go/leveldb/record"

type SimpleIndexer record.Writer

func NewSimpleIndexer(w io.Writer) *SimpleIndexer {
	return (*SimpleIndexer)(record.NewWriter(w))
}

func (si *SimpleIndexer) Index(f *IndexEntry) error {
	w, err := (*record.Writer)(si).Next()
	if err != nil {
		return err
	}

	var output []byte
	output, err = proto.Marshal(f)
	if err != nil {
		return err
	}

	_, err = w.Write(output)
	return err
}

func (si *SimpleIndexer) Close() error {
	return (*record.Writer)(si).Close()
}
