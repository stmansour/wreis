package ws

// These are general utilty routines to support w2ui grid components.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	db "wreis/db/lib"
	"wreis/session"
	util "wreis/util/lib"
)

// SvcErrorResponse is the generalized error structure to return errors to the grid widget
type SvcErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SvcStatusResponse is the response to return status when no other data
// needs to be returned
type SvcStatusResponse struct {
	Status string `json:"status"` // typically "success"
	Recid  int64  `json:"recid"`  // set to id of newly inserted record
}

// SvcStatus is the general success return value when no data is required
type SvcStatus struct {
	Status string `json:"status"`
}

// GenSearch describes a search condition
type GenSearch struct {
	Field    string `json:"field"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}

// ColSort is what the UI uses to indicate how the return values should be sorted
type ColSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// WebGridSearchRequestJSON is a struct suitable for describing a webservice operation.
// It is the wire format data. It will be merged into another object where JSONTime values
// are converted to time.Time
type WebGridSearchRequestJSON struct {
	Cmd         string      `json:"cmd"`         // get, save, delete
	Limit       int         `json:"limit"`       // max number to return
	Offset      int         `json:"offset"`      // solution set offset
	Selected    []int       `json:"selected"`    // selected rows
	SearchLogic string      `json:"searchLogic"` // OR | AND
	Search      []GenSearch `json:"search"`      // what fields and what values
	Sort        []ColSort   `json:"sort"`        // sort criteria
}

// WebGridSearchRequest is a struct suitable for describing a webservice operation.
type WebGridSearchRequest struct {
	Cmd         string      `json:"cmd"`         // get, save, delete
	Limit       int         `json:"limit"`       // max number to return
	Offset      int         `json:"offset"`      // solution set offset
	Selected    []int       `json:"selected"`    // selected rows
	SearchLogic string      `json:"searchLogic"` // OR | AND
	Search      []GenSearch `json:"search"`      // what fields and what values
	Sort        []ColSort   `json:"sort"`        // sort criteria
	GroupName   string      `json:"groupName"`   // filter on this group name
}

// WebFormRequest is a struct suitable for describing a webservice operation.
type WebFormRequest struct {
	Cmd      string      `json:"cmd"`    // get, save, delete
	Recid    int         `json:"recid"`  // max number to return
	FormName string      `json:"name"`   // solution set offset
	Record   interface{} `json:"record"` // selected rows
}

// WebTypeDownRequest is a search call made by a client while the user is
// typing in something to search for and the expecation is that the solution
// set will be sent back in realtime to aid the user.  Search is a string
// to search for -- it's what the user types in.  Max is the maximum number
// of matches to return.
type WebTypeDownRequest struct {
	Search string `json:"search"`
	Max    int    `json:"max"`
}

// WebGridDelete is a generic command structure returned when records are
// deleted from a grid. the Selected struct will contain the list of ids
// (recids which should map to the record type unique identifier) that are
// to be deleted.
type WebGridDelete struct {
	Cmd      string  `json:"cmd"`
	Selected []int64 `json:"selected"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
}

// ServiceData is the generalized data gatherer for svcHandler. It allows all
// the common data to be centrally parsed and passed to a handler, which may
// need to parse further to get its unique data.  It includes fields for
// common data elements in web svc requests
type ServiceData struct { // position 0 is 'v1'
	Service       string               // the service requested (position 1)
	GID           int64                // which group (position 2)
	ID            int64                // the numeric id parsed from position 3
	DetVal        string               // the string for the 4th param if provided
	wsSearchReq   WebGridSearchRequest // what did the search requester ask for
	wsTypeDownReq WebTypeDownRequest   // fast for typedown
	data          string               // the raw unparsed data as string
	b             []byte               // the raw unparsed bytes
	GetParams     map[string]string    // parameters when HTTP GET is used
	sess          *session.Session     // user session
}

