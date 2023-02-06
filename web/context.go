package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleWare
	handlers 		[]HandleFunc
	handlerIndex 	int
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		handlerIndex: -1,
	}
}

// Parser params from request

// PostForm returns the first value for the named component of the query.
// 返回请求体中的key对应的值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query parses RawQuery and returns the corresponding value
// 返回请求url中key对应的请求参数值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Param parses request path and returns the value of corresponding path key
// Param 解析路径中的参数，并返回key对应的值
func (c *Context) Param(key string) string {
	v := c.Params[key]
	return v
}

// Construct response info

// Status sends an HTTP response header with the provided status code.
// 返回体中写入响应码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader sets the header entries associated with key to the single element value.
// 响应头赋值
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String format string http reply
// 构建string格式的http返回
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON format json http reply
// 构建json格式的http返回
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Next() {
	c.handlerIndex++
	// can't make sure middlerWares will call Context.Next(), Using c.handlers[c.handlerIndex] is not right
	s := len(c.handlers)
	for ; c.handlerIndex < s; c.handlerIndex++ {
		c.handlers[c.handlerIndex](c)
	}
}