package index

type Indexer interface {
	Index(ie *IndexEntry) error
	Close() error
}
