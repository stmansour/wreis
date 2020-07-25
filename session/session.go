package session

import (
	"bytes"
	"context"
	"encoding/json"
	"extres"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	util "wreis/util/lib"
)

// WSPerson is the person information we pass over the web service call.
// We can add information as we need it after reviewing the security
// implications.
//-------------------------------------------------------------------------
type WSPerson struct {
	UID           int64
	FirstName     string
	MiddleName    string
	LastName      string
	PreferredName string
}

// UserInfoResponse is the directory server's response to a user info request
type UserInfoResponse struct {
	Status  string   `json:"status"`
	Record  WSPerson `json:"record"`
	Message string
}

// DirectoryPerson is the structure of person in Accord Directory
// with publicly viewable data.
type DirectoryPerson struct {
	UID           int64
	UserName      string
	LastName      string
	MiddleName    string
	FirstName     string
	PreferredName string
	OfficePhone   string
	CellPhone     string
}

// sessSvcStatus is the generalized error structure to return errors to the grid widget
type sessSvcStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// sessSvcStatusResponse is the response to return status when no other data
// needs to be returned
type sessSvcStatusResponse struct {
	Status string `json:"status"` // typically "success"
	Recid  int64  `json:"recid"`  // set to id of newly inserted record
}

// appConfig defines the external resources information
var appConfig extres.ExternalResources

// Session is the structure definition for  a user session
// with this program.
type Session struct {
	Token    string    // this is the md5 hash, unique id
	Username string    // associated username
	Name     string    // user's preferred name if it exists, otherwise the FirstName
	UID      int64     // user's db uid
	CoCode   int64     // logged in user's company (from Accord Directory)
	ImageURL string    // user's picture
	Expire   time.Time // when does the cookie expire
	RoleID   int64     // security role id
}

// UnrecognizedCookie is the error string associated with an
// unrecognized value for the airoller cookie
var UnrecognizedCookie = string("unrecognized cookie")

// ReqSessionMem is the channel int used to request write permission to the session list
var ReqSessionMem chan int

// ReqSessionMemAck is the channel int used to handshake for access to the session list
var ReqSessionMemAck chan int

// CleanupTime defines the time interval between the routine that removes
// expired sessions.
var CleanupTime time.Duration

// sessions is the session list managed by this code
var sessions map[string]*Session

// SessionTimeout defines how long a session can remain idle before it expires.
var SessionTimeout time.Duration // in minutes

// SessionCookieName is the name of the Roller cookie where the session
// token is stored.
var SessionCookieName = string("air")

// Init must be called prior to using the session subsystem. It
// initializes structures and starts the dispatcher
//
// INPUT
//  timeout - the number of minutes before a session times out
//  x       - external resources definition
//
// RETURNS
//  nothing at this time
//-----------------------------------------------------------------------------
func Init(timeout int, x extres.ExternalResources) {
	appConfig = x
	sessions = make(map[string]*Session)
	ReqSessionMem = make(chan int)
	ReqSessionMemAck = make(chan int)
	CleanupTime = time.Duration(1)
	SessionTimeout = time.Duration(timeout) * time.Minute
	go Dispatcher()
	go Cleanup()
}

// GetSessionCookieName simply returns a string containing the session
// cookie name. We want this to be a private / unchangeable name.
//-----------------------------------------------------------------------------
func GetSessionCookieName() string {
	return SessionCookieName
}

// GetSessionTable returns a copy of the session table.  This is for use
// in the administrators command to view the session table of the active
// server.
//-----------------------------------------------------------------------------
func GetSessionTable() map[string]*Session {
	return sessions
}

// Dispatcher is a Go routine that controls access to shared memory.
//-----------------------------------------------------------------------------
func Dispatcher() {
	for {
		select {
		case <-ReqSessionMem:
			ReqSessionMemAck <- 1 // tell caller go ahead
			<-ReqSessionMemAck    // block until caller is done with mem
		}
	}
}

