package myrouter

import (
	"testing"
)

func TestarrayContainsString(t *testing.T) {
	if !arrayContainsString([]string{"one", "two"}, "two") {
		t.Fail()
	}
}

func TestStringNotInArray(t *testing.T) {
	if arrayContainsString([]string{"one", "two"}, "three") {
		t.Fail()
	}
}

func TestStringCompareEqualArray(t *testing.T) {
	var a = []string{"a", "b"}
	var b = []string{"a", "b"}
	if !arrayCompareString(a, b) {
		t.Fail()
	}
}

func TestStringCompareEqualUnorderedArray(t *testing.T) {
	var a = []string{"a", "b"}
	var b = []string{"b", "a"}
	if !arrayCompareString(a, b) {
		t.Fail()
	}
}

func TestStringCompareNotEqualArray(t *testing.T) {
	var a = []string{"a", "b"}
	var b = []string{"a", "B"}
	if arrayCompareString(a, b) {
		t.Fail()
	}
}

func TestStringCompareNoCaseEqualArray(t *testing.T) {
	var a = []string{"a", "b"}
	var b = []string{"a", "B"}
	if !arrayCompareStringNoCase(a, b) {
		t.Fail()
	}
}
