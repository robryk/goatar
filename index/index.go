package index

import "hash"
import "crypto/sha256"

type Indexer struct {
	file *IndexEntry
	hash.Hash
}

func NewIndexer(f *IndexEntry) *Indexer {
	return &Indexer{
		file: f,
		Hash: sha256.New(),
	}
}

func (i *Indexer) Close() *IndexEntry {
	i.file.Hash = i.Hash.Sum(nil)
	i.Hash = nil
	return i.file
}
