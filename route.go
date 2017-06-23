package myrouter

import (
	"regexp"
	"strings"
)

// Route represent single http route
type Route struct {
	name        string
	methods     []string
	schema      string
	host        string
	port        int
	path        string
	matchRegexp *regexp.Regexp
	parameters  []string
}

// NewRoute create new route
func NewRoute(name string, methods []string, schema string, host string, port int, path string, requirements map[string]string) *Route {
	var regexp = generateRegExpFromPath(path, requirements)
	var parameters = extractParamNames(path)
	var route = &Route{name, methods, schema, host, port, path, regexp, parameters}
	return route
}

// SetMethods replace list of methods
func (route *Route) SetMethods(methods []string) (*Route, bool) {
	if !arrayCompareStringNoCase(SupportedMethods, methods) {
		return route, false
	}
	var result []string
	for _, value := range methods {
		result = append(result, strings.ToLower(value))
	}
	route.methods = result
	return route, true
}

// AddMethod append method to list
func (route *Route) AddMethod(newMethod string) (*Route, bool) {
	newMethod = strings.ToLower(newMethod)
	if !arrayContainsString(SupportedMethods, newMethod) {
		return route, false
	}
	route.methods = append(route.methods, newMethod)
	return route, true
}

//RemoveMethod remove method from route
func (route *Route) RemoveMethod(toRemove string) (*Route, bool) {
	var result []string
	var lenA = len(route.methods)
	toRemove = strings.ToLower(toRemove)
	for _, value := range route.methods {
		if value != toRemove {
			result = append(result, value)
		}
	}
	var lenB = len(result)
	route.methods = result
	return route, lenA != lenB
}
