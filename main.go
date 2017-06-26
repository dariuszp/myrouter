package myrouter

// NewMyRouter create instance of a router
func NewMyRouter(schema string, host string, port int) *MyRouter {
	var router = &MyRouter{schema, host, port, make(map[string]map[string]*Route), make(map[string]*Route)}
	return router
}
