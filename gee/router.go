package gee

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// define another name of handler

type Router struct {
	roots    map[string]*node
	handlers map[string]Handler
}

func newRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]Handler),
	}
}

// parts stop at *
func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	arr := strings.Split(pattern, "/")
	for _, part := range arr {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) addRoute(method string, path string, f Handler) {
	root, ok := r.roots[method]
	if !ok {
		root = &node{
			part:     "/",
			children: make([]*node, 0),
			isWild:   false,
		}
		r.roots[method] = root
	}
	parts := parsePattern(path)
	root.insert(path, parts, 0)
	s := method + "-" + path
	r.handlers[s] = f
}

func (r *Router) parsePatternMap(n *node, path string) map[string]string {
	m := make(map[string]string)
	parts := parsePattern(n.pattern)
	arr := parsePattern(path)
	//originParts := r.parsePattern(n.pattern)
	for index, part := range parts {
		if part[0] == ':' {
			m[part[1:]] = arr[index]
		} else if part[0] == '*' && len(part) > 1 {
			m[part[1:]] = strings.Join(arr[index:], "/")
			break
		}
	}
	return m
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)

	if n != nil {
		m := r.parsePatternMap(n, path)
		return n, m
	}
	return nil, nil
}

func (r *Router) deleteRoute(method string, path string, f Handler) {
	root, ok := r.roots[method]
	if !ok {
		return
	}
	arr := strings.Split(path, "/")
	delNode := root.search(arr, len(arr))
	if delNode == nil {
		return
	}
	s := method + "-" + path
	delete(r.handlers, s)
}

func (r *Router) handle(c *Context) {

	n, m := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = m
		s := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[s])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "url not found method:%s path:%s", c.Method, c.Path)
		})
	}
	c.Next()
}

// unit test
func NewTestRouter() *Router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello", nil)
	r.addRoute("GET", "/max/:name", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	r.addRoute("POST", "/", nil)
	r.addRoute("POST", "/hello", nil)
	r.addRoute("POST", "/max/:name", nil)
	r.addRoute("POST", "/hello/:name", nil)
	r.addRoute("POST", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern fail")
	}
}

func TestGetRoute(t *testing.T) {
	r := NewTestRouter()
	n, ps := r.getRoute("GET", "/hello/oufan")
	if n == nil {
		t.Fatal("test getRoute fail: nil should not be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("test getRoute fail: n.pattern wrong")
	}

	if ps["name"] != "oufan" {
		t.Fatal("name should be oufan")
	}

	fmt.Printf("matched path:%s, params['name']:%s\n", n.pattern, ps["name"])
}
