package main

import (
	"fmt"
	"os"

	util "wreis/util/lib"
)

type udate struct {
	s string
	e int
}

func main() {
	var m = []udate{
		{"2018-06-20 00:00:00 UTC", 20},
		{"2020-10-01 10:37:45 UTC", 1},
	}

	for i := 0; i < len(m); i++ {
		dt, err := util.StringToDate(m[i].s)
		if err != nil {
			fmt.Printf("String = %s, expect no error.  got error: %s\n", m[i].s, err.Error())
			os.Exit(1)
		}
		if dt.Day() != m[i].e {
			fmt.Printf("Day number is wrong, expected %d, got %d\n", m[i].e, dt.Day())
			os.Exit(1)
		}
	}
	fmt.Printf("Success!  No errors found\n")
}
