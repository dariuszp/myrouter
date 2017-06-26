package myrouter

import (
	"errors"
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
	return path, nil
}

// GenerateURL combine host, port and path to create absolute url
func generateURL(schema string, host string, port int, path string, parameters map[string]string, requirements map[string]*regexp.Regexp) (string, error) {
	var hostname string
	if port > 0 {
		hostname = strings.Join([]string{schema, "://", host, ":", strconv.Itoa(port)}, "")
	} else {
		hostname = strings.Join([]string{schema, "://", host}, "")
	}
	var generatedPath, err = generatePath(path, parameters, requirements)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{hostname, generatedPath}, ""), nil
}
