package gee

import (
	"log"
	"net/http"
	"strings"
)

type Handler func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []Handler
	parent      *RouterGroup
	gee         *Gee //  share one gee by pointer
}

type Gee struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

func New() *Gee {
	gee := &Gee{
		router: newRouter(),
	}
	gee.RouterGroup = &RouterGroup{
		prefix:      "/",
		middlewares: make([]Handler, 0),
		parent:      nil,
		gee:         gee,
	}
	gee.groups = make([]*RouterGroup, 0)
	gee.groups = append(gee.groups, gee.RouterGroup)
	return gee
}

func (group *RouterGroup) Group(pre string) *RouterGroup {
	gee := group.gee
	g := &RouterGroup{
		prefix:      pre,
		middlewares: make([]Handler, 0),
		parent:      group,
		gee:         gee,
	}
	gee.groups = append(gee.groups, g)
	return g
}

func (group *RouterGroup) addRoute(method string, path string, fc Handler) {
	pattern := group.prefix + path
	log.Printf("Route %4s - %s", method, pattern)
	group.gee.router.addRoute(method, pattern, fc)
}

func (group *RouterGroup) GET(path string, fc Handler) {
	group.addRoute("GET", path, fc)
}

func (group *RouterGroup) POST(path string, fc Handler) {
	group.addRoute("POST", path, fc)
}

func (group *RouterGroup) Use(middlewares ...Handler) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (gee *Gee) Run(port string) {

	http.ListenAndServe(port, gee)
}

func (gee *Gee) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []Handler
	for _, group := range gee.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	gee.router.handle(c)
}
