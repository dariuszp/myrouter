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

//NewMyRouter create instance of MyRouter
func NewMyRouter(scheme string, host string, port int) *MyRouter {
	var verbs = make(map[string]map[string]*Route)
	for _, verb := range SupportedMethods {
		verbs[verb] = make(map[string]*Route)
	}
	var router = &MyRouter{scheme, "", host, port, verbs, make(map[string]*Route)}
	return router
}

//NewUnsecureMyRouter create instance of MyRouter
func NewUnsecureMyRouter(scheme string, user string, host string, port int) *MyRouter {
	var verbs = make(map[string]map[string]*Route)
	for _, verb := range SupportedMethods {
		verbs[verb] = make(map[string]*Route)
	}
	var router = &MyRouter{scheme, user, host, port, verbs, make(map[string]*Route)}
	return router
}

// MyRouter is just my router :-D
type MyRouter struct {
	defaultSchema       string
	defaultUnsecureUser string
	defaultHost         string
	defaultPort         int
	verbs               map[string]map[string]*Route
	routes              map[string]*Route
}

// Add register method for verbs
// name - name of the route
// methods - list of methods that works with this route
// path - path after the host and port
// requirements - map of regexp patterns (as strings) for route params
func (router *MyRouter) Add(name string, methods []string, path string, requirements map[string]string) (bool, error) {
	return router.AddCustom(name, methods, router.defaultSchema, router.defaultUnsecureUser, router.defaultHost, router.defaultPort, path, requirements)
}

// AddCustom register method for verbs
// name - name of the route
// methods - list of methods that works with this route
// scheme - http, https, ftp etc...
// host - website host, for example example.com
// port - leave empty if You don't want to change port
// path - path after the host and port
// requirements - map of regexp patterns (as strings) for route params
func (router *MyRouter) AddCustom(name string, methods []string, scheme string, unsecureUser string, host string, port int, path string, requirements map[string]string) (bool, error) {
	var err error
	var route *Route
	var _, ok = router.routes[name]
	if ok {
		err = errors.New(strings.Join([]string{"Route name already registered", name}, " "))
		return false, err
	}
	for _, method := range methods {
		method = strings.ToLower(method)
		if !arrayContainsStringNoCase(SupportedMethods, method) {
			var err = errors.New(strings.Join([]string{"Unsupported method", method}, " "))
			return false, err
		}
	}

	route, err = NewRoute(name, methods, scheme, host, port, path, requirements)

	if err != nil {
		return false, err
	}

	route.unsecureUser = unsecureUser

	for _, method := range methods {
		method = strings.ToLower(method)
		router.verbs[method][name] = route
	}
	router.routes[name] = route

	return true, nil
}

// Remove remove route by name
func (router *MyRouter) Remove(name string) bool {
	var _, ok = router.routes[name]
	if !ok {
		return false
	}

	for _, method := range router.routes[name].methods {
		method = strings.ToLower(method)
		delete(router.verbs[method], name)
	}
	delete(router.routes, name)
	return true
}

// Get return route by its name
func (router *MyRouter) Get(name string) *Route {
	var route, ok = router.routes[name]
	if !ok {
		return nil
	}
	return route
}

// MatchPath find route that match specified path
func (router *MyRouter) MatchPath(path string) (*Route, map[string]string, error) {
	return match(router, router.routes, path)
}

// MatchURL find route that match specified path
func (router *MyRouter) MatchURL(url string) (*Route, map[string]string, error) {
	return matchURL(router, router.routes, url)
}

// MatchPathByMethod find route that match specified path filtered by method
// Returns route, parameters and error, route will be nil if there is no match
func (router *MyRouter) MatchPathByMethod(method string, path string) (*Route, map[string]string, error) {
	method = strings.ToLower(method)
	var list, ok = router.verbs[method]
	if !ok {
		var err = errors.New(strings.Join([]string{"Method not found", method}, " "))
		return nil, map[string]string{}, err
	}

	return match(router, list, path)
}

// Path will get route by name and generate path for it
func (router *MyRouter) Path(name string, parameters map[string]string) (string, error) {
	var route = router.Get(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.Path(parameters)
}

// PathWithFragment will get route by name and generate path for it with anchor
func (router *MyRouter) PathWithFragment(name string, parameters map[string]string, fragment string) (string, error) {
	var route = router.Get(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.PathWithFragment(parameters, fragment)
}

// URL will get route by name and generate url for it
func (router *MyRouter) URL(name string, parameters map[string]string) (string, error) {
	var route = router.Get(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.URL(parameters)
}

// URLWithFragment will get route by name and generate url for it
func (router *MyRouter) URLWithFragment(name string, parameters map[string]string, fragment string) (string, error) {
	var route = router.Get(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.URLWithFragment(parameters, fragment)
}

// UnsecureURL will get route by name and generate url for it
func (router *MyRouter) UnsecureURL(name string, user string, parameters map[string]string) (string, error) {
	var route = router.Get(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.UnsecureURL(user, parameters)
}

// UnsecureURLWithFragment will get route by name and generate url for it
func (router *MyRouter) UnsecureURLWithFragment(name string, user string, parameters map[string]string, fragment string) (string, error) {
	var route = router.Get(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.UnsecureURLWithFragment(user, parameters, fragment)
}

func match(router *MyRouter, list map[string]*Route, path string) (*Route, map[string]string, error) {
	for _, route := range list {
		if route.Match(path) {
			var _, parameters, err = route.ParsePath(path)
			return route, parameters, err
		}
	}
	return nil, map[string]string{}, errors.New("No route match")
}

func matchURL(router *MyRouter, list map[string]*Route, url string) (*Route, map[string]string, error) {
	for _, route := range list {
		if route.MatchURL(url) {
			var _, parameters, err = route.ParseURL(url)
			return route, parameters, err
		}
	}
	return nil, map[string]string{}, errors.New("No route match")
}
