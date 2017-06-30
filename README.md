# MyRouter

package dariuszp/myrouter

My router for GO http server. I was pissed off by existing routers and way they handle paths and params.
**This router is still in development.**

I do not bundle this router with any manager for handlers, controllers or whatever you come up with. Each router is named. You are free to decorate this router with methods that will handle your kind of handlers.

## Concept

* Each route is named
* Rou need to provide scheme, host and port
* if you provide port above 0, it will be added to url. Leave 0 if you don't want port in url
* Route params are check against regexp, default regexp for route param is "[^/]+" so "/" in route param is not allowed by default
* when you match path, you get route data and params passed to route
* if more than one route match path, first will be working
* MatchByMethod is little quicker than Match because it will loop only over routes with that method

## Usage

This part will explain how to use "My Router"

### Creating router

> func NewMyRouter(scheme string, host string, port int) *MyRouter

To create router that will handle "http://madmanlabs.com" on default port, just call NewMyRouter with parameters

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
```

In case you want port in Your URL, just set value above 0

```go
var router = NewMyRouter("http", "madmanlabs.com", 443)
```

### Adding new route

> func (router *MyRouter) Add(name string, methods []string, path string, requirements map[string]string) (bool, error)

Too add new route, without requirements:

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
router.Add("user_profile", make([]string), "/user/{id}", make(map[string]string))
```

but often we have specific ID format, for example only numbers and we want to force this format on router:

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
var status, err = router.Add("user_profile", make([]string), "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
```

Adding router will result in bool value (check if route was added) and potential error.

### Adding custom route

> func (router *MyRouter) AddCustom(name string, methods []string, scheme string, unsecureUser string, host string, port int, path string, requirements map[string]string) (bool, error)

There is also an edge case when we want router with different host, port etc. In that case we need to add custom route:

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
var status, err = router.AddCustom("user_profile", make([]string), "https", "mylogin:mypassword", "secure.example.com", 440, "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
```

### Setting route method 

If we want route to work only with specific methods, we can set them while creating the route

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
var status, err = router.Add("user_profile", []string{"get", "post"}, "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
```

### Removing route

> func (router *MyRouter) Remove(name string) bool

Removing route is simple. Since all routes need to have internal name, you can just simply call

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
// Adding route
var status, err = router.Add("user_profile", []string{"get", "post"}, "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
// Removing aded route
router.Remove("user_profile")
```

### Getting the route

We can also retrive added routes

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
// Adding route
var status, err = router.Add("user_profile", []string{"get", "post"}, "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
// Get aded route
var route = router.Get("user_profile")
```

### Matching HTTP path against router

Ok, we have our router, we added bunch of routes, how do we know what route was used by the user?




License: **MIT**

