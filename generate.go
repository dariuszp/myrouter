package myrouter

import (
	"errors"
	"strconv"
	"strings"
)

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
func GenerateRegExpFromPath(path string, requirements map[string]string) {
}
