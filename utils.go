package spdy

import (
	"io"
	"net/http"
)

// cloneHeader returns a duplicate of the provided Header.
func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// updateHeader adds and new name/value pairs and replaces
// those already existing in the older header.
func updateHeader(older, newer http.Header) {
	for name, values := range newer {
		for i, value := range values {
			if i == 0 {
				older.Set(name, value)
			} else {
				older.Add(name, value)
			}
		}
	}
}

func bytesToUint16(b []byte) uint16 {
	return (uint16(b[0]) << 8) + uint16(b[1])
}

func bytesToUint24(b []byte) uint32 {
	return (uint32(b[0]) << 16) + (uint32(b[1]) << 8) + uint32(b[2])
}

func bytesToUint24Reverse(b []byte) uint32 {
	return (uint32(b[2]) << 16) + (uint32(b[1]) << 8) + uint32(b[0])
}

func bytesToUint32(b []byte) uint32 {
	return (uint32(b[0]) << 24) + (uint32(b[1]) << 16) + (uint32(b[2]) << 8) + uint32(b[3])
}

// read is used to ensure that the given number of bytes
// are read if possible, even if multiple calls to Read
// are required.
func read(r io.Reader, i int) ([]byte, error) {
	out := make([]byte, i)
	in := out[:]
	for i > 0 {
		if r == nil {
			return nil, ErrConnNil
		}
		if n, err := r.Read(in); err != nil {
			return nil, err
		} else {
			in = in[n:]
			i -= n
		}
	}
	return out, nil
}

// write is used to ensure that the given data is written
// if possible, even if multiple calls to Write are
// required.
func write(w io.Writer, data []byte) error {
	i := len(data)
	for i > 0 {
		if w == nil {
			return ErrConnNil
		}
		if n, err := w.Write(data); err != nil {
			return err
		} else {
			data = data[n:]
			i -= n
		}
	}
	return nil
}

// readCloser is a helper structure to allow
// an io.Reader to satisfy the io.ReadCloser
// interface.
type readCloser struct {
	io.Reader
}

func (r *readCloser) Close() error {
	return nil
}