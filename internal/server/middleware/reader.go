package middleware

import (
	"compress/gzip"
	"io"
)

type Reader struct {
	rc io.ReadCloser
	zr *gzip.Reader
}

func NewReader(reader io.ReadCloser) (*Reader, error) {
	zipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	return &Reader{
		rc: reader,
		zr: zipReader,
	}, nil
}

func (r *Reader) Read(p []byte) (int, error) {
	return r.zr.Read(p)
}

func (r *Reader) Close() error {
	if err := r.rc.Close(); err != nil {
		return err
	}

	return r.zr.Close()
}
