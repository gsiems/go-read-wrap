package srw

import (
	"io"
)

type multiReader struct {
	reader io.Reader
}

// Read reads from a MultiReader until either the supplied byte
// buffer is filled or the MultiReader returns an error/EOF.
func (mr *multiReader) readMulti(p []byte) (n int, err error) {

	var nc int

	// io.MultiReader appears to only read from one reader per
	// call. If the number of bytes requested exceeds the number
	// of bytes available to the current reader then additional
	// calls to Read are required until either the buffer fills
	// or the MultiReader returns an error.
	for {
		nc, err = mr.reader.Read(p[n:])
		n += nc
		if err != nil {
			break
		}
		if n >= len(p) {
			break
		}
	}
	return
}
