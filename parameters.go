package myrouter

import (
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
