package utils

import "io"

// ReadSeekerAt is a wrapper around io.ReadSeeker that implements io.ReaderAt
type ReadSeekerAt struct {
	io.ReadSeeker
}

// ReadAt reads len(p) bytes into p starting at offset off in the underlying input source.
func (r *ReadSeekerAt) ReadAt(p []byte, off int64) (n int, err error) {
	// Save the current position
	cur, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Go to the desired offset
	if _, err = r.Seek(off, io.SeekStart); err != nil {
		return 0, err
	}

	// Read at the desired offset
	n, err = r.Read(p)

	// Restore the original position
	_, seekErr := r.Seek(cur, io.SeekStart)
	if seekErr != nil {
		return 0, seekErr
	}

	return n, err
}
