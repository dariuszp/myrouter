package myrouter

import (
	"reflect"
	"strings"
)

// Check if string is in given array/slice
func arrayContainsString(list []string, item string) bool {
	for _, element := range list {
		if element == item {
			return true
		}
	}
	return false
}

// Check if lowercase string is in given array/slice
func arrayContainsStringNoCase(list []string, item string) bool {
	item = strings.ToLower(item)
	for _, element := range list {
		if strings.ToLower(element) == item {
			return true
		}
	}
	return false
}

// stringCompareArray compare array of strings
func arrayCompareString(listA []string, listB []string) bool {
	for _, element := range listA {
		if !arrayContainsString(listB, element) {
			return false
		}
	}
	return true
}

// stringCompareNoCaseArray compare array of lowercase strings
func arrayCompareStringNoCase(listA []string, listB []string) bool {
	for _, element := range listA {
		if !arrayContainsStringNoCase(listB, element) {
			return false
		}
	}
	return true
}

func isArray(v interface{}) bool {
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice
}
