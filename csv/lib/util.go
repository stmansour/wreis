package wcsv

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	util "wreis/util/lib"
)

func deepCopy(l1 []ColumnDef) []ColumnDef {
	var l2 []ColumnDef
	var count = len(l1)
	for i := 0; i < count; i++ {
		var x ColumnDef
		x.Name = l1[i].Name
		x.Required = l1[i].Required
		x.CaseSensitive = l1[i].CaseSensitive
		x.CanonicalIndex = l1[i].CanonicalIndex
		x.Index = l1[i].Index
		l2 = append(l2, x)
	}

	return l2
}

// ParseInt is a utility function that optimizes parsing of integer values
// for the csv subsystem. If there is an error in parsing it is appended to the
// to the supplied error list along with the line number on which the error
// occurred.
//------------------------------------------------------------------------------
func ParseInt(s string, line int, errlist []error) (int, []error) {
	var i int64
	i, errlist = (ParseInt64(s, line, errlist))
	return int(i), errlist
}

// ParseInt64 is a utility function that optimizes parsing of integer values
// for the csv subsystem. If there is an error in parsing it is appended to the
// to the supplied error list along with the line number on which the error
// occurred.
//------------------------------------------------------------------------------
func ParseInt64(s string, line int, errlist []error) (int64, []error) {
	if len(s) == 0 {
		return int64(0), errlist
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		errlist = append(errlist, fmt.Errorf("Line %d: %s", line, err.Error()))
	}

	return i, errlist
}

// ParseFloat64 is a utility function that optimizes parsing of float values
// for the csv subsystem. If there is an error in parsing it is appended to the
// to the supplied error list along with the line number on which the error
// occurred.
//------------------------------------------------------------------------------
func ParseFloat64(s string, line int, errlist []error) (float64, []error) {
	if len(s) == 0 {
		return float64(0), errlist
	}
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		errlist = append(errlist, fmt.Errorf("Line %d: %s", line, err.Error()))
	}

	return i, errlist
}

// ParseDate is a utility function that optimizes parsing of date values
// for the csv subsystem. If there is an error in parsing it is appended to the
// to the supplied error list along with the line number on which the error
// occurred.
//------------------------------------------------------------------------------
func ParseDate(s string, line int, errlist []error) (time.Time, []error) {
	if len(s) == 0 {
		return util.TIME0, errlist
	}

	i, err := util.StringToDate(s)
	if err != nil {
		errlist = append(errlist, fmt.Errorf("Line %d: %s", line, err.Error()))
	}

	return i, errlist
}

// YesNoToInt takes multiple forms of "Yes" and converts to integer 1, multiple
// forms of "No" to integer 0
//
// INPUTS
// si - input string
//
// RETURNS
// 1 if the string is Yes, True, or 1,  0 otherwise
// any error encountered
//------------------------------------------------------------------------------
func YesNoToInt(si string) (int64, error) {
	s := strings.ToUpper(strings.TrimSpace(si))
	switch {
	case s == "Y" || s == "YES" || s == "1" || s == "T" || s == "TRUE":
		return 1, nil
	case s == "N" || s == "NO" || s == "0" || s == "F" || s == "FALSE":
		return 0, nil
	default:
		err := fmt.Errorf("Unrecognized yes/no string: %s", si)
		return 0, err
	}
}

// GetBitFlagValue will return the supplied value v if the string is any
// recognized form of "yes" or true.  It will return 0 if the string is any
// recognized form of "no" or false.  If the string is not a recognized form
// of yes or no, it returns an error.
//
// INPUTS
// si      - input string
// errlist
//
// RETURNS
// v   - if the string is Yes, True, or 1
// 0   - if the string is No, False, or 0, or if an error is encountered
//------------------------------------------------------------------------------
func GetBitFlagValue(s string, v uint64, errlist []error) (uint64, []error) {
	x, err := YesNoToInt(s)
	if err != nil {
		errlist = append(errlist, err)
		return 0, errlist
	}
	if x > 0 {
		return v, errlist
	}
	return 0, errlist
}
