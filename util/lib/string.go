package util

import "strings"

// Stripchars returns a string with the characters from chars removed
func Stripchars(str, chars string) string {
	return strings.Map(func(r rune) rune {
		if !strings.ContainsRune(chars, r) {
			return r
		}
		return -1
	}, str)
}