// ServiceHandler describes the handler for all services
type ServiceHandler struct {
	Cmd           string
	AuthNRequired bool
	Handler       func(http.ResponseWriter, *http.Request, *ServiceData)
}

// Svcs is the table of all service handlers
var Svcs = []ServiceHandler{
	{Cmd: "authn", AuthNRequired: false, Handler: SvcAuthenticate},
	{Cmd: "discon", AuthNRequired: true, Handler: SvcDisableConsole},
	{Cmd: "sessions", AuthNRequired: true, Handler: SvcDumpSessions},
	{Cmd: "encon", AuthNRequired: true, Handler: SvcEnableConsole},
	{Cmd: "logoff", AuthNRequired: true, Handler: SvcLogoff}, // it handles properly if session has already timed out
	{Cmd: "property", AuthNRequired: true, Handler: SvcHandlerProperty},
	{Cmd: "ping", AuthNRequired: false, Handler: SvcHandlerPing},
	{Cmd: "rentsteps", AuthNRequired: true, Handler: SvcHandlerRentSteps},
	{Cmd: "renewoptions", AuthNRequired: true, Handler: SvcHandlerRenewOptions},
	{Cmd: "userprofile", AuthNRequired: true, Handler: SvcUserProfile},
}

// SvcHandlerPing is the most basic test that you can run against the server
// see if it is alive and taking requests. It will return its version number.
func SvcHandlerPing(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	fmt.Fprintf(w, "WREIS - Version %s\n", GetVersionNo())
}

// V1ServiceHandler is the main dispatch point for WEB SERVICE requests
//
// The expected input is of the form:
//		request=%7B%22cmd%22%3A%22get%22%2C%22selected%22%3A%5B%5D%2C%22limit%22%3A100%2C%22offset%22%3A0%7D
// This is exactly what the w2ui grid sends as a request.
//
// Decoded, this message looks something like this:
//		request={"cmd":"get","selected":[],"limit":100,"offset":0}
//
// The leading "request=" is optional. This routine parses the basic
// information, then contacts an appropriate handler for more detailed
// processing.  It will set the Cmd member variable.
//
// W2UI sometimes sends requests that look like this:
//      request=%7B%22search%22%3A%22s%22%2C%22max%22%3A250%7D
// using HTTP GET (rather than its more typical POST).  The command decodes to
// this: request={"search":"s","max":250}
//------------------------------------------------------------------------------
func V1ServiceHandler(w http.ResponseWriter, r *http.Request) {
	funcname := "V1ServiceHandler"
	svcDebugTxn(funcname, r)
	var d ServiceData
	d.ID = -1 // indicates it has not been set

	if e := svcGetPayload(w, r, &d); e != nil {
		SvcErrorReturn(w, e)
		return
	}
	sid := -1 // index to the service requested. Initialize to "not found"
	for i := 0; i < len(Svcs); i++ {
		if Svcs[i].Cmd == d.Service {
			sid = i
			break
		}
	}
	if sid < 0 {
		util.Console("**** YIPES! **** %s - Handler not found\n", r.RequestURI)
		e := fmt.Errorf("Service not recognized: %s", d.Service)
		util.Console("***ERROR IN URL***  %s", e.Error())
		SvcErrorReturn(w, e)
		return
	}

	//-----------------------------------------------------------------------
	// Is authentication required for this command?  If so validate that we
	// have a cookie.
	//-----------------------------------------------------------------------
	// util.Console("SVC: A\n")
	if Svcs[sid].AuthNRequired {
		// util.Console("SVC: B\n")
		c, err := session.ValidateSessionCookie(w, r, 0 /* get all user info */) // this updates the expire time
		if err != nil {
			SvcErrorReturn(w, err)
			return
		}
		// util.Console("SVC: C\n")
		if c.Status != "success" {
			//----------------------------------------------------------------------
			// user needs to log in.  If the command was logoff, just return success.
			// Otherwise return error that they need to login.
			//----------------------------------------------------------------------
			if Svcs[sid].Cmd == "logoff" {
				SvcWriteSuccessResponse(w)
				return
			}
			SvcErrorReturn(w, db.ErrSessionRequired)
			return
		}
		// util.Console("SVC: D\n")
		//----------------------------------------------------------------------
		// The air cookie is valid.  Create (or get) the internal session. This
		// is needed to identify the person associated with the request. All
		// database updates performed by this user will be updated captured
		// in the database writes/updates. We maintain info about this user
		// so we don't have to look it up every time they make a db change.
		//----------------------------------------------------------------------
		if d.sess, err = session.GetSession(r.Context(), w, r, &c); err != nil {
			SvcErrorReturn(w, err)
			return
		}
		// util.Console("\n\n**** %s: d.CookieUpdated = %t ****\n\n", funcname, d.CookieUpdated)
		//----------------------------------------------------------------------
		// If we make it here, we have a good session.  Add it to the request
		// context
		//----------------------------------------------------------------------
		ctx := session.SetSessionContextKey(r.Context(), d.sess)
		r = r.WithContext(ctx)

		// util.Console("SVC: E\n")
	}
	// util.Console("SVC: F\n")

	svcDebugURL(r, &d)
	showRequestHeaders(r)
	Svcs[sid].Handler(w, r, &d)
	svcDebugTxnEnd()
}

