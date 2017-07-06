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
	Name         string
	Methods      []string
	Scheme       string
	UnsecureUser string
	Host         string
	Port         int
	Path         string
	Parameters   []string
	MatchRegexp  *regexp.Regexp
	Requirements map[string]*regexp.Regexp
}

// NewRoute create new route
func NewRoute(name string, methods []string, scheme string, host string, port int, path string, requirements Requirements) (*Route, error) {
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
	route.Methods = result
	return route, true
}

// Match path to find route
// @returns boolean
func (route *Route) Match(path string) bool {
	return route.MatchRegexp.MatchString(path)
}

// MatchURL path to find route
// @returns boolean
func (route *Route) MatchURL(urlAddress string) bool {
	var parsed, err = url.Parse(urlAddress)
	if err != nil {
		return false
	}
	if !strings.EqualFold(parsed.Scheme, route.Scheme) {
		return false
	}
	var host = route.Host
	if route.Port > 0 {
		host = strings.Join([]string{host, strconv.Itoa(route.Port)}, ":")
	}
	if !strings.EqualFold(parsed.Host, host) {
		return false
	}
	return route.MatchRegexp.MatchString(parsed.Path)
}

// MatchMethod match route against specific method
// If no method is specified in route, this check will pass
func (route *Route) MatchMethod(method string, path string) bool {
	if len(route.Methods) > 0 && !arrayContainsString(route.Methods, method) {
		return false
	}
	return route.Match(path)
}

// MatchURLMethod match route against specific method
// If no method is specified in route, this check will pass
func (route *Route) MatchURLMethod(method string, urlAddress string) bool {
	if len(route.Methods) > 0 && !arrayContainsString(route.Methods, method) {
		return false
	}
	var parsed, err = url.Parse(urlAddress)
	if err != nil {
		return false
	}
	if !strings.EqualFold(parsed.Scheme, route.Scheme) {
		return false
	}
	var host = route.Host
	if route.Port > 0 {
		host = strings.Join([]string{host, strconv.Itoa(route.Port)}, ":")
	}
	if !strings.EqualFold(parsed.Host, host) {
		return false
	}
	return route.Match(parsed.Path)
}

// ParsePath will parse path, check for a match and then extract route parameters
// Returns match bool, parameetsrs and error in case parameters parse does not work
func (route *Route) ParsePath(path string) (*MyURL, error) {
	var parameters PathParameters

	// Parse url
	var parsed, err = url.Parse(path)
	if err != nil {
		return NewEmptyMyURL(), err
	}

	path = parsed.Path
	var match = route.Match(path)
	if !match {
		return NewEmptyMyURL(), errors.New("No route match")
	}

	parameters, err = extractParamsFromRoute(route, path)

	var queryParameters URLParameters
	for key, value := range parsed.Query() {
		queryParameters[key] = value
	}

	return NewMyURL("", "", "", "", 0, path, parameters, queryParameters, parsed.Fragment, route), nil
}

// ParseURL will parse path, check for a match and then extract route parameters
// Returns match bool, parameetsrs and error in case parameters parse does not work
func (route *Route) ParseURL(urlAddress string) (*MyURL, error) {
	var parameters PathParameters

	// Parse url
	var parsed, err = url.Parse(urlAddress)
	if err != nil {
		return NewEmptyMyURL(), err
	}

	var path = parsed.Path
	var match = route.Match(path)
	if !match {
		return NewEmptyMyURL(), errors.New("No route match")
	}

	// Extrac host and port
	var parsedHost = parsed.Host
	var splitParsedHost = strings.Split(parsedHost, ":")
	var onlyHost = ""
	var onlyPort = 0
	if len(splitParsedHost) == 2 {
		onlyHost = splitParsedHost[0]
		onlyPort, err = strconv.Atoi(splitParsedHost[1])
		if err != nil {
			onlyPort = 0
		}
	} else {
		onlyHost = parsedHost
	}

	var username = ""
	var password = ""

	if parsed.User != nil {
		var ok bool
		username = parsed.User.Username()
		password, ok = parsed.User.Password()
		if !ok {
			password = ""
		}
	}

	parameters, err = extractParamsFromRoute(route, path)

	var queryParameters = make(URLParameters)
	for key, value := range parsed.Query() {
		queryParameters[key] = value
	}

	return NewMyURL(parsed.Scheme, username, password, onlyHost, onlyPort, parsed.Path, parameters, queryParameters, parsed.Fragment, route), nil
}

// GeneratePath generate path from route
func (route *Route) GeneratePath(parameters URLParameters) (string, error) {
	return generatePath(route.Path, parameters, route.Requirements)
}

// GeneratePathWithFragment generate path from route and add anchor at the end
func (route *Route) GeneratePathWithFragment(parameters URLParameters, fragment string) (string, error) {
	var path, err = generatePath(route.Path, parameters, route.Requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{path, fragment}, "#"), nil
}

// GenerateURL generate path from route
func (route *Route) GenerateURL(parameters URLParameters) (string, error) {
	return generateURL(route.Scheme, route.UnsecureUser, route.Host, route.Port, route.Path, parameters, route.Requirements)
}

// GenerateURLWithFragment generate path from route and add anchor at the end
func (route *Route) GenerateURLWithFragment(parameters URLParameters, fragment string) (string, error) {
	var url, err = generateURL(route.Scheme, route.UnsecureUser, route.Host, route.Port, route.Path, parameters, route.Requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{url, fragment}, "#"), nil
}

// UnsecureURL generate path from route with user
func (route *Route) UnsecureURL(user string, parameters URLParameters) (string, error) {
	return generateURL(route.Scheme, user, route.Host, route.Port, route.Path, parameters, route.Requirements)
}

// UnsecureURLWithFragment generate path from route and add anchor at the end with user
func (route *Route) UnsecureURLWithFragment(user string, parameters URLParameters, fragment string) (string, error) {
	var url, err = generateURL(route.Scheme, user, route.Host, route.Port, route.Path, parameters, route.Requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{url, fragment}, "#"), nil
}
