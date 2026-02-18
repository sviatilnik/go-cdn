package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// If request content compressed
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			r.Body = gzReader
		}

		// If client doesn't accept gzip encoded content
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gzWriter.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzResponseWriter{
			ResponseWriter: w,
			Writer:         gzWriter,
		}, r)
	})
}

type gzResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