// svcGetPayload simply pulls the common data out of the request and puts it
// into the expected locations.
//------------------------------------------------------------------------------
func svcGetPayload(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	var err error
	//-----------------------------------------------------------------------
	// pathElements:  0   1            2
	//               /v1/{subservice}/{ID}
	//-----------------------------------------------------------------------
	ss := strings.Split(r.RequestURI[1:], "?") // it could be GET command
	pathElements := strings.Split(ss[0], "/")
	lpe := len(pathElements)
	if lpe > 1 { // look for the requested service
		d.Service = pathElements[1]
	}
	if lpe > 2 { // subservice, if any
		d.ID, err = util.IntFromString(pathElements[2], "bad ID")
		if err != nil {
			return fmt.Errorf("ID in URL is invalid: %s", err.Error())
		}
	}
	switch r.Method {
	case "POST":
		if err = getPOSTdata(w, r, d); err != nil {
			return err
		}
	case "GET":
		if err = getGETdata(w, r, d); err != nil {
			return err
		}
	}
	showWebRequest(d)
	return nil
}

// SvcGetInt64 tries to read an int64 value from the supplied string.
// If it fails for any reason, it sends writes an error message back
// to the caller and returns the error.  Otherwise, it returns an
// int64 and returns nil
//------------------------------------------------------------------------------
func SvcGetInt64(s, errmsg string, w http.ResponseWriter) (int64, error) {
	i, err := util.IntFromString(s, "not an integer number")
	if err != nil {
		err = fmt.Errorf("%s: %s", errmsg, err.Error())
		SvcErrorReturn(w, err)
		return i, err
	}
	return i, nil
}

// // SvcExtractIDFromURI extracts an int64 id value from position pos of the supplied uri.
// // The URI is of the form returned by http.Request.RequestURI .  In particular:
// //
// //	pos:     0    1      2  3
// //  uri:    /v1/rentable/34/421
// //
// // So, in the example uri above, a call where pos = 3 would return int64(421). errmsg
// // is a string that will be used in the error message if the requested position had an
// // error during conversion to int64. So in the example above, pos 3 is the RID, so
// // errmsg would probably be set to "RID"
// //-----------------------------------------------------------------------------
// func SvcExtractIDFromURI(uri, errmsg string, pos int, w http.ResponseWriter) (int64, error) {
// 	var ID = int64(0)
// 	var err error
//
// 	sa := strings.Split(uri[1:], "/")
// 	// util.Console("uri parts:  %v\n", sa)
// 	if len(sa) < pos+1 {
// 		err = fmt.Errorf("Expecting at least %d elements in URI: %s, but found only %d", pos+1, uri, len(sa))
// 		// util.Console("err = %s\n", err)
// 		SvcErrorReturn(w, err)
// 		return ID, err
// 	}
// 	// util.Console("sa[pos] = %s\n", sa[pos])
// 	ID, err = SvcGetInt64(sa[pos], errmsg, w)
// 	return ID, err
// }

