package myrouter

import "regexp"

// URLParameters containts url params and query string params
type URLParameters map[string][]string

// PathParameters contains params extracted from path
type PathParameters map[string]string

// Requirements contains regexps for route params
// This do not apply for query strings
type Requirements map[string]string

// CompiledRequirements contains Requirements turn to regexp
type CompiledRequirements map[string]*regexp.Regexp

// Routes contain list of routes
type Routes map[string]*Route

// Verbs contains map of verbs with list of routes
type Verbs map[string]map[string]*Route
