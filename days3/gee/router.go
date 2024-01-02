package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots   map[string]*node
	handler map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:   make(map[string]*node),
		handler: make(map[string]HandlerFunc),
	}
}

// 分割url
func parsePath(path string) []string {
	vs := strings.Split(path, "/")

	parts := make([]string, 0)

	//防止 home//view 情况出现
	//即将 home//view 视为home/view
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (r *router) addRoute(method string, path string, handler HandlerFunc) {
	parts := parsePath(path)

	log.Printf("Route %4s - %4s", method, path)
	key := method + "-" + path
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
	r.handler[key] = handler
}

// 得到对应处理方法路径&&获得动态参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchPath := parsePath(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchPath, 0)

	if n != nil {
		parts := parsePath(n.path)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchPath[index]
			}

			//如果匹配上* 并且带有名字 则后面路径全为该参数
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchPath[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.path
		r.handler[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND")
	}
}
