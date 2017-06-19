package myrouter

// Route represent single http route
type Route struct {
	name    string
	methods []string
	schema  string
	host    string
	port    int
	path    string
	handler func()
}
