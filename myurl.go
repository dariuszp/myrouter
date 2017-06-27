package myrouter

// MyURL is a result from parsing url/path
type MyURL struct {
	Scheme     string
	Login      string
	Password   string
	Host       string
	Port       string
	Path       string
	Query      map[string]string
	RawQuery   string
	Fragment   string
	Parameters map[string]string
}
