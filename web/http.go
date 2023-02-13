package web

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

type HandleFunc func(*Context)

type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func NewEngine() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// insert middleWares into context
	var middleWares []HandleFunc
	for _, m := range e.groups {
		if strings.HasPrefix(req.URL.Path, m.prefix) {
			middleWares = append(middleWares, m.middleWares...)
		}
	}
	c := NewContext(w, req, e)
	c.handlers = middleWares
	e.router.Handle(c)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

// LoadHTMLGlob loads all template html file into cache
// ParseGlob open all files match the pattern and store all files content
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

type RouterGroup struct {
	prefix      string
	middleWares []HandleFunc
	parent      *RouterGroup
	engine      *Engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	engine := r.engine
	newRouterGroup := &RouterGroup{
		prefix: r.prefix + prefix,
		parent: r,
		engine: engine,
	}
	engine.groups = append(engine.groups, newRouterGroup)
	return newRouterGroup
}

func (r *RouterGroup) addRouter(method string, pattern string, handler HandleFunc) {
	r.engine.router.addRouter(method, r.prefix+pattern, handler)
}

func (r *RouterGroup) Use(middleWare HandleFunc) {
	r.middleWares = append(r.middleWares, middleWare)
}

func (r *RouterGroup) GET(url string, handler HandleFunc) {
	r.addRouter("GET", url, handler)
}

func (r *RouterGroup) POST(url string, handler HandleFunc) {
	r.addRouter("POST", url, handler)
}

// Static register router to server static file server
func (r *RouterGroup) Static(relativePath string, root string) {
	handler := r.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*path")
	r.GET(urlPattern, handler)
}

func (r *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	stripPath := path.Join(r.prefix, relativePath)
	fileServer := http.StripPrefix(stripPath, http.FileServer(fs))
	return func(c *Context) {
		filePath := c.Param("path")
		// check if file exists and/or if we have permission to access it
		if _, err := fs.Open(filePath); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
