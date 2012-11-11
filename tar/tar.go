package tar

import "archive/tar"
import "errors"
import "io"

import "github.com/robryk/goutils/teller"
import "github.com/robryk/goatar/index"

var ErrIndexMismatch = errors.New("goatar: supplied index entry doesn't match the tar file")
var ErrIncompleteIndex = errors.New("goatar: supplied index entry is incomplete")

type Reader struct {
	r io.ReadSeeker
}

func (r *Reader) GetFile(f index.IndexEntry) (metadata *tar.Header, contents io.Reader, err error) {
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
	currentFile *index.Indexer
	fileOutput  func(*index.IndexEntry)
}

func NewWriter(w io.Writer, indexer func(*index.IndexEntry)) *Writer {
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

	w.fileOutput(w.currentFile.Close())
}

func (w *Writer) WriteHeader(hdr *tar.Header) error {
	w.finishFile()
	w.Flush()
	offset, err := w.ioWriter.Seek(1, 0)
	if err != nil {
		return err
	}

	nextFile := &index.IndexEntry{
		Path:   &hdr.Name,
		Offset: &offset,
	}
	w.currentFile = index.NewIndexer(nextFile)

	return w.Writer.WriteHeader(hdr)
}

func (w *Writer) Write(b []byte) (n int, err error) {
	n, err = w.Writer.Write(b)
	w.currentFile.Write(b[:n])
	return
}

func (w *Writer) Close() error {
	w.finishFile()
	return w.Writer.Close()
}

const tarBlockSize = 512

func Index(r io.Reader, indexer func(*index.IndexEntry) error) error {
	sr := teller.NewReader(r)
	tr := tar.NewReader(sr)
	for {
		// This depends on the implementation of tar.Reader
		// We want to find the exact offset of the beginning of the header. The offset we get before the header is read is at the end of the previous file's data (because we've read all of its data).
		// Alas, tar pads files to multiples of block size and the padding is read only when Reader.Next() is invoked. So we get the offset of the end of the previous file's data and round up to
		// a multiple of blockSize bytes.
		offset, err := sr.Seek(1, 0)
		if err != nil {
			return err
		}
		offset = (offset + tarBlockSize - 1) & ^(tarBlockSize - 1)

		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		file := &index.IndexEntry{
			Path:   &hdr.Name,
			Offset: &offset,
		}
		idx := index.NewIndexer(file)

		_, err = io.Copy(idx, tr)
		if err != nil {
			return err
		}

		err = indexer(idx.Close())
		if err != nil {
			return err
		}
	}
	panic("notreached")
}
