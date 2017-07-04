package myrouter

// URLParameters containts url params and query string params
type URLParameters map[string][]string

// PathParameters contains params extracted from path
type PathParameters map[string]string

// Requirements contains regexps for route params
// This do not apply for query strings
type Requirements map[string]string

// Routes contain list of routes
type Routes map[string]*Route

// Verbs contains map of verbs with list of routes
type Verbs map[string]map[string]*Route
