package gee

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

// 初始化gee.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 注册路由
func (engine *Engine) addRouter(method string, path string, handler HandlerFunc) {
	// 优化字符串拼接
	var key strings.Builder
	key.WriteString(method)
	key.WriteString("-")
	key.WriteString(path)

	engine.router[key.String()] = handler
}

func (engine *Engine) Get(path string, handler HandlerFunc) {
	engine.addRouter("GET", path, handler)
}

func (engine *Engine) Post(path string, handler HandlerFunc) {
	engine.addRouter("POST", path, handler)
}

// ListenAndServe 需求interface{}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var key strings.Builder
	key.WriteString(req.Method)
	key.WriteString("-")
	key.WriteString(req.URL.Path)

	if handler, ok := engine.router[key.String()]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NO FOUND")
	}
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
