package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// 初始化gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 设置路由组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 设置中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 注册路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	path := group.prefix + comp
	log.Printf("Route %4s - %4s", method, path)
	group.engine.router.addRoute(method, path, handler)
}

// 注册get方法
func (group *RouterGroup) GET(path string, handler HandlerFunc) {
	group.addRoute("GET", path, handler)
}

// 注册post方法
func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.addRoute("POST", path, handler)
}

// ListenAndServe 需求interface{}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//获取对应分组所需的中间件
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

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
