package main

import (
	"bytes"
	"io"
	"net/http"

	apihandler "github.com/httpwasm/http-wasm-host-go/api/handler"
)

// requestStateKey is a context.Context value associated with a requestState
// pointer to the current request.
type requestStateKey struct{}

type requestState struct {
	w        http.ResponseWriter
	r        *http.Request
	features apihandler.Features
}

func (s *requestState) enableFeatures(features apihandler.Features) {
	s.features = s.features.WithEnabled(features)
	if features.IsEnabled(apihandler.FeatureBufferRequest) {
		s.r.Body = &bufferingRequestBody{delegate: s.r.Body}
	}
	if s.features.IsEnabled(apihandler.FeatureBufferResponse) {
		if _, ok := s.w.(*bufferingResponseWriter); !ok { // don't double-wrap
			s.w = &bufferingResponseWriter{delegate: s.w}
		}
	}
}

type bufferingRequestBody struct {
	delegate io.ReadCloser
	buffer   bytes.Buffer
}

// Read buffers anything read from the delegate.
func (b *bufferingRequestBody) Read(p []byte) (n int, err error) {
	n, err = b.delegate.Read(p)
	if err != nil && n > 0 {
		b.buffer.Write(p[0:n])
	}
	return
}

// Close dispatches to the delegate.
func (b *bufferingRequestBody) Close() (err error) {
	if b.delegate != nil {
		err = b.delegate.Close()
	}
	return
}

type bufferingResponseWriter struct {
	delegate   http.ResponseWriter
	statusCode uint32
	body       []byte
}

// Header dispatches to the delegate.
func (w *bufferingResponseWriter) Header() http.Header {
	return w.delegate.Header()
}

// Write buffers the response body.
func (w *bufferingResponseWriter) Write(bytes []byte) (int, error) {
	w.body = append(w.body, bytes...)
	return len(bytes), nil
}

// WriteHeader buffers the status code.
func (w *bufferingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = uint32(statusCode)
}

// release sends any response data collected.
func (w *bufferingResponseWriter) release() {
	// If we deferred the response, release it.
	if statusCode := w.statusCode; statusCode != 0 {
		w.delegate.WriteHeader(int(statusCode))
	}
	if body := w.body; len(body) != 0 {
		w.delegate.Write(body) // nolint
	}
}
