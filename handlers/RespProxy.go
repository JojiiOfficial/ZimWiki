package handlers

import (
	"io"
	"net/http"

	"github.com/klauspost/compress/gzip"
)

type writerData struct {
	didWrite bool
}

// ResponseProxy proxy response
type ResponseProxy struct {
	http.ResponseWriter
	acceptGzip bool
	writer     io.Writer

	gzipWriter *gzip.Writer
	useGzip    bool

	wd *writerData
}

// NewResponseProxy proxy sending response
func NewResponseProxy(acceptGzip bool, w http.ResponseWriter) *ResponseProxy {
	resProxy := ResponseProxy{
		ResponseWriter: w,
		acceptGzip:     acceptGzip,

		// We need a pointer
		wd: &writerData{
			didWrite: false,
		},
	}

	// Use a gzip writer
	// if client accepts compression
	if acceptGzip {
		resProxy.gzipWriter = gzip.NewWriter(w)
		resProxy.writer = resProxy.gzipWriter
		resProxy.useGzip = true
	} else {
		resProxy.writer = w
	}

	return &resProxy
}

func (rp ResponseProxy) Write(b []byte) (int, error) {
	if !rp.wd.didWrite {
		// Set Content-Encoding header before writing
		// to the client the first time
		if rp.useGzip {
			rp.Header().Set("Content-Encoding", "gzip")
		}

		rp.wd.didWrite = true
	}

	return rp.writer.Write(b)
}

// Done does all the rest
func (rp ResponseProxy) Done() {
	if rp.useGzip && rp.gzipWriter != nil {
		rp.gzipWriter.Close()
	}
}
