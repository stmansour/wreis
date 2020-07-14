package util

import (
	"fmt"
	"log"
	"runtime/debug"
)

// ErrCheck - saves a bunch of typing, prints error if it exists
//            and provides a traceback as well
func ErrCheck(err error) {
	if err != nil {
		fmt.Printf("error = %v\n", err)
		debug.PrintStack()
		log.Fatal(err)
	}
}
