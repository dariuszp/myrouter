package myrouter

import (
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
