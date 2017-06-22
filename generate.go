package myrouter

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var defaultReadParameterRegexp = "[^/]+"

// GeneratePath fill url pattern with parameters
func generatePath(path string, parameters map[string]string) (string, error) {
	var extracted = extractParamNames(path)
	for _, parameterName := range extracted {
		var value, ok = parameters[parameterName]
		if !ok {
			return "", errors.New(strings.Join([]string{"Invalid path parameter", parameterName}, " "))
		}
		parameterName = strings.Join([]string{"{", parameterName, "}"}, "")
		path = strings.Replace(path, parameterName, value, -1)
	}
	return path, nil
}

// GenerateURL combine host, port and path to create absolute url
func generateURL(schema string, host string, port int, path string, parameters map[string]string) (string, error) {
	var hostname string
	if port > 0 {
		hostname = strings.Join([]string{schema, "://", host, ":", strconv.Itoa(port)}, "")
	} else {
		hostname = strings.Join([]string{schema, "://", host}, "")
	}
	var generatedPath, err = generatePath(path, parameters)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{hostname, generatedPath}, ""), nil
}

//GenerateRegExpFromPath turns path to regexp pattern
func GenerateRegExpFromPath(path string, requirements map[string]string) *regexp.Regexp {
	var parameterEscapedName, pattern, escapedPath, result, patternReplace string
	var ok bool
	var parameters = extractParamNames(path)

	escapedPath = regexp.QuoteMeta(path)
	if len(parameters) == 0 {
		return regexp.MustCompile(escapedPath)
	}

	result = escapedPath
	for _, parameterName := range parameters {
		pattern, ok = requirements[parameterName]
		if ok && len(pattern) > 0 {
			patternReplace = pattern
		} else {
			patternReplace = defaultReadParameterRegexp
		}
		parameterEscapedName = strings.Join([]string{"\\{", parameterName, "\\}"}, "")
		patternReplace = strings.Join([]string{"(", patternReplace, ")"}, "")
		result = strings.Replace(result, parameterEscapedName, patternReplace, -1)
	}

	return regexp.MustCompile(result)
}
