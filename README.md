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

To create router that will handle "http://madmanlabs.com" on default port, just call NewMyRouter with parameters

```go
// func NewMyRouter(scheme string, host string, port int) *MyRouter
var router = NewMyRouter("http", "madmanlabs.com", 0)
```

In case you want port in Your URL, just set value above 0

```go
// func NewMyRouter(scheme string, host string, port int) *MyRouter
var router = NewMyRouter("http", "madmanlabs.com", 443)
```

License: **MIT**

