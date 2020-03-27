package util

import (
	"fmt"
	"strings"
	"time"
)

// DateToString rounds the supplied amount to the nearest cent.
func DateToString(t time.Time) string {
	return t.Format("01/02/2006")
}

// AcceptedDateFmts is the array of string formats that StringToDate accepts
var AcceptedDateFmts = []string{
	RRDATEINPFMT,
	RRDATEFMT2,
	RRDATEFMT,
	RRDATEFMT3,
	RRJSUTCDATETIME,
	RRDATETIMEW2UIFMT,
	RRDATETIMEINPFMT,
	RRDATETIMEFMT,
	RRDATEREPORTFMT,
	RRDATERECEIPTFMT,
}

// StringToDate tries to convert the supplied string to a time.Time value. It will use the
// formats called out in dbtypes.go:  RRDATEFMT, RRDATEINPFMT, RRDATEINPFMT2, ...
//
// for further experimentation, try: https://play.golang.org/p/JNUnA5zbMoz
//----------------------------------------------------------------------------------
func StringToDate(s string) (time.Time, error) {
	// try the ansi std date format first
	var Dt time.Time
	var err error
	s = strings.TrimSpace(s)
	for i := 0; i < len(AcceptedDateFmts); i++ {
		Dt, err = time.Parse(AcceptedDateFmts[i], s)
		if nil == err {
			return Dt, nil
		}
	}
	return Dt, fmt.Errorf("Date could not be decoded")
}
