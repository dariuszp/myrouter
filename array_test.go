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

func TestStringCompareEqualArray(t *testing.T) {
	var a = []string{"a", "b"}
	var b = []string{"a", "b"}
	if !stringCompareArray(a, b) {
		t.Fail()
	}
}

func TestStringCompareEqualUnorderedArray(t *testing.T) {
	var a = []string{"a", "b"}
	var b = []string{"b", "a"}
	if !stringCompareArray(a, b) {
		t.Fail()
	}
}

func TestStringCompareNotEqualArray(t *testing.T) {
	var a = []string{"a", "b"}
	var b = []string{"a", "B"}
	if stringCompareArray(a, b) {
		t.Fail()
	}
}

func TestStringCompareNoCaseEqualArray(t *testing.T) {
	var a = []string{"a", "b"}
	var b = []string{"a", "B"}
	if !stringCompareNoCaseArray(a, b) {
		t.Fail()
	}
}
