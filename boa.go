package boa

import (
	"net/http"
	"sync"
	"log"
	"os"
)

const (
	Version = "0.0.1"
	DEV     = "development"
	PROD    = "production"
	TEST    = "test"
)

var Env string

type Boa struct {
	debug           bool
	name            string
	router          *Router
	logger          Logger
	pool            sync.Pool
	errorHandler    ErrorHandleFunc
	middleware      []HandlerFunc
	notFoundHandler HandlerFunc
}

type Middleware interface{}

type HandlerFunc func(*Context)

type ErrorHandleFunc func(error, *Context)

func New() *Boa {
	b := new(Boa)
	b.middleware = make([]HandlerFunc, 0)
	b.pool = sync.Pool{
		New: func() interface{} {
			return NewContext(nil, nil, b)
		},
	}

	if Env != PROD {
		b.debug = true
	}

	b.router = NewRouter(b)
	b.logger = log.New(os.Stderr, "[Boa] ", log.LstdFlags)
	b.notFoundHandler = func(ctx *Context) {
		ctx.String(404, "Not found")
	}
	return b
}

func (b *Boa) Server(addr string) *http.Server {
	s := &http.Server{Addr: addr}
	return s
}

func (b *Boa) Run(addr string) {
	b.run(b.Server(addr))
}

func (b *Boa) run(s *http.Server) {
	s.Handler = b
	println("Listen", s.Addr)
	s.ListenAndServe()
}

func (b *Boa) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := b.pool.Get().(*Context)
	defer b.pool.Put(c)
	c.Reset(w, r)
	h, name := b.router.Match(r.Method, r.URL.Path, c)
	c.routeName = name
	if h == nil {
		c.handlers = append(c.handlers, b.notFoundHandler)
	} else {
		c.handlers = append(c.handlers, h...)
	}
	c.Next()
}

func (b *Boa) GET(pattern string, h ...HandlerFunc) {
	b.router.Add("GET", pattern, h)
}