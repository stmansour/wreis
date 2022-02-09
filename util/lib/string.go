package util

import (
	"fmt"
	"strconv"
	"strings"
)

// Stripchars returns a string with the characters from chars removed
func Stripchars(str, chars string) string {
	return strings.Map(func(r rune) rune {
		if !strings.ContainsRune(chars, r) {
			return r
		}
		return -1
	}, str)
}

// Tline returns a string of dashes that is the specified length
func Tline(n int) string {
	p := make([]byte, n)
	for i := 0; i < n; i++ {
		p[i] = '-'
	}
	return string(p)
}

// Mkstr returns a string of n of the supplied character that is the specified length
func Mkstr(n int, c byte) string {
	p := make([]byte, n)
	for i := 0; i < n; i++ {
		p[i] = c
	}
	return string(p)
}

// SafeURLFilename is a function that makes filenames safe for urls.
// Filenames that will be used in a URL could contain characters that will
// cause problems.  This function replaces those characters with dashes.
//
// INPUTS
//  f = filename to change
//
// Mapping:
//  ' ' -> '-'
//  ',' -> '-'
//------------------------------------------------------------------------------
func SafeURLFilename(f string) string {
	p := []byte(f)
	n := len(p)
	for i := 0; i < n; i++ {
		if p[i] == ' ' || p[i] == ',' {
			p[i] = '-'
		}
	}
	return string(p)
}

// IntFromString converts the supplied string to an int64 value. If there
// is a problem in the conversion, it generates an error message. To suppress
// the error message, pass in "" for errmsg.
func IntFromString(sa string, errmsg string) (int64, error) {
	var n = int64(0)
	s := strings.TrimSpace(sa)
	if len(s) > 0 {
		i, err := strconv.Atoi(s)
		if err != nil {
			if errmsg != "" {
				return 0, fmt.Errorf("IntFromString: %s: %s", errmsg, s)
			}
			return n, err
		}
		n = int64(i)
	}
	return n, nil
}
