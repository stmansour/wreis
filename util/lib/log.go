package util

import (
	"fmt"
	"log"
)

// Ulog is Phonebooks's standard logger
func Ulog(format string, a ...interface{}) {
	p := fmt.Sprintf(format, a...)
	log.Print(p)
	// debug.PrintStack()
}
