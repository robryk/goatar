package index

import "hash"
import "crypto/sha256"

type Indexer struct {
	file *File
	hash.Hash
}

func NewIndexer(f *File) *Indexer {
	return &Indexer{
		file: f,
		Hash: sha256.New(),
	}
}

func (i *Indexer) Close() *File {
	i.file.Hash = i.Hash.Sum(nil)
	i.Hash = nil
	return i.file
}
