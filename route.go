package myrouter

import (
	"errors"
	"regexp"
	"strings"
)

// Route represent single http route
type Route struct {
	name         string
	methods      []string
	schema       string
	host         string
	port         int
	path         string
	parameters   []string
	matchRegexp  *regexp.Regexp
	requirements map[string]*regexp.Regexp
}

// NewRoute create new route
func NewRoute(name string, methods []string, schema string, host string, port int, path string, requirements map[string]string) (*Route, error) {
	var regexpFromPath, err = generateRegExpFromPath(path, requirements)
	if err != nil {
		return nil, err
	}

	var parameters = extractParamNames(path)
	var regexpRequirements = make(map[string]*regexp.Regexp)

	for name, requirement := range requirements {
		regexpRequirements[name], err = regexp.Compile(requirement)
		if err != nil {
			err = errors.New(strings.Join([]string{"Cannot compile requirement:", name, "=", requirement}, " "))
			return nil, err
		}
	}
	return &Route{name, methods, schema, host, port, path, parameters, regexpFromPath, regexpRequirements}, nil
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

// Match path to find route
// @returns boolean
func (route *Route) Match(path string) bool {
	return route.matchRegexp.MatchString(path)
}

// MatchMethod match route against specific method
// If no method is specified in route, this check will pass
func (route *Route) MatchMethod(method string, path string) bool {
	if len(route.methods) > 0 && !arrayContainsString(route.methods, method) {
		return false
	}
	return route.Match(path)
}

// ParsePath will parse path, check for a match and then extract route parameters
// Returns match bool, parameetsrs and error in case parameters parse does not work
func (route *Route) ParsePath(path string) (bool, map[string]string, error) {
	var match = route.Match(path)
	var parameters, err = extractParamsFromRoute(route, path)
	return match, parameters, err
}

// GeneratePath generate path from route
func (route *Route) GeneratePath(parameters map[string]string) (string, error) {
	return generatePath(route.path, parameters, route.requirements)
}

// GenerateURL generate path from route
func (route *Route) GenerateURL(parameters map[string]string) (string, error) {
	return generateURL(route.schema, route.host, route.port, route.path, parameters, route.requirements)
}
