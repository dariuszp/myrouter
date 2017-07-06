package myrouter

import (
	"errors"
	"regexp"
)

const (
	parameterSubpattern = `(\{([a-z]+[a-z0-9_]*)\})`
)

var parameterRegexp = regexp.MustCompile(parameterSubpattern)

// Extract parameters from path
func extractParamNames(path string) []string {
	var parameters []string

	var extracted = parameterRegexp.FindAllStringSubmatch(path, -1)
	for _, element := range extracted {
		parameters = append(parameters, element[2])
	}

	return parameters
}

func extractParamsFromRoute(route *Route, path string) (PathParameters, error) {
	var result = make(map[string]string)
	var name, value string
	var err error

	var parameters = route.MatchRegexp.FindAllStringSubmatch(path, -1)
	var countMatch = true
	if len(parameters) == 0 || (len(parameters[0]) != (len(route.Parameters) + 1)) {
		countMatch = false
		err = errors.New("Parameters count does not match")
	}

	for i := 0; i < len(route.Parameters); i++ {
		name = route.Parameters[i]
		if countMatch {
			value = parameters[0][i+1]
		} else {
			value = ""
		}
		result[name] = value
	}

	return result, err
}
