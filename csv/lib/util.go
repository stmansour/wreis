package wcsv

import (
	"fmt"
	"strconv"
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
	i, err := util.StringToDate(s)
	if err != nil {
		errlist = append(errlist, fmt.Errorf("Line %d: %s", line, err.Error()))
	}

	return i, errlist
}
