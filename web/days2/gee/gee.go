package gee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	router *router
}

// 初始化gee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 注册路由
func (engine *Engine) addRouter(method string, path string, handler HandlerFunc) {
	engine.router.addRouter(method, path, handler)
}

func (engine *Engine) GET(path string, handler HandlerFunc) {
	engine.addRouter("GET", path, handler)
}

func (engine *Engine) POST(path string, handler HandlerFunc) {
	engine.addRouter("POST", path, handler)
}

// ListenAndServe 需求interface{}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
