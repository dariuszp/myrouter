package myrouter

import (
	"strconv"
	"strings"
)

// GeneratePath fill url pattern with parameters
func GeneratePath(path string, parameters map[string]string) string {
	var parameterName string
	for name, value := range parameters {
		parameterName = strings.Join([]string{"{", name, "}"}, "")
		path = strings.Replace(path, parameterName, value, -1)
	}

	return path
}

// GenerateURL combine host, port and path to create absolute url
func GenerateURL(schema string, host string, port int, path string, parameters map[string]string) string {
	var hostname string
	if port > 0 {
		hostname = strings.Join([]string{schema, "://", host, ":", strconv.Itoa(port)}, "")
	} else {
		hostname = strings.Join([]string{schema, "://", host}, "")
	}

	return strings.Join([]string{hostname, GeneratePath(path, parameters)}, "")
}
