package compiler_util

// returns true if a given string 'a' is in a given array 'list'
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// returns true if a given string 'a' is in a given map 'list'
func StringInMap(a string, list map[string]Node) bool {
	_, ok := list[a]
	if ok {
		return true
	} else {
		return false
	}
}
