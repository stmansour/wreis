package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	db "wreis/db/lib"
	util "wreis/util/lib"
)

// UserInfo is the data the directory server returns to us for the UID
type UserInfo struct {
	UID           int64
	FirstName     string
	MiddleName    string
	LastName      string
	PreferredName string
}

// UserInfoTD is the data the directory server returns to us for the UID
type UserInfoTD struct {
	Recid int64 `json:"recid"` // this will hold the UID
	UID   int64
	Name  string
}

// UserInfoResponse is the directory server's response to a user info request
type UserInfoResponse struct {
	Status  string   `json:"status"`
	Record  UserInfo `json:"record"`
	Message string
}

// UserInfoResponseTD is the directory server's response to a user info request
type UserInfoResponseTD struct {
	Status  string       `json:"status"`
	Records []UserInfoTD `json:"records"`
	Message string
}

// UserListInfoResponse is the directory server's response to a user info request
type UserListInfoResponse struct {
	Status  string     `json:"status"`
	Total   int64      `json:"total"`
	Records []UserInfo `json:"records"`
	Message string
}

// UserInfoRequest has all the command info needed to make a request for
// user information to the Directory Service.
//-----------------------------------------------------------------------------
type UserInfoRequest struct {
	Cmd  string `json:"cmd"` // get, save, delete
	UIDs []int64
}

// SvcUserTypeDown handles typedown requests for Users.  It returns
// FirstName, LastName, and TCID
//-----------------------------------------------------------------------------
func SvcUserTypeDown(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	const funcname = "SvcUserTypeDown"
	var g UserInfoResponseTD
	// var err error

	util.Console("Entered %s\n", funcname)
	// util.Console("handle typedown: GetUsersTypeDown( id=%d, search=%s, limit=%d\n", d.ID, d.wsTypeDownReq.Search, d.wsTypeDownReq.Max)

	// build a url like this to send to the directory service:
	//  http://wherever/v1/peopletd?request=%7B%22search%22%3A%22m%22%2C%22max%22%3A100%7D
	q := fmt.Sprintf(`request={"search":%q,"max":10}`, d.wsTypeDownReq.Search)
	url := fmt.Sprintf("%sv1/peopletd?%s", db.Wdb.Config.AuthNHost, url.QueryEscape(q))
	util.Console("search = %s\nurl = %s\n", d.wsTypeDownReq.Search, url)

	resp, err := http.Get(url)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	defer resp.Body.Close()

	// util.Console("Response status: %s\n", resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}

	// util.Console("body = %s\n", string(body))

	if err = json.Unmarshal(body, &g); err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	for i := 0; i < len(g.Records); i++ {
		g.Records[i].Recid = g.Records[i].UID
	}

	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// GetUserInfo contacts the directory service and gets information about
// the user with the supplied UID.
//
// INPUTS
//   uid = User ID of person of interest
//
// RETURNS
//   Name information about the person
//   any errors encountered
//-----------------------------------------------------------------------------
func GetUserInfo(uid int64) (UserInfo, error) {
	funcname := "ws.getUserInfo"
	var p UserInfo
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
	url := fmt.Sprintf("%sv1/people/0/%d", db.Wdb.Config.AuthNHost, uid)
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

// GetUserListInfo contacts the directory service and gets information about
// the user with the supplied UID.
//
// INPUTS
//   uids = slice User IDs of persons of interest
//
// RETURNS
//   Name information about the users in the uid list
//   any errors encountered
//-----------------------------------------------------------------------------
func GetUserListInfo(uids []int64) ([]UserInfo, error) {
	funcname := "ws.getUserInfo"
	var g []UserInfo
	var r = UserInfoRequest{Cmd: "getlist"}
	r.UIDs = append(r.UIDs, uids...)
	b, err := json.Marshal(&r)
	if err != nil {
		e := fmt.Errorf("Error marshaling json data: %s", err.Error())
		util.Ulog("%s: %s\n", funcname, err.Error())
		return g, e
	}

	//----------------------------------------------------------------------
	// the business portion of the URL is ignored.  We snap it to 0
	//----------------------------------------------------------------------
	url := fmt.Sprintf("%sv1/people/0", db.Wdb.Config.AuthNHost)
	// util.Console("userInfo request: %s\n", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return g, fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
	}
	defer resp.Body.Close()

	// util.Console("response Status: %s\n", resp.Status)
	// util.Console("response Headers: %s\n", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// util.Console("response Body: %s\n", string(body))

	var foo UserListInfoResponse
	if err := json.Unmarshal([]byte(body), &foo); err != nil {
		return g, fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
	}

	g = foo.Records

	switch foo.Status {
	case "success":
		return g, nil
	case "error":
		return g, fmt.Errorf("%s", foo.Message)
	default:
		return g, fmt.Errorf("%s: Unexpected response from authentication service:  %s", funcname, foo.Status)
	}
}
