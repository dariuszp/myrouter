# MyRouter

[![Build Status](https://travis-ci.org/dariuszp/myrouter.png?branch=master)](https://travis-ci.org/dariuszp/myrouter) Master
[![Build Status](https://travis-ci.org/dariuszp/myrouter.png?branch=develop)](https://travis-ci.org/dariuszp/myrouter) Develop

## WARNING

> This package is in heavy development, it's very unstable

## Requirements

Go 1.8

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

## Router

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

## Route

### Adding new route

> func (router *MyRouter) Add(name string, methods []string, path string, requirements map[string]string) (bool, error)

Too add new route, without requirements:

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
router.Add("user_profile", []string{}, "/user/{id}", map[string]string{})
```

but often we have specific ID format, for example only numbers and we want to force this format on router:

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
var route, err = router.Add("user_profile", []string{}, "/user/{id}", map[string]string{"id": "[1-9]+[0-9]*"})
```

Adding router will result in bool value (check if route was added) and potential error.

### Adding custom route

> func (router *MyRouter) AddCustom(name string, methods []string, scheme string, unsecureUser string, host string, port int, path string, requirements map[string]string) (bool, error)

There is also an edge case when we want router with different host, port etc. In that case we need to add custom route:

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
var route, err = router.AddCustom("user_profile", []string{}, "https", "mylogin:mypassword", "secure.example.com", 440, "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
```

### Setting route method 

If we want route to work only with specific methods, we can set them while creating the route

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
var route, err = router.Add("user_profile", []string{"get", "post"}, "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
```

### Removing route

> func (router *MyRouter) Remove(name string) bool

Removing route is simple. Since all routes need to have internal name, you can just simply call

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
// Adding route
var route, err = router.Add("user_profile", []string{"get", "post"}, "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
// Removing aded route
router.Remove("user_profile")
```

### Getting the route

We can also retrive added routes

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
// Adding route
var route, err = router.Add("user_profile", []string{"get", "post"}, "/user/{id}", map[string]string{ "id": "[1-9]+[0-9]*" })
// Get aded route
var route = router.Get("user_profile")
```

### Matching HTTP path against router

Ok, we have our router, we added bunch of routes, how do we know what route was used by the user?
Let say that our user is calling us using url:

> GET http:/madmanlabs.com/user/poltorak-dariusz?tab=contacts

We need our router to be ready for him:

```go
var router = NewMyRouter("http", "madmanlabs.com", 0)
var _, err = router.Add("user_profile", []string{}, "/user/{slug}", map[string]string{ "slug": "[a-z]+[a-z\\-]*" })
```

then we can try to match url he used with our service:

```go
    var result, err = router.Match("GET", "http:/madmanlabs.com/user/poltorak-dariusz?tab=contacts")
```

This will create instance of MyURL (types are at the bottom of readme) and return optional error. If error is nil, you will have access to: 

| Scheme     | string              | For example http, https etc.                                                                                                                                           |
|------------|---------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| User       | string              | Username, if we use basic auth and provide "dariusz:poltorak", user will be "dariusz"                                                                                  |
| Password   | string              | If we use basic auth and provide "dariusz:poltorak", password will be "poltorak"                                                                                       |
| Host       | string              | For example "madmanlabs.com"                                                                                                                                           |
| Port       | int                 | Default is 0. When port is provided, you will get exact number                                                                                                         |
| Path       | string              | For example "/user/poltorak-dariusz"                                                                                                                                   |
| Parameters | map[string]string   | URL parameters, when your pattern is "/user/{slug}" and You provide path "/user/dariusz-poltorak" you will get result: map[string]string{ "slug": "dariusz-poltorak" } |
| Query      | map[string][]string | Same as route parameters but query contains arrays in values. Mostly because GET method allow you to provide arrays in query string                                    |
| Fragment   | string              | Everything after hash ("#")                                                                                                                                            |
| Route      | *Route              | Instance of added route                                                                                                                                                |
#### No match

**In case you have no match, you will still get instance of MyURL but second parameter will contain error instead of nil.**

#### Other type of matches

Sometimes you just want to match path, in that is the case, just use

```go
    var result, err = router.MatchPath("/user/poltorak-dariusz?tab=contacts")
```

or just orl while ignoring the method

```go
    var result, err = router.MatchURL("http:/madmanlabs.com/user/poltorak-dariusz?tab=contacts")
```

Match is alias to MatchURLByMethod, you can call it directly:

```go
    var result, err = router.MatchURLByMethod("GET", "http:/madmanlabs.com/user/poltorak-dariusz?tab=contacts")
```

You can do the same with the patch

```go
    var result, err = router.MatchPathByMethod("GET", "/user/poltorak-dariusz?tab=contacts")
```

**There is also simple, default Match(method, url) method. It's an alias to MatchURLByMethod(method, url)**


### Generating Path / URL

The reason we name our routes is that we also want to generate addresses based on that name, for example:

```go
    var router = NewMyRouter("http", "example.com", 3000)

    router.Add("dashboard", []string{}, "/dashboard", make(map[string]string))
    router.Add("profile", []string{"GET"}, "/user/{slug}", map[string]string{"slug": "[a-z\\-]+"})
    router.Add("message", []string{"POST", "PUT"}, "/message/{channel}/{type}", map[string]string{"type": "error|success"})
```

Now we want to generate path to message:

```go
    var url, err = router.URL("message", map[string][]string{"channel": []string{"sms"}, "type": []string{"error"}})
```

Now err is nil and url equals:

> http://example.com:3000/message/sms/error

Remember that each route param can be mentioned in requirements. Default requirement for route value is NOT to contain slash ("/", that means "[^/]"). Everything else will pass. So doing something like this:

```go
    var url, err = router.URL("message", map[string][]string{"channel": []string{"sms"}, "type": []string{"warning"}})
```

Will result in error because "type" accept only "error or success" with regexp "error|success".


You can also notice the fact that you provide parameters in different format. You use map[string][]string instead of map[string]string. That's because parameters you use must be compatibile with query string format. Because unused parameters go to query string, like this:

```go
    var url, err = router.URL("message", map[string][]string{"channel": []string{"sms"}, "type": []string{"error"}, "ids": []string{"5", "6"}})
```

Now err is nil and url equals:

> http://example.com:3000/message/sms/error?ids=5&ids=6

Usually when You just want internal url in your system, you only need path, not entire url. For that you have:

```go
    var path, err = router.Path("message", map[string][]string{"channel": []string{"sms"}, "type": []string{"error"}, "ids": []string{"5", "6"}})
```

Now err is nil and url equals:

> /message/sms/error?ids=5&ids=6

Each generator return path/url and error. That is because there are few instances where you can get empty url/path and error. For example if You have route that 
request parameter "slug" containing letters only, if there is any other character, url will be empty and error will be returned.

## Types

#### MyURL 

| Scheme     | string              | For example http, https etc.                                                                                                                                           |
|------------|---------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| User       | string              | Username, if we use basic auth and provide "dariusz:poltorak", user will be "dariusz"                                                                                  |
| Password   | string              | If we use basic auth and provide "dariusz:poltorak", password will be "poltorak"                                                                                       |
| Host       | string              | For example "madmanlabs.com"                                                                                                                                           |
| Port       | int                 | Default is 0. When port is provided, you will get exact number                                                                                                         |
| Path       | string              | For example "/user/poltorak-dariusz"                                                                                                                                   |
| Parameters | map[string]string   | URL parameters, when your pattern is "/user/{slug}" and You provide path "/user/dariusz-poltorak" you will get result: map[string]string{ "slug": "dariusz-poltorak" } |
| Query      | map[string][]string | Same as route parameters but query contains arrays in values. Mostly because GET method allow you to provide arrays in query string                                    |
| Fragment   | string              | Everything after hash ("#")                                                                                                                                            |
| Route      | *Route              | Instance of added route                                                                                                                                                |

#### Route 

| Name                                                                                    | string                    | Name of the route                                                                                                                                                          |
|-----------------------------------------------------------------------------------------|---------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Methods                                                                                 | string                    | List of supported methods (empty to support all)                                                                                                                           |
| Scheme                                                                                  | string                    | Scheme like "https"                                                                                                                                                        |
| UnsecureUser                                                                            | string                    | Unsecure in name suggests it's a very bad idea to use it                                                                                                                   |
| Host                                                                                    | string                    | Host, like "example.com"                                                                                                                                                   |
| Port                                                                                    | int                       | Port, put 0 to omit                                                                                                                                                        |
| Path                                                                                    | string                    | Like "/test"                                                                                                                                                               |
| Parameters                                                                              | []string                  | List of url parameters, in case of path "/user/{id}" you will have "[]string{"id"}" there                                                                                  |
| MatchRegexp                                                                             | *regexp.Regexp            | Path turned to regexp                                                                                                                                                      |
| Requirements                                                                            | map[string]*regexp.Regexp | Requirements for parameters. For example using ID you want to force only numbers. In that case you can type: map[string]string{ "id": regexp.MustCompile("[1-9]+[0-9]*") } |
| SetMethods(methods []string)                                                            | (bool, error)             | Replace method list                                                                                                                                                        |
| Match(path string)                                                                      | bool                      | Match route with provided path                                                                                                                                             |
| MatchURL(urlAddress string)                                                             | bool                      | Match route with provided URL                                                                                                                                              |
| MatchMethod(method string, path string)                                                 | bool                      | Match route with provided path and method                                                                                                                                  |
| MatchURLMethod(method string, urlAddress string)                                        | bool                      | Match route against method and url                                                                                                                                         |
| ParsePath(path string)                                                                  | (*MyURL, error)           | Parse given path and return instance of MyURL                                                                                                                              |
| ParseURL(urlAddress string)                                                             | (*MyURL, error)           | Parse given url and return instance of MyURL                                                                                                                               |
| GeneratePath(parameters URLParameters)                                                  | (string, error)           | Generate path from route using parameters                                                                                                                                  |
| GeneratePathWithFragment(parameters URLParameters, fragment string)                     | (string, error)           | Same as GeneratePath but allow you to provide hash string (#something) at the end                                                                                          |
| GenerateURL(parameters URLParameters)                                                   | (string, error)           | Generate full url from route                                                                                                                                               |
| GenerateURLWithFragment(parameters URLParameters, fragment string)                      | (string, error)           | Same as GenerateURL but with hash string                                                                                                                                   |
| GenerateUnsecureURL(user string, parameters URLParameters)                              | (string, error)           | Same as generate URL but allow you to pass login and password. Very bad idea but it's there.                                                                               |
| GenerateUnsecureURLWithFragment(user string, parameters URLParameters, fragment string) | (string, error)           | Same as GenerateUnsecureURL but with hash string. Bad idea to use it.                                                                                                      |
| Generate(parameters URLParameters)                                                      | (string, error)           | Alias to GeneratePath                                                                                                                                                      |
| GenerateWithFragment(parameters URLParameters, fragment string)                         | (string, error)           | Alias to GeneratePathWithFragment                                                                                                                                          |



License: **MIT**

