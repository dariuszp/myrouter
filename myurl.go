package myrouter

// MyURL contains parsed url
type MyURL struct {
	Scheme     string
	User       string
	Password   string
	Host       string
	Port       int
	Path       string
	Parameters PathParameters
	Query      URLParameters
	Fragment   string
	Route      *Route
}

// NewMyURL create instance of MyURL
func NewMyURL(scheme string, user string, password string, host string, port int, path string, parameters PathParameters, query URLParameters, fragment string, route *Route) *MyURL {
	return &MyURL{scheme, user, password, host, port, path, parameters, query, fragment, route}
}

// NewEmptyMyURL creates empty MyURL instance
func NewEmptyMyURL() *MyURL {
	return &MyURL{"", "", "", "", 0, "", make(PathParameters), make(URLParameters), "", nil}
}
