package myrouter

// NewMyRouter create instance of a router
func NewMyRouter() *MyRouter {
	var router = &MyRouter{make(map[string]map[string]*Route), make(map[string]*Route)}
	return router
}
