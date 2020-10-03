package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// UserInfoResponse is the directory server's response to a user info request
type UserInfoResponse struct {
	Status  string   `json:"status"`
	Record  UserInfo `json:"record"`
	Message string
}

// UserInfoRequest has all the command info needed to make a request for
// user information to the Directory Service.
//-----------------------------------------------------------------------------
type UserInfoRequest struct {
	Cmd string `json:"cmd"` // get, save, delete
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
