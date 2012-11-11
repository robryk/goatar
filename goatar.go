package goatar

import "errors"
import "io"
import "archive/tar"
import "hash"
import "crypto/sha256"

//import "code.google.com/p/leveldb-go/leveldb/table"
import "github.com/robryk/goutils/teller"

var ErrIndexMismatch = errors.New("goatar: supplied index entry doesn't match the tar file")
var ErrIncompleteIndex = errors.New("goatar: supplied index entry is incomplete")

type Reader struct {
	r io.ReadSeeker
}

func (r *Reader) GetFile(f File) (metadata *tar.Header, contents io.Reader, err error) {
	if f.Offset == nil {
		err = ErrIncompleteIndex
		return
	}
	r.r.Seek(*f.Offset, 0)
	tarReader := tar.NewReader(r.r)
	if metadata, err = tarReader.Next(); err != nil {
		err = ErrIndexMismatch
		return
	}
	contents = tarReader
	return
}

type Writer struct {
	*tar.Writer
	ioWriter    io.WriteSeeker
	currentFile *File
	fileOutput  func(*File)
	hasher      hash.Hash
}

func NewWriter(w io.Writer, indexer func(*File)) *Writer {
	ws := teller.NewWriter(w)
	return &Writer{
		Writer:      tar.NewWriter(ws),
		ioWriter:    ws,
		currentFile: nil,
		fileOutput:  indexer,
	}
}

func (w *Writer) finishFile() {
	if w.currentFile == nil {
		return
	}

	w.currentFile.Hash = w.hasher.Sum(nil)
	w.hasher = nil

	w.fileOutput(w.currentFile)
	w.currentFile = nil
}

func (w *Writer) WriteHeader(hdr *tar.Header) error {
	w.finishFile()
	w.Flush()
	offset, err := w.ioWriter.Seek(1, 0)
	if err != nil {
		return err
	}
	w.currentFile = &File{
		Path:   &hdr.Name,
		Offset: &offset,
	}
	w.hasher = sha256.New()
	return w.Writer.WriteHeader(hdr)
}

func (w *Writer) Write(b []byte) (n int, err error) {
	n, err = w.Writer.Write(b)
	w.hasher.Write(b[:n])
	return
}

func (w *Writer) Close() error {
	w.finishFile()
	return w.Writer.Close()
}
