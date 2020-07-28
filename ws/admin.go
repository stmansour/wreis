package ws

import (
	"net/http"
	"wreis/session"
	util "wreis/util/lib"
)

// SvcDisableConsole disables console messages from printing out
func SvcDisableConsole(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.DisableConsole()
	SvcWriteSuccessResponse(w)
}

// SvcEnableConsole enables console messages to print out
func SvcEnableConsole(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.EnableConsole()
	SvcWriteSuccessResponse(w)
}

// SvcDumpSessions enables console messages to print out
func SvcDumpSessions(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	session.DumpSessions()
	SvcWriteSuccessResponse(w)
}
