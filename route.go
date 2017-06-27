package myrouter

import (
	"errors"
	"regexp"
	"strings"
)

// Route represent single http route
type Route struct {
	name             string
	methods          []string
	schema           string
	unsecureLogin    string
	unsecurePassword string
	host             string
	port             int
	path             string
	parameters       []string
	matchRegexp      *regexp.Regexp
	requirements     map[string]*regexp.Regexp
}

// NewRoute create new route
func NewRoute(name string, methods []string, schema string, host string, port int, path string, requirements map[string]string) (*Route, error) {
	if len(path) == 0 {
		return nil, errors.New("Route path cannot be empty")
	}

	var regexpFromPath, rx *regexp.Regexp
	var err error
	var requirementsRegexp = make(map[string]*regexp.Regexp)

	for name, pattern := range requirements {
		rx, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		requirementsRegexp[name] = rx
	}

	regexpFromPath, err = generateRegexpFromPath(path, requirementsRegexp)
	if err != nil {
		return nil, err
	}

	var parameters = extractParamNames(path)
	return &Route{name, methods, schema, "", "", host, port, path, parameters, regexpFromPath, requirementsRegexp}, nil
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

// GeneratePathWithFragment generate path from route and add anchor at the end
func (route *Route) GeneratePathWithFragment(parameters map[string]string, fragment string) (string, error) {
	var path, err = generatePath(route.path, parameters, route.requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{path, fragment}, "#"), nil
}

// GenerateURL generate path from route
func (route *Route) GenerateURL(parameters map[string]string) (string, error) {
	return generateURL(route.schema, route.unsecureLogin, route.unsecurePassword, route.host, route.port, route.path, parameters, route.requirements)
}

// GenerateURLWithFragment generate path from route and add anchor at the end
func (route *Route) GenerateURLWithFragment(parameters map[string]string, fragment string) (string, error) {
	var url, err = generateURL(route.schema, route.unsecureLogin, route.unsecurePassword, route.host, route.port, route.path, parameters, route.requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{url, fragment}, "#"), nil
}

// GenerateUnsecureURL generate path from route with login and password
func (route *Route) GenerateUnsecureURL(login string, password string, parameters map[string]string) (string, error) {
	return generateURL(route.schema, login, password, route.host, route.port, route.path, parameters, route.requirements)
}

// GenerateUnsecureURLWithFragment generate path from route and add anchor at the end with login and password
func (route *Route) GenerateUnsecureURLWithFragment(login string, password string, parameters map[string]string, fragment string) (string, error) {
	var url, err = generateURL(route.schema, login, password, route.host, route.port, route.path, parameters, route.requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{url, fragment}, "#"), nil
}
