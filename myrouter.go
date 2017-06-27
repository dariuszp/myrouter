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
func NewMyRouter(schema string, host string, port int) *MyRouter {
	var verbs = make(map[string]map[string]*Route)
	for _, verb := range SupportedMethods {
		verbs[verb] = make(map[string]*Route)
	}
	var router = &MyRouter{schema, host, port, verbs, make(map[string]*Route)}
	return router
}

// MyRouter is just my router :-D
type MyRouter struct {
	defaultSchema string
	defaultHost   string
	defaultPort   int
	verbs         map[string]map[string]*Route
	routes        map[string]*Route
}

// AddRoute register method for verbs
// name - name of the route
// methods - list of methods that works with this route
// path - path after the host and port
// requirements - map of regexp patterns (as strings) for route params
func (router *MyRouter) AddRoute(name string, methods []string, path string, requirements map[string]string) (bool, error) {
	return router.AddCustomRoute(name, methods, router.defaultSchema, "", "", router.defaultHost, router.defaultPort, path, requirements)
}

// AddCustomRoute register method for verbs
// name - name of the route
// methods - list of methods that works with this route
// schema - http, https, ftp etc...
// host - website host, for example example.com
// port - leave empty if You don't want to change port
// path - path after the host and port
// requirements - map of regexp patterns (as strings) for route params
func (router *MyRouter) AddCustomRoute(name string, methods []string, schema string, unsecureLogin string, unsecurePassword string, host string, port int, path string, requirements map[string]string) (bool, error) {
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

	route, err = NewRoute(name, methods, schema, host, port, path, requirements)

	if err != nil {
		return false, err
	}

	route.unsecureLogin = unsecureLogin
	route.unsecurePassword = unsecurePassword

	for _, method := range methods {
		method = strings.ToLower(method)
		router.verbs[method][name] = route
	}
	router.routes[name] = route

	return true, nil
}

// RemoveRoute remove route by name
func (router *MyRouter) RemoveRoute(name string) bool {
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

// MatchPath find route that match specified path
func (router *MyRouter) MatchPath(path string) (*Route, map[string]string, error) {
	return match(router, router.routes, path)
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

// GetRouteByName return route by its name
func (router *MyRouter) GetRouteByName(name string) *Route {
	var route, ok = router.routes[name]
	if !ok {
		return nil
	}
	return route
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

// GeneratePath will get route by name and generate path for it
func (router *MyRouter) GeneratePath(name string, parameters map[string]string) (string, error) {
	var route = router.GetRouteByName(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.GeneratePath(parameters)
}

// GeneratePathWithFragment will get route by name and generate path for it with anchor
func (router *MyRouter) GeneratePathWithFragment(name string, parameters map[string]string, fragment string) (string, error) {
	var route = router.GetRouteByName(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.GenerateURLWithFragment(parameters, fragment)
}

// GenerateURL will get route by name and generate url for it
func (router *MyRouter) GenerateURL(name string, parameters map[string]string) (string, error) {
	var route = router.GetRouteByName(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.GenerateURL(parameters)
}

// GenerateURLWithFragment will get route by name and generate url for it
func (router *MyRouter) GenerateURLWithFragment(name string, parameters map[string]string, fragment string) (string, error) {
	var route = router.GetRouteByName(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.GenerateURLWithFragment(parameters, fragment)
}

// GenerateUnsecureURL will get route by name and generate url for it
func (router *MyRouter) GenerateUnsecureURL(name string, login string, password string, parameters map[string]string) (string, error) {
	var route = router.GetRouteByName(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.GenerateUnsecureURL(login, password, parameters)
}

// GenerateUnsecureURLWithFragment will get route by name and generate url for it
func (router *MyRouter) GenerateUnsecureURLWithFragment(name string, login string, password string, parameters map[string]string, fragment string) (string, error) {
	var route = router.GetRouteByName(name)
	if route == nil {
		return "", errors.New(strings.Join([]string{"Invalid route name:", name}, " "))
	}
	return route.GenerateUnsecureURLWithFragment(login, password, parameters, fragment)
}
