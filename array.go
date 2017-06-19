package myrouter

// Check if string is in given array/slice
func stringInArray(list []string, item string) bool {
	for _, element := range list {
		if element == item {
			return true
		}
	}
	return false
}
