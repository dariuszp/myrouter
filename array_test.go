package myrouter

import (
	"testing"
)

func TestStringInArray(t *testing.T) {
	if !stringInArray([]string{"one", "two"}, "two") {
		t.Fail()
	}
}

func TestStringNotInArray(t *testing.T) {
	if stringInArray([]string{"one", "two"}, "three") {
		t.Fail()
	}
}
