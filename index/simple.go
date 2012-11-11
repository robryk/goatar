package index

import "code.google.com/p/goprotobuf/proto"
import "code.google.com/p/leveldb-go/leveldb/record"

type SimpleIndexer record.Writer

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
