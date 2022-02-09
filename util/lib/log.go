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

// LogAndPrintError is Phonebooks's standard logger
func LogAndPrintError(f string, err error) {
	p := fmt.Sprintf("%s: %s", f, err.Error())
	log.Print(p)
	fmt.Print(p)
}
