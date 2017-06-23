package myrouter

import (
	"errors"
	"strings"
)

const (
	// MethodOptions string identificator
	MethodOptions = "options"
	// MethodGet string identificator
	MethodGet = "get"
	// MethodHead string identificator
	MethodHead = "head"
	// MethodPost string identificator
	MethodPost = "post"
	// MethodPut string identificator
	MethodPut = "put"
	// MethodDelete string identificator
	MethodDelete = "delete"
	// MethodTrace string identificator
	MethodTrace = "trace"
	// MethodConnect string identificator
	MethodConnect = "connect"
)

// SupportedMethods contains all supported HTTP verbs
var SupportedMethods = []string{MethodOptions, MethodGet, MethodHead, MethodPost, MethodPut, MethodDelete, MethodTrace, MethodConnect}

// MyRouter is just my router :-D
type MyRouter struct {
	verbs  map[string]map[string]*Route
	routes map[string]*Route
}

// AddRoute register method for verbs
// name - name of the route
// methods - list of methods that works with this route
// schema - http, https, ftp etc...
// host - website host, for example example.com
// port - leave empty if You don't want to change port
// path - path after the host and port
func (router *MyRouter) AddRoute(name string, methods []string, schema string, host string, port int, path string) (*MyRouter, error) {
	var _, ok = router.routes[name]
	if ok {
		var err = errors.New(strings.Join([]string{"Route name already registered", name}, " "))
		return router, err
	}
	for _, method := range methods {
		method = strings.ToLower(method)
		if !arrayContainsStringNoCase(SupportedMethods, method) {
			var err = errors.New(strings.Join([]string{"Unsupported method", method}, " "))
			return router, err
		}
	}
	var route = &Route{name, methods, schema, host, port, path}
	for _, method := range methods {
		router.verbs[method][name] = route
	}
	router.routes[name] = route
	return router, nil
}

// UpdateRoute register method for verbs
// name - name of the route
// methods - list of methods that works with this route
// schema - http, https, ftp etc...
// host - website host, for example example.com
// port - leave empty if You don't want to change port
// path - path after the host and port
func (router *MyRouter) UpdateRoute(name string, methods []string, schema string, host string, port int, path string) (*MyRouter, error) {
	var route, ok = router.routes[name]
	for _, method := range methods {
		method = strings.ToLower(method)
		if !arrayContainsStringNoCase(SupportedMethods, method) {
			var err = errors.New(strings.Join([]string{"Unsupported method", method}, " "))
			return router, err
		}
	}

	if ok && !arrayCompareString(methods, route.methods) {
		for _, method := range router.routes[name].methods {
			delete(router.verbs[method], name)
			if len(router.verbs[method]) == 0 {
				delete(router.verbs, method)
			}
		}
		for _, method := range methods {
			router.verbs[method][name] = route
		}
	}

	route.methods = methods
	route.schema = schema
	route.host = host
	route.port = port
	route.path = path

	return router, nil
}

// RemoveRoute remove route by name
func (router *MyRouter) RemoveRoute(name string) bool {
	var _, ok = router.routes[name]
	if !ok {
		return false
	}

	for _, method := range router.routes[name].methods {
		delete(router.verbs[method], name)
		if len(router.verbs[method]) == 0 {
			delete(router.verbs, method)
		}
	}
	delete(router.routes, name)
	return true
}
