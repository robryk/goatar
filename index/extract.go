package index

import "hash"
import "crypto/sha256"

type Extractor struct {
	file *IndexEntry
	hash.Hash
}

func NewExtractor(f *IndexEntry) *Extractor {
	return &Extractor{
		file: f,
		Hash: sha256.New(),
	}
}

func (e *Extractor) Close() *IndexEntry {
	e.file.Hash = e.Hash.Sum(nil)
	e.Hash = nil
	return e.file
}
