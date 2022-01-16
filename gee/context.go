package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// req and res
	W   http.ResponseWriter
	Req *http.Request

	// req type
	Method string
	Path   string
	Params map[string]string

	// status code
	StatusCode int

	// middleware
	handlers []Handler
	index    int
}

func newContext(writer http.ResponseWriter, reqest *http.Request) *Context {
	return &Context{
		W:      writer,
		Req:    reqest,
		Method: reqest.Method,
		Path:   reqest.URL.Path,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) GetParams(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.W.WriteHeader(code)
}

func (c *Context) SetHeader(k string, v string) {
	c.W.Header().Set(k, v)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.W.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	if _, err := c.W.Write(data); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	c.W.Write([]byte(html))
}

func (c *Context) Fail(code int, message string) {
	//c.Status(code)
	c.index = len(c.handlers)
	c.JSON(code, H{
		"message": message,
	})
}
