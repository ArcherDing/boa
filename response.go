package boa

import (
	"net/http"
	"io"
)

type Response struct {
	status int
	res http.ResponseWriter
	writer io.Writer
	written int64
	wroteHeader bool
	boa *Boa
}

func NewResponse(w http.ResponseWriter, b *Boa) *Response {
	r := new(Response)
	r.res = w
	r.writer = w
	r.boa = b
	return r
}

func (r *Response) rest(w http.ResponseWriter) {
	r.res = w
	r.writer = w
	r.status = http.StatusOK
}

func (r *Response) Header() http.Header {
	return r.res.Header()
}

func (r *Response) WriterHeader(code int) {
	if r.wroteHeader {
		return
	}
	r.wroteHeader = true
	r.status = code
	r.res.WriteHeader(code)
}

func (r *Response) Write(bytes []byte) (int, error) {
	if !r.wroteHeader {
		r.res.WriteHeader(http.StatusOK)
	}
	n, err := r.writer.Write(bytes)
	r.written += int64(n)
	return n, err
}
