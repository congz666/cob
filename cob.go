package cob

import (
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by cob
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router // Implement prefix tree and storage HandlerFunc
	// Group
	*RouterGroup                // Group acting on the global
	groups       []*RouterGroup // storage all groups
}

// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// New is the constructor of cob.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Process requests
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
