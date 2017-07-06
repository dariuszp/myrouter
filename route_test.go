package myrouter

import (
	"fmt"
	"strings"
	"testing"
)

func TestRouteCreation(t *testing.T) {
	var route *Route
	var err error
	var methods = []string{"get", "post"}
	route, err = NewRoute("test", methods, "https", "example.com", 0, "/api/user/{id}", map[string]string{"id": "[1-9]+[0-9]*"})

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if route.Name != "test" {
		fmt.Println("Invalid route name")
		t.Fail()
	}

	if !arrayCompareString(route.Methods, methods) {
		fmt.Println("Route methods does not match")
		t.Fail()
	}

	if route.Scheme != "https" || route.Host != "example.com" || route.Port != 0 {
		fmt.Println("Invalid route base url")
		t.Fail()
	}

	if route.Path != "/api/user/{id}" {
		fmt.Println("Invalid path")
		t.Fail()
	}

	if route.MatchRegexp.String() != "/api/user/([1-9]+[0-9]*)" {
		fmt.Println("Invalid regexp")
		t.Fail()
	}

	if len(route.Parameters) != 1 {
		fmt.Println("Invalid parameters count")
		t.Fail()
	}

	if route.Parameters[0] != "id" {
		fmt.Println("Invalid parameters")
		t.Fail()
	}
}

func TestRouteCreationWithNoParams(t *testing.T) {
	var route *Route
	var err error
	var methods = []string{"get", "post"}
	route, err = NewRoute("test", methods, "https", "example.com", 0, "/api/user", map[string]string{"id": "[1-9]+[0-9]*"})

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if route.Name != "test" {
		fmt.Println("Invalid route name")
		t.Fail()
	}

	if !arrayCompareString(route.Methods, methods) {
		fmt.Println("Route methods does not match")
		t.Fail()
	}

	if route.Scheme != "https" || route.Host != "example.com" || route.Port != 0 {
		fmt.Println("Invalid route base url")
		t.Fail()
	}

	if route.Path != "/api/user" {
		fmt.Println("Invalid path")
		t.Fail()
	}

	if route.MatchRegexp.String() != "/api/user" {
		fmt.Println("Invalid regexp")
		t.Fail()
	}

	if len(route.Parameters) != 0 {
		fmt.Println("Invalid parameters count")
		t.Fail()
	}
}

func TestURL(t *testing.T) {
	var route *Route
	var err error
	var url string
	var methods = []string{"get", "post"}
	route, err = NewRoute("test", methods, "https", "example.com", 0, "/api/user/{id}", map[string]string{"id": "[1-9]+[0-9]*"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	url, err = route.GenerateURL(map[string][]string{"id": []string{"5"}})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if url != "https://example.com/api/user/5" {
		fmt.Println(strings.Join([]string{"Invalid url", url}, " "))
		t.Fail()
	}
}

func TestURLWithNoParams(t *testing.T) {
	var route *Route
	var err error
	var url string
	var methods = []string{"get", "post"}
	route, err = NewRoute("test", methods, "https", "example.com", 0, "/api/user", map[string]string{"id": "[1-9]+[0-9]*"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	url, err = route.GenerateURL(map[string][]string{"id": []string{"5"}})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if url != "https://example.com/api/user?id=5" {
		fmt.Println(strings.Join([]string{"Invalid url", url}, " "))
		t.Fail()
	}
}

func TestURLWithUser(t *testing.T) {
	var route *Route
	var err error
	var url string
	var methods = []string{"get", "post"}
	route, err = NewRoute("test", methods, "https", "example.com", 0, "/api/user", map[string]string{"id": "[1-9]+[0-9]*"})
	route.UnsecureUser = "darek:poltorak"

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	url, err = route.GenerateURL(map[string][]string{"id": []string{"5"}})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if url != "https://darek:poltorak@example.com/api/user?id=5" {
		fmt.Println(strings.Join([]string{"Invalid url", url}, " "))
		t.Fail()
	}
}
