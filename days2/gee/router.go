package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handler map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handler: make(map[string]HandlerFunc)}
}

func (r *router) addRouter(method string, path string, handler HandlerFunc) {
	log.Printf("Route %4s - %4s", method, path)
	// key = method + "-" + path 优化
	var key strings.Builder
	key.WriteString(method)
	key.WriteString("-")
	key.WriteString(path)

	r.handler[key.String()] = handler
}

func (r *router) handle(c *Context) {
	var key strings.Builder
	key.WriteString(c.Method)
	key.WriteString("-")
	key.WriteString(c.Path)

	if handler, ok := r.handler[key.String()]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NO FOUND")
	}
}
