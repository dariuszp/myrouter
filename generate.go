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

//GenerateRegExpFromPath turns path to regexp pattern
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

// GeneratePath fill url pattern with parameters
func generatePath(path string, parameters map[string]string, requirements map[string]*regexp.Regexp) (string, error) {
	var extracted = extractParamNames(path)
	for _, parameterName := range extracted {
		var value, ok = parameters[parameterName]

		if !ok {
			return "", errors.New(strings.Join([]string{"Invalid path parameter", parameterName}, " "))
		}

		delete(parameters, parameterName)

		var rx, rxok = requirements[parameterName]
		if !rxok {
			rx = defaultReadParameterRegexp
		}
		if !rx.MatchString(value) {
			return "", errors.New(strings.Join([]string{"Parameter ", parameterName, "does not meet requirements", rx.String()}, " "))
		}

		parameterName = strings.Join([]string{"{", parameterName, "}"}, "")
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

// GenerateURL combine host, port and path to create absolute url
func generateURL(schema string, login string, password string, host string, port int, path string, parameters map[string]string, requirements map[string]*regexp.Regexp) (string, error) {
	var hostname string

	var basicAuth = ""
	if len(login) > 0 || len(password) > 0 {
		basicAuth = strings.Join([]string{login, ":", password, "@"}, "")
	}

	if port > 0 {
		hostname = strings.Join([]string{schema, "://", basicAuth, host, ":", strconv.Itoa(port)}, "")
	} else {
		hostname = strings.Join([]string{schema, "://", basicAuth, host}, "")
	}
	var generatedPath, err = generatePath(path, parameters, requirements)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{hostname, generatedPath}, ""), nil
}

func generateQueryString(parameters map[string]string) string {
	var queryStringParameters []string
	var queryParam string
	for name, value := range parameters {
		name = url.QueryEscape(name)
		value = url.QueryEscape(value)
		queryParam = strings.Join([]string{name, value}, "=")
		queryStringParameters = append(queryStringParameters, queryParam)
	}

	return strings.Join(queryStringParameters, "&")
}
