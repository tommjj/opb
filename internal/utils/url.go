package utils

import "strings"

// BuildQuery replaces instance of s with arg in order,
// returning the replaced url and the number of arg used
func BuildQuery(url string, s string, arg ...string) (string, int) {
	if len(arg) == 0 {
		return url, 0
	}

	before, after, found := strings.Cut(url, s)
	if !found {
		return url, 0
	}

	last, num := BuildQuery(after, s, arg[1:]...)

	return before + arg[0] + last, 1 + num
}
