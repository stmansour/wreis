package ws

import (
	"fmt"
	"mojo/util"
	"net/http"
	db "wreis/db/lib"
	"wreis/session"
)

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// Dashboard contains the data fields for the dashboard
type Dashboard struct {
	PropertyCount       int64
	CompletedProperties int64
	ActiveProperties    int64
	YourQueue           int64
}

// GetDashboard is the struct used to xfer Dashboard data to a requester
type GetDashboard struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Record  Dashboard `json:"record"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SvcHandlerDashboard returns dashboard stats
// with the w2ui Form
//
// The server command can be:
//      get
//------------------------------------------------------------------------------
func SvcHandlerDashboard(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerDashboard, d.ID = %d\n", d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		getDashboard(w, r, d)
		break
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// getDashboard returns the Dashboard stats
// wsdoc {
//  @Title  Get Dashboard
//	@URL /v1/Dashboard/:PRID
//  @Method  GET
//	@Synopsis Get Dashboard statistics
//  @Description  Return all fields for Dashboard :PRID
//	@Input WebGridSearchRequest
//  @Response GetDashboard
// wsdoc }
//-----------------------------------------------------------------------------
func getDashboard(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getDashboard"
	util.Console("entered %s\n", funcname)

	var g GetDashboard

	row := db.Wdb.DB.QueryRow(`SELECT
	COUNT(PRID) AS Total,
	COUNT(CASE WHEN FlowState=7 THEN 1 END) AS Closed,
    COUNT(CASE WHEN FlowState IN (1,2,3,4,5,6) THEN 1 END) AS Open
FROM Property`)
	if err := row.Scan(&g.Record.PropertyCount, &g.Record.CompletedProperties, &g.Record.ActiveProperties); err != nil {
		SvcErrorReturn(w, err)
		return
	}

	var sess *session.Session
	var ok bool
	if sess, ok = session.GetSessionFromContext(r.Context()); !ok {
		SvcErrorReturn(w, db.ErrSessionRequired)
		return
	}

	q := fmt.Sprintf(`SELECT DISTINCT COUNT(Property.PRID) FROM Property
LEFT JOIN StateInfo ON (StateInfo.FlowState = Property.FlowState)
WHERE (StateInfo.FLAGS & 0x4) =0 AND ((%d = StateInfo.OwnerUID AND (StateInfo.FLAGS & 2)=0) OR (%d = ApproverUID AND (StateInfo.FLAGS & 2)=2));
`, sess.UID, sess.UID)
	row = db.Wdb.DB.QueryRow(q)
	if err := row.Scan(&g.Record.YourQueue); err != nil {
		SvcErrorReturn(w, err)
		return
	}

	g.Status = "success"
	SvcWriteResponse(&g, w)
}
