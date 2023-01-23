package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break;
			}
		}
	}
	return parts
}

func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	parts := parsePattern(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(parts, 0)
	if n != nil {
		params := make(map[string]string)
		getpart := parsePattern(n.pattern)
		for index, part := range getpart {
			if part[0] == ':' {
				params[part[1:]] = parts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(parts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *Router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	}else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}

func (r *Router)getRoutes(method string) []*node{
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}