// Cleanup a Go routine to periodically spin through the session list
// and remove any sessions which have timed out.
//-----------------------------------------------------------------------------
func Cleanup() {
	for {
		select {
		case <-time.After(CleanupTime * time.Minute):
			ReqSessionMem <- 1                 // ask to access the shared mem, blocks until granted
			<-ReqSessionMemAck                 // make sure we got it
			ss := make(map[string]*Session, 0) // here's the new Session list
			n := 0                             // total number removed
			now := time.Now()                  // this is the timestamp we'll compare against
			// util.Console("Cleanup time: %s\n", now.Format(RRDATETIMEINPFMT))
			for k, v := range sessions { // look at every Session
				// util.Console("Found session: %s, expire time: %s\n", v.Name, v.Expire.Format(RRDATETIMEINPFMT))
				if now.After(v.Expire) { // if it's still active...
					n++ // removed another
				} else {
					ss[k] = v // ...copy it to the new list
				}
			}
			sessions = ss         // set the new list
			ReqSessionMemAck <- 1 // tell Dispatcher we're done with the data
			// util.Console("Cleanup completed. %d removed. Current Session list size = %d\n", n, len(sessions))
		}
	}
}

// Get returns the session associated with the supplied token, if it
// exists. It may no longer exist because it timed out.
//
// INPUT
//  token -  the index into the session table for this session. This is the
//           value that is stored in a web session cookie
//
// RETURNS
//  session - pointer to the session if the bool is true
//  bool    - true if the session was found, false otherwise
//-----------------------------------------------------------------------------
func Get(token string) (*Session, bool) {
	s, ok := sessions[token]
	return s, ok
}

// ToString is the stringer for sessions
//
// RETURNS
//  a string representation of the session entry
//-----------------------------------------------------------------------------
func (s *Session) ToString() string {
	if nil == s {
		return "nil"
	}
	return fmt.Sprintf("User(%s) Name(%s) UID(%d) Token(%s)",
		s.Username, s.Name, s.UID, s.Token)
}

// DumpSessions sends the contents of the session table to the consol.
//
// RETURNS
//  a string representation of the session entry
//-----------------------------------------------------------------------------
func DumpSessions() {
	i := 0
	for _, v := range sessions {
		util.Console("%2d. %s\n", i, v.ToString())
		i++
	}
}

// New creates a new session and adds it to the session list
//
// INPUT
//  token    - the unique token string. This will be used to index the session
//             list
//  username - the username from the authentication service
//  name     - the name to use in the session
//  uid      - the userid associated with username. From the auth server.
//  rid      - security role id
//
// RETURNS
//  session - pointer to the new session
//-----------------------------------------------------------------------------
func New(token, username, name string, uid int64, imgurl string, rid int64, expire *time.Time) *Session {
	util.Console("New Session:  username = %s, name = %s, UID = %d\n", username, name, uid)
	s := new(Session)
	s.Token = token
	s.Username = username
	s.Name = name
	s.UID = uid
	s.Expire = *expire

	switch appConfig.AuthNType {
	case "Accord Directory":
		s.ImageURL = imgurl
	}

	ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-ReqSessionMemAck // make sure we got it
	sessions[token] = s
	ReqSessionMemAck <- 1 // tell Dispatcher we're done with the data

	return s
}

