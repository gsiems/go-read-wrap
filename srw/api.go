// Package srw is a simple io.Reader wrapper.
// The intent is to provide a buffered reader, multi-reader, and
// buffered multi-reader that all behave the same.
package srw

import (
	"io"
)

// BuffReader returns an io.Reader. This wraps an io.Reader while
// adding a buffer to minimize disk reads when performing mostly
// small reads.
func BuffReader(sz int, reader io.Reader) io.Reader {

	var mr buffReader
	mr.reader = reader
	if sz <= 0 {
		mr.buff = make([]byte, defaultBufSz)
	} else {
		mr.buff = make([]byte, sz)
	}
	return &mr
}

// MultiReader returns an io.Reader. This wraps the io.MultiReader to
// ensure that switching from one reader to the next is transparent
// and appears as though the data was read from one continuous reader.
func MultiReader(readers ...io.Reader) io.Reader {

	var mr multiReader
	mr.reader = io.MultiReader(readers...)
	return &mr
}

// BuffMultiReader returns an io.Reader. This wraps a MultiReader
// while adding a buffer to minimize disk reads when performing mostly
// small reads.
func BuffMultiReader(sz int, readers ...io.Reader) io.Reader {

	var mr buffReader
	mr.reader = MultiReader(readers...)
	if sz <= 0 {
		mr.buff = make([]byte, defaultBufSz)
	} else {
		mr.buff = make([]byte, sz)
	}
	return &mr
}

// Read reads from either a BuffReader or a BuffMultiReader until
// either the supplied byte buffer is filled or the reader returns
// an error/EOF.
func (mr *buffReader) Read(p []byte) (n int, err error) {
	n, err = mr.readBuffer(p)
	return
}

// Read reads from a MultiReader until either the supplied byte
// buffer is filled or the MultiReader returns an error/EOF.
func (mr *multiReader) Read(p []byte) (n int, err error) {
	n, err = mr.readMulti(p)
	return
}
