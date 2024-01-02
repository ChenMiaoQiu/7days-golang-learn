package gee

import (
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
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
	router        *router
	groups        []*RouterGroup     //所有的路由组
	htmlTemplates *template.Template //html模板渲染
	funcMap       template.FuncMap   //html模板渲染
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

// 注册模板
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPath := path.Join(relativePath, "/*filepath")
	group.GET(urlPath, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absPath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absPath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		//检查文件是否存在
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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
	c.engine = engine
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