// CreateSession creates an HTTP Cookie with the token for this session
//
// INPUT
//  w    - where to write the set cookie
//  r - the request where w should look for the cookie
//
// RETURNS
//  session - pointer to the new session
//-----------------------------------------------------------------------------
func CreateSession(ctx context.Context, a *ValidateCookieResponse) (*Session, error) {
	var err error
	expiration := time.Time(a.Expire)

	// expiration := time.Now().Add(15 * time.Minute)

	//----------------------------------------------
	// TODO: lookup username in address book data
	//----------------------------------------------
	// util.Console("Calling GetDirectoryPerson with UID = %d\n", a.UID)
	var c DirectoryPerson
	// err := sessDB.PBsql.GetDirectoryPerson.QueryRow(a.UID).Scan(&c.UID, &c.UserName, &c.LastName, &c.MiddleName, &c.FirstName, &c.PreferredName, &c.PreferredName, &c.OfficePhone, &c.CellPhone)

	if c, err = getUserInfo(a.UID); err != nil {
		return nil, err
	}

	// util.Console("Successfully read info from directory for UID = %d\n", c.UID)

	RoleID := int64(0) // we haven't yet implemented Role
	name := c.FirstName
	if len(c.PreferredName) > 0 {
		name = c.PreferredName
	}
	s := New(a.Token, c.UserName, name, a.UID, a.ImageURL, RoleID, &expiration)
	return s, nil
}

// UserInfoRequest has all the command info needed to make a request for
// user information to the Directory Service.
//-----------------------------------------------------------------------------
type UserInfoRequest struct {
	Cmd string `json:"cmd"` // get, save, delete
}

// getUserInfo creates an HTTP Cookie with the token for this session
//
// INPUT
//  w    - where to write the set cookie
//  r - the request where w should look for the cookie
//
// RETURNS
//  session - pointer to the new session
//-----------------------------------------------------------------------------
func getUserInfo(uid int64) (DirectoryPerson, error) {
	funcname := "session.getUserInfo"
	var p DirectoryPerson
	var r = UserInfoRequest{Cmd: "get"}
	b, err := json.Marshal(&r)
	if err != nil {
		e := fmt.Errorf("Error marshaling json data: %s", err.Error())
		util.Ulog("%s: %s\n", funcname, err.Error())
		return p, e
	}

	//----------------------------------------------------------------------
	// the business portion of the URL is ignored.  We snap it to 0
	//----------------------------------------------------------------------
	url := fmt.Sprintf("%sv1/people/0/%d", appConfig.AuthNHost, uid)
	// util.Console("userInfo request: %s\n", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return p, fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
	}
	defer resp.Body.Close()

	// util.Console("response Status: %s\n", resp.Status)
	// util.Console("response Headers: %s\n", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// util.Console("response Body: %s\n", string(body))

	var foo UserInfoResponse
	if err := json.Unmarshal([]byte(body), &foo); err != nil {
		return p, fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
	}

	// util.Console("before migrate: foo.record = %#v\n", foo.Record)
	util.MigrateStructVals(&foo.Record, &p)

	switch foo.Status {
	case "success":
		return p, nil
	case "error":
		return p, fmt.Errorf("%s", foo.Message)
	default:
		return p, fmt.Errorf("%s: Unexpected response from authentication service:  %s", funcname, foo.Status)
	}
}

// IsUnrecognizedCookieError returns true if the error is UnrecognizedCookie.
//
// INPUT
//  err - the error to check
//
// RETURNS
//  bool - true means it is an UnrecognizedCookie error
//         false means it is not
//-----------------------------------------------------------------------------
func IsUnrecognizedCookieError(err error) bool {
	return strings.Contains(err.Error(), UnrecognizedCookie)
}

