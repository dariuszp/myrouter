package myrouter

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Route represent single http route
type Route struct {
	name         string
	methods      []string
	scheme       string
	unsecureUser string
	host         string
	port         int
	path         string
	parameters   []string
	matchRegexp  *regexp.Regexp
	requirements map[string]*regexp.Regexp
}

// NewRoute create new route
func NewRoute(name string, methods []string, scheme string, host string, port int, path string, requirements map[string]string) (*Route, error) {
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
	return &Route{name, methods, scheme, "", host, port, path, parameters, regexpFromPath, requirementsRegexp}, nil
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

// MatchURL path to find route
// @returns boolean
func (route *Route) MatchURL(urlAddress string) bool {
	var parsed, err = url.Parse(urlAddress)
	if err != nil {
		return false
	}
	if !strings.EqualFold(parsed.Scheme, route.scheme) {
		return false
	}
	var host = route.host
	if route.port > 0 {
		host = strings.Join([]string{host, strconv.Itoa(route.port)}, ":")
	}
	if !strings.EqualFold(parsed.Host, host) {
		return false
	}
	return route.matchRegexp.MatchString(parsed.Path)
}

// MatchMethod match route against specific method
// If no method is specified in route, this check will pass
func (route *Route) MatchMethod(method string, path string) bool {
	if len(route.methods) > 0 && !arrayContainsString(route.methods, method) {
		return false
	}
	return route.Match(path)
}

// MatchURLMethod match route against specific method
// If no method is specified in route, this check will pass
func (route *Route) MatchURLMethod(method string, urlAddress string) bool {
	if len(route.methods) > 0 && !arrayContainsString(route.methods, method) {
		return false
	}
	var parsed, err = url.Parse(urlAddress)
	if err != nil {
		return false
	}
	if !strings.EqualFold(parsed.Scheme, route.scheme) {
		return false
	}
	var host = route.host
	if route.port > 0 {
		host = strings.Join([]string{host, strconv.Itoa(route.port)}, ":")
	}
	if !strings.EqualFold(parsed.Host, host) {
		return false
	}
	return route.Match(parsed.Path)
}

// ParsePath will parse path, check for a match and then extract route parameters
// Returns match bool, parameetsrs and error in case parameters parse does not work
func (route *Route) ParsePath(path string) (bool, map[string]string, error) {
	var match = route.Match(path)
	var parameters, err = extractParamsFromRoute(route, path)
	return match, parameters, err
}

// ParseURL will parse path, check for a match and then extract route parameters
// Returns match bool, parameetsrs and error in case parameters parse does not work
func (route *Route) ParseURL(urlAddress string) (bool, map[string]string, error) {
	var parameters map[string]string
	var parsed, err = url.Parse(urlAddress)
	if err != nil {
		return false, make(map[string]string), err
	}
	var path = parsed.Path
	var match = route.Match(path)
	parameters, err = extractParamsFromRoute(route, path)
	return match, parameters, err
}

// Path generate path from route
func (route *Route) Path(parameters map[string]string) (string, error) {
	return generatePath(route.path, parameters, route.requirements)
}

// PathWithFragment generate path from route and add anchor at the end
func (route *Route) PathWithFragment(parameters map[string]string, fragment string) (string, error) {
	var path, err = generatePath(route.path, parameters, route.requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{path, fragment}, "#"), nil
}

// URL generate path from route
func (route *Route) URL(parameters map[string]string) (string, error) {
	return generateURL(route.scheme, route.unsecureUser, route.host, route.port, route.path, parameters, route.requirements)
}

// URLWithFragment generate path from route and add anchor at the end
func (route *Route) URLWithFragment(parameters map[string]string, fragment string) (string, error) {
	var url, err = generateURL(route.scheme, route.unsecureUser, route.host, route.port, route.path, parameters, route.requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{url, fragment}, "#"), nil
}

// UnsecureURL generate path from route with user
func (route *Route) UnsecureURL(user string, parameters map[string]string) (string, error) {
	return generateURL(route.scheme, user, route.host, route.port, route.path, parameters, route.requirements)
}

// UnsecureURLWithFragment generate path from route and add anchor at the end with user
func (route *Route) UnsecureURLWithFragment(user string, parameters map[string]string, fragment string) (string, error) {
	var url, err = generateURL(route.scheme, user, route.host, route.port, route.path, parameters, route.requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{url, fragment}, "#"), nil
}
