package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	db "wreis/db/lib"
	"wreis/session"
	util "wreis/util/lib"
)

// AuthenticateData is the struct with the username and password
// used for authentication
type AuthenticateData struct {
	User       string `json:"user"`
	Pass       string `json:"pass"`
	FLAGS      uint64 `json:"flags"`
	UserAgent  string `json:"useragent"`
	RemoteAddr string `json:"remoteaddr"`
}

// SvcAuthenticate handles authentication requests from clients.
//
//	wsdoc {
//	 @Title Authenticate
//	 @URL /v1/authn
//	 @Method  POST
//	 @Synopsis Authenticate a user
//	 @Descr It contacts Accord Directory server to authenticate users. If successful,
//	 @Descr it creates a session for the user and sends a response with Status
//	 @Descr set to "success".  If it is not successful, it sends  response
//	 @Descr with Status set to "error" and provides the reason as given by
//	 @Descr the Accord Directory server.
//	 @Input AuthenticateData
//	 @Response SvcStatus
//
// wsdoc }
// -----------------------------------------------------------------------------
func SvcAuthenticate(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var funcname = "SvcAuthenticate"
	var a AuthenticateData
	var b session.ValidateCookieResponse

	util.Console("Entered %s\n", funcname)
	// util.Console("record data = %s\n", d.data)  // this has user's password, so try not to print (you may forget to remove it)

	if err := json.Unmarshal([]byte(d.data), &a); err != nil {
		e := fmt.Errorf("%s: error with json.Unmarshal:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	//-----------------------------------------------------------------------
	// fill in what the auth server needs...
	//-----------------------------------------------------------------------
	a.RemoteAddr = r.RemoteAddr // this needs to be the user's value, not our server's value
	a.UserAgent = r.UserAgent() // this needs to be the user's value, not our server's value
	fwdaddr := r.Header.Get("X-Forwarded-For")
	util.Console("Forwarded-For address: %q\n", fwdaddr)
	if len(fwdaddr) > 0 {
		a.RemoteAddr = fwdaddr
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
	url := db.Wdb.Config.AuthNHost + "v1/authenticate"
	util.Console("posting request to: %s\n", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pbr))
	if err != nil {
		e := fmt.Errorf("%s: failed to post:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	util.Console("\n*** req = %#v\n\n", req)
	client := &http.Client{}
	util.Console("\n*** client = %#v\n\n", client)
	resp, err := client.Do(req)
	if err != nil {
		e := fmt.Errorf("%s: failed to execute client.Do:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	util.Console("response Body: %s\n", string(body))

	if err := json.Unmarshal([]byte(body), &b); err != nil {
		e := fmt.Errorf("%s: error with json.Unmarshal:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	util.Console("Successfully unmarshaled response: %s\n", string(body))

	switch b.Status {
	case "success":
		util.Console("Authentication succeeded\n")
	default:
		e := fmt.Errorf("%s: unexpected response from authentication service:  %s", funcname, b.Status)
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	util.Console("Directory Service Expire time = %s\n", time.Time(b.Expire).Format(util.RRDATETIMEINPFMT))
	s, err := session.CreateSession(r.Context(), &b)
	if err != nil {
		SvcFuncErrorReturn(w, err, funcname)
		return
	}
	util.Console("Session Created.  b = %#v\n", b)
	cookie := http.Cookie{Name: session.SessionCookieName, Value: b.Token, Expires: s.Expire, Path: "/"}

	http.SetCookie(w, &cookie) // a cookie cannot be set after writing anything to a response writer
	b.ImageURL = s.ImageURL
	b.Name = s.Name
	util.Ulog("user %s (%d) logged in\n", s.Name, s.UID)
	util.Console("Session Table:\n")
	session.DumpSessions()
	util.Console("Created session: %#v\n", s)
	util.Console("Created response: %#v\n", b)
	SvcWriteResponse(&b, w)
}