// getPostData parses the posted data from the client and stores in d
//-----------------------------------------------------------------------------
func getPOSTdata(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	funcname := "getPOSTdata"
	var err error
	d.b, err = ioutil.ReadAll(r.Body)
	if err != nil {
		e := fmt.Errorf("%s: Error reading message Body: %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return e
	}

	// THIS IS A DEBUG FEATURE, USUALLY TURNED OFF
	if false {
		fname := "./aws-" + time.Now().String()
		err = ioutil.WriteFile(fname, d.b, 0644)
		if err != nil {
			util.Console("Error with ioutil.WriteFile: %s\n", err.Error())
		}
	}

	util.Console("\t- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	util.Console("\td.b = %s\n", d.b)
	if len(d.b) == 0 {
		d.wsSearchReq.Cmd = "?"
		return nil
	}
	u, err := url.QueryUnescape(string(d.b))
	if err != nil {
		e := fmt.Errorf("%s: Error with QueryUnescape: %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return e
	}
	d.data = strings.TrimPrefix(u, "request=") // strip off "request=" if it is present
	d.b = []byte(d.data)                       // keep the byte array around too
	util.Console("\tUnescaped d.b = %s\n", d.data)

	var wjs WebGridSearchRequestJSON
	err = json.Unmarshal(d.b, &wjs)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return e
	}
	util.MigrateStructVals(&wjs, &d.wsSearchReq)
	return err
}

func getGETdata(w http.ResponseWriter, r *http.Request, d *ServiceData) error {
	funcname := "getGETdata"
	s, err := url.QueryUnescape(strings.TrimSpace(r.URL.String()))
	if err != nil {
		e := fmt.Errorf("%s: Error with url.QueryUnescape:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return e
	}
	util.Console("Unescaped query = %s\n", s)
	w2uiPrefix := "request="
	n := strings.Index(s, w2uiPrefix)
	util.Console("n = %d\n", n)
	if n > 0 {
		util.Console("Will process as Typedown\n")
		d.data = s[n+len(w2uiPrefix):]
		util.Console("%s: will unmarshal: %s\n", funcname, d.data)
		if err = json.Unmarshal([]byte(d.data), &d.wsTypeDownReq); err != nil {
			e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
			SvcErrorReturn(w, e)
			return e
		}
		d.wsSearchReq.Cmd = "typedown"
	} else {
		util.Console("Will process as web search command\n")
		d.wsSearchReq.Cmd = r.URL.Query().Get("cmd")
	}
	return nil
}

func showRequestHeaders(r *http.Request) {
	util.Console("Headers:\n")
	for k, v := range r.Header {
		util.Console("\t%s: ", k)
		for i := 0; i < len(v); i++ {
			util.Console("%q  ", v[i])
		}
		util.Console("\n")
	}
}

func showWebRequest(d *ServiceData) {
	if d.wsSearchReq.Cmd == "typedown" {
		util.Console("Typedown:\n")
		util.Console("\tSearch  = %q\n", d.wsTypeDownReq.Search)
		util.Console("\tMax     = %d\n", d.wsTypeDownReq.Max)
	} else {
		util.Console("\tSearchReq:\n")
		util.Console("\t\tCmd           = %s\n", d.wsSearchReq.Cmd)
		util.Console("\t\tLimit         = %d\n", d.wsSearchReq.Limit)
		util.Console("\t\tOffset        = %d\n", d.wsSearchReq.Offset)
		util.Console("\t\tsearchLogic   = %s\n", d.wsSearchReq.SearchLogic)
		for i := 0; i < len(d.wsSearchReq.Search); i++ {
			util.Console("\t\tsearch[%d] - Field = %s,  Type = %s,  Value = %s,  Operator = %s\n", i, d.wsSearchReq.Search[i].Field, d.wsSearchReq.Search[i].Type, d.wsSearchReq.Search[i].Value, d.wsSearchReq.Search[i].Operator)
		}
		for i := 0; i < len(d.wsSearchReq.Sort); i++ {
			util.Console("\t\tsort[%d] - Field = %s,  Direction = %s\n", i, d.wsSearchReq.Sort[i].Field, d.wsSearchReq.Sort[i].Direction)
		}
	}
}

func svcDebugTxn(funcname string, r *http.Request) {
	util.Console("\n%s\n", util.Mkstr(80, '-'))
	util.Console("URL:      %s\n", r.URL.String())
	util.Console("METHOD:   %s\n", r.Method)
	util.Console("Handler:  %s\n", funcname)
}

func svcDebugURL(r *http.Request, d *ServiceData) {
	//-----------------------------------------------------------------------
	// pathElements: 0         1     2
	// Break up {subservice}/{BUI}/{ID} into an array of strings
	// BID is common to nearly all commands
	//-----------------------------------------------------------------------
	ss := strings.Split(r.RequestURI[1:], "?") // it could be GET command
	pathElements := strings.Split(ss[0], "/")
	util.Console("\t%s\n", r.URL.String()) // print before we strip it off
	for i := 0; i < len(pathElements); i++ {
		util.Console("\t\t%d. %s\n", i, pathElements[i])
	}
}

func svcDebugTxnEnd() {
	util.Console("END\n")
}

// SvcWriteResponse finishes the transaction with the W2UI client
func SvcWriteResponse(g interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(g)
	if err != nil {
		e := fmt.Errorf("Error marshaling json data: %s", err.Error())
		util.Ulog("SvcWriteResponse: %s\n", err.Error())
		SvcErrorReturn(w, e)
		return
	}
	SvcWrite(w, b)
}

// SvcWrite is a general write routine for service calls... it is a bottleneck
// where we can place debug statements as needed.
func SvcWrite(w http.ResponseWriter, b []byte) {
	util.Console("first 300 chars of response: %-300.300s\n", string(b))
	// util.Console("\nResponse Data:  %s\n\n", string(b))
	w.Write(b)
}

// SvcWriteSuccessResponse is used to complete a successful write operation on w2ui form save requests.
func SvcWriteSuccessResponse(w http.ResponseWriter) {
	var g = SvcStatusResponse{Status: "success"}
	w.Header().Set("Content-Type", "application/json")
	SvcWriteResponse(&g, w)
}

// SvcWriteSuccessResponseWithID is used to complete a successful write operation on w2ui form save requests.
func SvcWriteSuccessResponseWithID(w http.ResponseWriter, id int64) {
	var g = SvcStatusResponse{Status: "success", Recid: id}
	w.Header().Set("Content-Type", "application/json")
	SvcWriteResponse(&g, w)
}

// SvcFuncErrorReturn formats an error return to the grid widget and sends it
func SvcFuncErrorReturn(w http.ResponseWriter, err error, funcname string) {
	// util.Console("<Function>: %s | <Error Message>: <<begin>>\n%s\n<<end>>\n", funcname, err.Error())
	util.Console("%s: %s\n", funcname, err.Error())
	var e SvcErrorResponse
	e.Status = "error"
	e.Message = err.Error()
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(e)
	SvcWrite(w, b)
}

// SvcErrorReturn formats an error return to the grid widget and sends it
func SvcErrorReturn(w http.ResponseWriter, err error) {
	var e SvcErrorResponse
	e.Status = "error"
	e.Message = fmt.Sprintf("Error: %s\n", err.Error())
	b, _ := json.Marshal(e)
	SvcWrite(w, b)
}
