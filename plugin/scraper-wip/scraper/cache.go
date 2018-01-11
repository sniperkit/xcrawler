package scraper

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// fileResponseWriter wraps an http.ResponseWriter and a File
// passing it to an http.Handler's ServeHTTP
// will write to both the file and the response.

type fileResponseWriter struct {
	file  io.Writer
	resp  http.ResponseWriter
	multi io.Writer
}

func newFileResponseWriter(file io.Writer, resp http.ResponseWriter) http.ResponseWriter {
	multi := io.MultiWriter(file, resp)
	return &fileResponseWriter{
		file:  file,
		resp:  resp,
		multi: multi,
	}
}

// implement http.ResponseWriter
// https://golang.org/pkg/net/http/#ResponseWriter
func (w *fileResponseWriter) Header() http.Header {
	return w.resp.Header()
}

func (w *fileResponseWriter) Write(b []byte) (int, error) {
	return w.multi.Write(b)
}

func (w *fileResponseWriter) WriteHeader(i int) {
	w.resp.WriteHeader(i)
}

// Example http handler middleware.
// proxies to backend app,
// and then caches the response in file named after request PATH
type RequestCacher struct {
	app      http.Handler
	cacheDir string
}

func NewRequestCacher(app http.Handler, cacheDir string) http.Handler {
	handler := &RequestCacher{
		app:      app,
		cacheDir: cacheDir,
	}

	return handler
}

func (h *RequestCacher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var writer http.ResponseWriter
	writer = rw

	path := filepath.Join(h.cacheDir, req.URL.Path)
	dir := filepath.Dir(path)

	// Make sure the cache dir exists
	err := os.MkdirAll(dir, 0764)
	if err != nil {
		log.Println("mkdir", err)
		h.app.ServeHTTP(writer, req)
		return
	}

	// create file to cache response into.
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Println("create file", err)
		h.app.ServeHTTP(writer, req)
		return
	}

	// wrap file and original response writer
	writer = newFileResponseWriter(file, rw)

	// backend app will write to both
	// file and http response
	h.app.ServeHTTP(writer, req)
	log.Println("cached", path)
}
