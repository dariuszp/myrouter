package myrouter

import (
	"fmt"
	"testing"
)

func TestextractParamNamesEmptyPath(t *testing.T) {
	var parameters = extractParamNames("")
	if cap(parameters) > 0 {
		t.Fail()
	}
}

func TestextractParamNamesRootPath(t *testing.T) {
	var parameters = extractParamNames("/")
	if cap(parameters) > 0 {
		t.Fail()
	}
}

func TestextractParamNamesSingleArgument(t *testing.T) {
	var parameters = extractParamNames("/user/{id}")
	if cap(parameters) != 1 {
		t.Fail()
	}
	if parameters[0] != "id" {
		t.Fail()
	}
}

func TestextractParamNamesTwoArguments(t *testing.T) {
	var parameters = extractParamNames("/user/{id}/{slug}")
	if cap(parameters) != 2 {
		t.Fail()
	}
	if parameters[0] != "id" {
		t.Fail()
	}
	if parameters[1] != "slug" {
		t.Fail()
	}
}

func TestExtractParamsFromRoute(t *testing.T) {
	var route *Route
	var err error
	var parameters map[string]string
	route, err = NewRoute("test", []string{"GET"}, "http", "example.com", 0, "/api/group-{group}/user-{id}", map[string]string{})
	parameters, err = extractParamsFromRoute(route, "/api/group-global/user-5")

	if err != nil {
		fmt.Println("Error while extracting parameters")
		t.Fail()
	}

	if len(parameters) != 2 {
		fmt.Println("Parameters count does not match")
		t.Fail()
	}

	if parameters["group"] != "global" {
		fmt.Println("Invalid global variable")
		t.Fail()
	}

	if parameters["id"] != "5" {
		fmt.Println("Invalid user variable")
		t.Fail()
	}
}

func TestExtractParamsFromRouteWithDifferentRoute(t *testing.T) {
	var route *Route
	var err error
	var parameters map[string]string
	route, err = NewRoute("test", []string{"GET"}, "http", "example.com", 0, "/api/group-{group}/user-{id}", map[string]string{})
	parameters, err = extractParamsFromRoute(route, "/api/user-global/data-5")

	if err == nil {
		fmt.Println("Error expected")
		t.Fail()
	}

	if len(parameters) != 2 {
		fmt.Println("Parameters count should still match")
		t.Fail()
	}

	if parameters["group"] != "" {
		fmt.Println("Invalid global variable")
		t.Fail()
	}

	if parameters["id"] != "" {
		fmt.Println("Invalid user variable")
		t.Fail()
	}
}

func TestExtractParamsFromRouteWithNoParams(t *testing.T) {
	var route *Route
	var err error
	var parameters map[string]string
	route, err = NewRoute("test", []string{"GET"}, "http", "example.com", 0, "/api/group", map[string]string{})
	parameters, err = extractParamsFromRoute(route, "/api/group")

	if err != nil {
		fmt.Println("Error while extracting parameters")
		t.Fail()
	}

	if len(parameters) != 0 {
		fmt.Println("Parameters count does not match, it should be 0")
		t.Fail()
	}

	var _, ok1 = parameters["group"]
	if ok1 {
		fmt.Println("Invalid global variable")
		t.Fail()
	}

	var _, ok2 = parameters["id"]
	if ok2 {
		fmt.Println("Invalid user variable")
		t.Fail()
	}
}
