package myrouter

import (
	"strings"
)

// Route represent single http route
type Route struct {
	name    string
	methods []string
	schema  string
	host    string
	port    int
	path    string
	handler func()
}

// SetMethods replace list of methods
func (route *Route) SetMethods(methods []string) bool {
	if !stringCompareNoCaseArray(SupportedMethods, methods) {
		return false
	}
	route.methods = methods
	return true
}

// AddMethod append method to list
func (route *Route) AddMethod(newMethod string) bool {
	newMethod = strings.ToLower(newMethod)
	if !stringInArray(SupportedMethods, newMethod) {
		return false
	}
	route.methods = append(route.methods, newMethod)
	return true
}
