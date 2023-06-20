package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	db "wreis/db/lib"
	"wreis/session"
	util "wreis/util/lib"
)

// ValidateCookie describes the data sent by an AIR app to check
// whether or not a cookie value is valid.
type ValidateCookie struct {
	CookieVal string `json:"cookieval"`
	FLAGS     uint64 `json:"flags"`
}

// SvcLogoff handles authentication requests from clients.
//
//	@Title Logoff
//	@URL /v1/logoff
//	@Method  POST
//	@Synopsis Logoff a user
//	@Descr It removes the user's session from the session table if it exists
//	@Input n/a
//	@Response SvcStatus
//
// wsdoc }
// -----------------------------------------------------------------------------
func SvcLogoff(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var funcname = "SvcLogoff"
	util.Console("Entered %s\n", funcname)
	if d.sess == nil {
		util.Console("%s:  d.sess is nil\n", funcname)
		err := fmt.Errorf("%s: cannot delete nil session", funcname)
		SvcFuncErrorReturn(w, err, funcname)
		return
	}

	//-----------------------------------------------------------------
	// If we get this far, it means that we do have a session (d.sess)
	// Just delete the session.  This will also expire the cookie
	//-----------------------------------------------------------------
	if nil != d.sess {
		session.Delete(d.sess, w, r)
	}

	// The logoff command uses the same data struct as ValidateCookie
	var a = ValidateCookie{
		CookieVal: d.sess.Token, // this is the cookie val we want to delete
	}

	//-----------------------------------------------------------------------
	// Marshal together a new request buffer...
	//-----------------------------------------------------------------------
	pbr, err := json.Marshal(&a)
	if err != nil {
		e := fmt.Errorf("error marshaling json data: %s", err.Error())
		util.Ulog("%s: %s\n", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	util.Console("Request to auth server:  %s\n", string(pbr))

	//-----------------------------------------------------------------------
	// Send to the authenication server
	//-----------------------------------------------------------------------
	url := db.Wdb.Config.AuthNHost + "v1/logoff"
	util.Console("posting request to: %s\n", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pbr))
	if err != nil {
		e := fmt.Errorf("%s: failed to execute post request:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e := fmt.Errorf("%s: failed to execute client.Do:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	util.Console("response Body: %s\n", string(body))

	var b SvcStatus
	if err := json.Unmarshal([]byte(body), &b); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	util.Console("Status response: %s\n", b.Status)
	SvcWriteSuccessResponse(w)
	util.Ulog("user %s logged off\n", d.sess.Username)
}