// ValidateSessionCookie is used to ensure that the session is still valid.
// Even if the session is found in our internal table, the 'air' is cookie used
// other applications in the suite. Someone may have logged out from a
// different app. If the cookie is not validated, then destroy the session.
// If the cookie token is valid, its expire time will be updated.
//
// INPUTS
//  r  - pointer to the http request, which may be updated after we add the
//       context value to it.
//  d  - our service data struct
//
// RETURNS
//  cookie - the http cookie or nil if it doesn't exist
//  getData - if 0 then only the existence of the cookie is checked
//            if 1 the all the data associated with the directory service
//            cookie is returned -- this includes the UID, the expire time, ...
//-----------------------------------------------------------------------------
func ValidateSessionCookie(r *http.Request, getData int) (ValidateCookieResponse, error) {
	funcname := "ValidateSessionCookie"
	// util.Console("Entered %s\n", funcname)
	var vc ValidateCookieRequest
	var vr ValidateCookieResponse
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		if strings.Contains(err.Error(), "no air cookie in request headers") {
			return vr, nil
		}
		return vr, nil
	}
	vc.CookieVal = c.Value
	vc.IP = r.RemoteAddr
	vc.UserAgent = r.UserAgent()

	if getData > 0 {
		vc.FLAGS |= 1 // set bit 0
	}
	vc.FLAGS |= 1 << 1 // set bit 1 validate AND reset the expire time. This says 15 min (or whatever it is) from now is the new expire time.

	pbr, err := json.Marshal(&vc)
	if err != nil {
		return vr, fmt.Errorf("Error marshaling json data: %s", err.Error())
	}

	//-----------------------------------------------------------------------
	// Send to the authenication server
	//-----------------------------------------------------------------------
	url := appConfig.AuthNHost + "v1/validatecookie"
	// util.Console("posting request to: %s\n", url)
	// util.Console("              data: %s\n", string(pbr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pbr))
	req.Header.Set("Content-Type", "application/json")
	// util.Console("\n*** req = %#v\n\n", req)
	client := &http.Client{}
	// util.Console("\n*** client = %#v\n\n", client)
	resp, err := client.Do(req)
	if err != nil {
		return vr, fmt.Errorf("%s: failed to execute client.Do:  %s", funcname, err.Error())
	}
	defer resp.Body.Close()

	// util.Console("Response status = %s, status code = %d\n", resp.Status, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	// util.Console("*** Directory Service *** response Body: %s\n", string(body))

	if err := json.Unmarshal([]byte(body), &vr); err != nil {
		return vr, fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
	}
	// util.Console("Successfully unmarshaled response: %s\n", string(body))
	// if vr.Status != "success" {
	// 	vr.CookieVal = ""
	// }
	return vr, nil
}

// GetSession returns the session based on the cookie in the supplied
// HTTP connection. If the "air" cookie is valid, it will either find the
// existing session or create a new session.
//
// INPUT
//  r - the request where we look for the cookie
//
// RETURNS
//  session - pointer to the new session
//  error   - any error encountered
//-----------------------------------------------------------------------------
func GetSession(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Session, error) {
	// funcname := "GetSession"
	// var b AIRAuthenticateResponse
	var ok bool
	var sess *Session

	// util.Console("GetSession 1\n")
	// util.Console("\nSession Table:\n")
	DumpSessions()
	// util.Console("\n")
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		// util.Console("GetSession 2\n")
		if strings.Contains(err.Error(), "cookie not present") {
			// util.Console("GetSession 3\n")
			return nil, nil
		}
		// util.Console("GetSession 4\n")
		return nil, err
	}
	// util.Console("GetSession 5\n")
	sess, ok = sessions[cookie.Value]
	if !ok || sess == nil {
		b, err := ValidateSessionCookie(r, 1)
		if err != nil {
			return sess, err
		}
		// util.Console("ValidateSessionCookie returned b = %#v\n", b)
		// util.Console("Directory Service Expire time = %s\n", time.Time(b.Expire).Format(util.RRDATETIMEINPFMT))
		sess, err = CreateSession(ctx, &b)
		if err != nil {
			return nil, err
		}
		cookie := http.Cookie{Name: SessionCookieName, Value: b.Token, Expires: sess.Expire, Path: "/"}
		http.SetCookie(w, &cookie) // a cookie cannot be set after writing anything to a response writer
		// util.Console("*** NEW SESSION CREATED ***\n")
	}
	return sess, nil
}

// Refresh updates the cookie and Session with a new expire time.
//
// INPUT
//  w - where to write the set cookie
//  r - the request where w should look for the cookie
//
// RETURNS
//  session - pointer to the new session
//-----------------------------------------------------------------------------
func (s *Session) Refresh(w http.ResponseWriter, r *http.Request) int {
	cookie, err := r.Cookie(SessionCookieName)
	if nil != cookie && err == nil {
		cookie.Expires = time.Now().Add(SessionTimeout)
		ReqSessionMem <- 1        // ask to access the shared mem, blocks until granted
		<-ReqSessionMemAck        // make sure we got it
		s.Expire = cookie.Expires // update the Session information
		ReqSessionMemAck <- 1     // tell Dispatcher we're done with the data
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		return 0
	}
	return 1
}

