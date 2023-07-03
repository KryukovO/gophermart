package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

var acceptTypes = []string{"application/json", "text/html"}

type Writer struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func NewWriter(writer http.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: writer,
		zw:             gzip.NewWriter(writer),
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	defer w.Close()

	contentType := w.Header().Get("Content-Type")
	for _, typ := range acceptTypes {
		if strings.Contains(contentType, typ) {
			return w.zw.Write(p)
		}
	}

	return w.ResponseWriter.Write(p)
}

func (w *Writer) WriteHeader(statusCode int) {
	contentType := w.Header().Get("Content-Type")
	for _, typ := range acceptTypes {
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
