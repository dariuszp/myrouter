package myrouter

import (
	"testing"
)

func TestExtractEmptyPath(t *testing.T) {
	var parameters = Extract("")
	if cap(parameters) > 0 {
		t.Fail()
	}
}

func TestExtractRootPath(t *testing.T) {
	var parameters = Extract("/")
	if cap(parameters) > 0 {
		t.Fail()
	}
}

func TestExtractSingleArgument(t *testing.T) {
	var parameters = Extract("/user/{id}")
	if cap(parameters) != 1 {
		t.Fail()
	}
	if parameters[0] != "id" {
		t.Fail()
	}
}

func TestExtractTwoArguments(t *testing.T) {
	var parameters = Extract("/user/{id}/{slug}")
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
