package web

import (
	"strings"
)

func parserPattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	paths := make([]string, 0)
	for _, v := range vs {
		// 忽略路由中空项，比如//
		if v != "" {
			paths = append(paths, v)
			// 以'*' 开头匹配项默认匹配url剩余所有路径
			if v[0] == '*' {
				break
			}
		}
	}
	return paths
}

type router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandleFunc), roots: make(map[string]*node)}
}

func (r *router) addRouter(method string, pattern string, handler HandleFunc) {
	paths := parserPattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, paths, 0)
	r.handlers[method+"-"+pattern] = handler
}

// findRouter 寻找匹配路由，将路由中参数键值对返回
func (r *router) findRouter(context *Context) (*node, map[string]string) {
	searchPaths := parserPattern(context.Path)
	params := make(map[string]string)
	v, ok := r.roots[context.Method]
	if !ok {
		return nil, nil
	}
	n := v.search(searchPaths, 0)
	if n != nil {
		paths := parserPattern(n.pattern)
		for i, p := range paths {
			if p[0] == ':' {
				params[p[1:]] = searchPaths[i]
			}
			if p[0] == '*' {
				params[p[1:]] = strings.Join(searchPaths[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) Handle(context *Context) {
	n, params := r.findRouter(context)
	if n != nil {
		context.Params = params
		key := context.Method + "-" + n.pattern
		context.handlers = append(context.handlers, r.handlers[key])
	} else {
		context.handlers = append(context.handlers, func(c *Context) {
			c.String(404, "404 NOT FOUND: %s\n", c.Method)
		})
	}
	context.Next()
}
