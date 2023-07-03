package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type Writer struct {
	http.ResponseWriter
	zw          *gzip.Writer
	acceptTypes []string
}

func NewWriter(writer http.ResponseWriter, acceptTypes []string) *Writer {
	return &Writer{
		ResponseWriter: writer,
		zw:             gzip.NewWriter(writer),
		acceptTypes:    acceptTypes,
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	defer w.Close()

	contentType := w.Header().Get("Content-Type")
	for _, typ := range w.acceptTypes {
		if strings.Contains(contentType, typ) {
			return w.zw.Write(p)
		}
	}

	return w.ResponseWriter.Write(p)
}

func (w *Writer) WriteHeader(statusCode int) {
	contentType := w.Header().Get("Content-Type")
	for _, typ := range w.acceptTypes {
		if strings.Contains(contentType, typ) {
			w.Header().Set("Content-Encoding", "gzip")

			break
		}
	}

	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *Writer) Close() error {
	return w.zw.Close()
}
