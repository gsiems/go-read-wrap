package srw

import (
	"io"
)

const (
	defaultBufSz = 4096
)

type buffReader struct {
	reader io.Reader
	buff   []byte
	bix    int // buff offset to start reading from
	bct    int // count of bytes read into buff
	err    error
}

func (mr *buffReader) readBuffer(p []byte) (n int, err error) {

	pLen := len(p)
	if pLen == 0 {
		return
	}

	var needed int

	for {

		// If the buffer is empty and there is no more to read then return
		if mr.atEOF() {
			err = io.EOF
			break
		}

		// If there is anything in the buffer then copy as much as is
		// needed/available
		avail := mr.bct - mr.bix
		needed = pLen - n
		if avail > 0 && needed > 0 {
			if avail > needed {
				// read as much as is needed
				copy(p[n:], mr.buff[mr.bix:mr.bix+needed])

				// zero out the read bytes
				mr.clearBytes(needed)

				n += needed
				mr.bix += needed
			} else {
				// read what there is
				copy(p[n:], mr.buff[mr.bix:])

				// zero out the read bytes
				mr.clearBytes(mr.bct - mr.bix)

				n += avail
				mr.bix = mr.bct
			}
		}

		// if the buffer is empty then refill it
		mr.checkFillBuffer()
		if mr.hasErr() {
			break
		}

		// if p is full then return
		if n == pLen {
			break
		}
	}

	if mr.hasErr() {
		err = mr.err
	}
	return
}

func (mr *buffReader) clearBytes(count int) {
	for i := 0; i < count; i++ {
		mr.buff[i+mr.bix] = 0
	}
}

func (mr *buffReader) hasErr() (t bool) {
	if mr.err != nil && mr.err != io.EOF {
		return true
	}
	return false
}

func (mr *buffReader) atEOF() (t bool) {
	if mr.bct == 0 && mr.err == io.EOF {
		return true
	}
	return false
}

// checkFillBuffer checks and if needed re-fills the buffer.
func (mr *buffReader) checkFillBuffer() {
	if mr.bix >= mr.bct {

		var n int
		mr.bix = 0
		mr.bct = 0

		// MultiReader appears to only read from one reader per
		// call. If the number of bytes requested exceeds the number
		// of bytes available to the current reader then additional calls
		// to Read are required until either the buffer fills or the
		// MultiReader returns an error.
		for mr.bct < len(mr.buff) {
			n, mr.err = mr.reader.Read(mr.buff[mr.bct:])
			mr.bct += n
			if mr.err != nil {
				break
			}
		}
	}
}
