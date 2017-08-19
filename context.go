package boa

import "net/http"

const (
	CharsetUTF8 = "charset=utf-8"
	TextPlain = "text/plain"
	TextPlainCharsetUTF8 = TextPlain + "; " + CharsetUTF8
)

type Context struct {
	Req *http.Request
	Res *Response
	boa *Boa
	routeName string
	pathNames []string
	pathValues []string
	hi int
	handlers []HandlerFunc
}

func NewContext(w http.ResponseWriter, r *http.Request, b *Boa) *Context {
	c := new(Context)
	c.Res = NewResponse(w, b)
	c.boa = b
	c.pathNames = make([]string, 0, 16)
	c.pathValues = make([]string, 0, 16)
	c.handlers = make([]HandlerFunc, len(b.middleware), len(b.middleware) + 3)
	copy(c.handlers, b.middleware)
	c.Reset(w, r)
	return c
}

func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.Req = r
	c.Res.rest(w)
	c.handlers = c.handlers[:len(c.boa.middleware)]
	c.routeName = ""
	c.pathNames = c.pathNames[:0]
	c.pathValues = c.pathNames[:0]
	c.hi = 0
}

func (c *Context) Next() {
	if c.hi >= len(c.handlers) {
		return
	}

	i := c.hi
	c.hi++
	if c.handlers[i] != nil {
		c.handlers[i](c)
	} else {
		c.Next()
	}
}

func (c *Context) String(code int, s string)  {
	c.Res.Header().Set("Content-Type", TextPlainCharsetUTF8)
	c.Res.WriterHeader(code)
	c.Res.Write([]byte(s))
}