// ExpireCookie expires the cookie associated with this session now
//
// INPUT
//  w - where to write the set cookie
//  r - the request where w should look for the cookie
//
// RETURNS
//  nothing at this time
//-----------------------------------------------------------------------------
func (s *Session) ExpireCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionCookieName)
	if nil != cookie && err == nil {
		cookie.Expires = time.Now()
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}
}

// Delete removes the supplied Session.
// if there is a better idiomatic way to do this, please let me know.
// INPUT
//  Session  - pointer to the session to tdelete
//             list
//  w        - where to write the set cookie
//  r        - the request where w should look for the cookie
//
// RETURNS
//  session  - pointer to the new session
//-----------------------------------------------------------------------------
func Delete(s *Session, w http.ResponseWriter, r *http.Request) {
	if nil == s {
		util.Console("Delete: supplied session is nil\n")
		return
	}
	util.Console("Session being deleted: %s\n", s.ToString())
	util.Console("sessions before delete:\n")
	DumpSessions()

	ss := make(map[string]*Session, 0)

	ReqSessionMem <- 1 // ask to access the shared mem, blocks until granted
	<-ReqSessionMemAck // make sure we got it
	for k, v := range sessions {
		if s.Token != k {
			ss[k] = v
		}
	}
	sessions = ss
	ReqSessionMemAck <- 1 // tell Dispatcher we're done with the data
	s.ExpireCookie(w, r)
	util.Console("sessions after delete:\n")
	DumpSessions()
}

// sessSvcWrite is a general write routine for service calls... it is a bottleneck
// where we can place debug statements as needed.
func sessSvcWrite(w http.ResponseWriter, b []byte) {
	charsToPrint := 500
	format := fmt.Sprintf("First %d chars of response: %%-%d.%ds\n", charsToPrint, charsToPrint, charsToPrint)
	// util.Console("Format string = %q\n", format)
	util.Console(format, string(b))
	// util.Console("\nResponse Data:  %s\n\n", string(b))
	w.Write(b)
}

// sessSvcWriteResponse finishes the transaction with the W2UI client
func sessSvcWriteResponse(g interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(g)
	if err != nil {
		e := fmt.Errorf("Error marshaling json data: %s", err.Error())
		util.Ulog("SvcWriteResponse: %s\n", err.Error())
		sessSvcErrorReturn(w, e, "sessSvcWriteResponse")
		return
	}
	sessSvcWrite(w, b)
}

// sessSvcErrorReturn formats an error return to the grid widget and sends it
func sessSvcErrorReturn(w http.ResponseWriter, err error, funcname string) {
	// util.Console("<Function>: %s | <Error Message>: <<begin>>\n%s\n<<end>>\n", funcname, err.Error())
	util.Console("%s: %s\n", funcname, err.Error())
	var e sessSvcStatus
	e.Status = "error"
	e.Message = err.Error()
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(e)
	sessSvcWrite(w, b)
}

// sessSvcSuccessResponse is used to complete a successful write operation on w2ui form save requests.
func sessSvcSuccessResponse(w http.ResponseWriter) {
	var g = sessSvcStatusResponse{Status: "success"}
	w.Header().Set("Content-Type", "application/json")
	sessSvcWriteResponse(&g, w)
}

// Check encapsulates 6 lines of code that was repeated in every call
//
// INPUTS
//  ctx  the context, which should have session
//
// RETURNS
//  the session
//  ok == true - session was required but not found
//        false - session was found or session not required
//-----------------------------------------------------------------------------
func Check(ctx context.Context) (*Session, bool) {
	return GetSessionFromContext(ctx)
}
