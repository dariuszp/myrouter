package myrouter

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var defaultReadParameterPattern = "[^/]+"
var defaultReadParameterRegexp = regexp.MustCompile(defaultReadParameterPattern)

//RegExpFromPath turns path to regexp pattern
func generateRegexpFromPath(path string, requirements map[string]*regexp.Regexp) (*regexp.Regexp, error) {
	var parameterEscapedName, escapedPath, result, patternReplace string
	var ok bool
	var rx *regexp.Regexp
	var parameters = extractParamNames(path)

	escapedPath = regexp.QuoteMeta(path)
	if len(parameters) == 0 {
		return regexp.Compile(escapedPath)
	}

	result = escapedPath
	for _, parameterName := range parameters {
		rx, ok = requirements[parameterName]
		if ok && len(rx.String()) > 0 {
			patternReplace = rx.String()
		} else {
			patternReplace = defaultReadParameterPattern
		}
		parameterEscapedName = strings.Join([]string{"\\{", parameterName, "\\}"}, "")
		patternReplace = strings.Join([]string{"(", patternReplace, ")"}, "")
		result = strings.Replace(result, parameterEscapedName, patternReplace, -1)
	}

	return regexp.Compile(result)
}

// Path fill url pattern with parameters
func generatePath(path string, parameters map[string][]string, requirements map[string]*regexp.Regexp) (string, error) {
	var extracted = extractParamNames(path)
	for _, parameterName := range extracted {
		var values []string
		var value string
		var ok bool

		values, ok = parameters[parameterName]
		if !ok {
			return "", errors.New(strings.Join([]string{"Invalid path parameter", parameterName}, " "))
		}

		if len(values) != 1 {
			return "", errors.New(strings.Join([]string{"Parameters in path must contain only 1 element:", parameterName}, " "))
		}

		value = values[0]

		delete(parameters, parameterName)

		var rx, rxok = requirements[parameterName]
		if !rxok {
			rx = defaultReadParameterRegexp
		}
		if !rx.MatchString(value) {
			return "", errors.New(strings.Join([]string{"Parameter ", parameterName, "does not meet requirements", rx.String()}, " "))
		}

		parameterName = strings.Join([]string{"{", parameterName, "}"}, "")

		value = url.QueryEscape(value)
		path = strings.Replace(path, parameterName, value, -1)
	}

	if len(parameters) > 0 {
		var queryString = generateQueryString(parameters)
		if len(queryString) > 0 {
			path = strings.Join([]string{path, queryString}, "?")
		}
	}

	return path, nil
}

// URL combine host, port and path to create absolute url
func generateURL(scheme string, user string, host string, port int, path string, parameters map[string][]string, requirements map[string]*regexp.Regexp) (string, error) {
	if port > 0 {
		host = strings.Join([]string{host, strconv.Itoa(port)}, ":")
	}

	var urlUser *url.Userinfo
	if len(user) > 0 {
		var userData = strings.Split(user, ":")
		if len(userData) == 2 {
			urlUser = url.UserPassword(userData[0], userData[1])
		} else {
			urlUser = url.User(user)
		}

	}

	var url = &url.URL{
		Scheme: scheme,
		User:   urlUser,
		Host:   host,
		Path:   "",
	}

	var hostname = url.String()

	var generatedPath, err = generatePath(path, parameters, requirements)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{hostname, generatedPath}, ""), nil
}

func generateQueryString(parameters map[string][]string) string {
	var values = url.Values(parameters)
	return values.Encode()
}
