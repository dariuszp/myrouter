package myrouter

import (
	"strings"
)

// Check if string is in given array/slice
func stringInArray(list []string, item string) bool {
	for _, element := range list {
		if element == item {
			return true
		}
	}
	return false
}

// Check if lowercase string is in given array/slice
func stringInArrayNoCase(list []string, item string) bool {
	item = strings.ToLower(item)
	for _, element := range list {
		if strings.ToLower(element) == item {
			return true
		}
	}
	return false
}

// stringCompareArray compare array of strings
func stringCompareArray(listA []string, listB []string) bool {
	for _, a := range listA {
		if !stringInArray(listB, a) {
			return false
		}
	}
	return true
}

// stringCompareNoCaseArray compare array of lowercase strings
func stringCompareNoCaseArray(listA []string, listB []string) bool {
	for _, a := range listA {
		if !stringInArrayNoCase(listB, a) {
			return false
		}
	}
	return true
}
