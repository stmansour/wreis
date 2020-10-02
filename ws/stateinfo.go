package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	db "wreis/db/lib"
	util "wreis/util/lib"
)

// StateInfoGrid contains the data from StateInfo that is targeted to the UI Grid that displays
// a list of StateInfo structs

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// StateInfo is an individual member of the rent step list.
type StateInfo struct {
	Recid        int64             `json:"recid"`
	SIID         int64             // unique id for this record
	PRID         int64             // id of property to which this record belongs
	InitiatorUID int64             // date/time this state was initiated
	InitiatorDt  util.JSONDateTime // date/time this state was initiated
	ApproverUID  int64             // date/time this state was approved
	ApproverDt   util.JSONDateTime // date/time this state was approved
	FlowState    int64             // state being described
	FLAGS        uint64            // 1<<0 :  0 -> Opt is valid, 1 -> Dt is valid
	LastModTime  util.JSONDateTime // when was the record last written
	LastModBy    int64             // id of user that did the modify
	CreateTime   util.JSONDateTime // when was this record created
	CreateBy     int64             // id of user that created it
}

// SearchStateInfoResponse is the response data for a Rental Agreement Search
type SearchStateInfoResponse struct {
	Status  string      `json:"status"`
	Total   int64       `json:"total"`
	Records []StateInfo `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveStateInfo is sent to save one of open time slots as a reservation
type SaveStateInfo struct {
	Cmd     string      `json:"cmd"`
	Records []StateInfo `json:"records"`
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetStateInfo is the struct returned on a request for a reservation.
type GetStateInfo struct {
	Status  string      `json:"status"`
	Records []StateInfo `json:"records"`
}

//-----------------------------------------------------------------------------
//#############################################################################
//-----------------------------------------------------------------------------

// SvcHandlerStateInfo formats a complete data record for an assessment for use
// with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
//    /v1/StateInfo/PRID
//
// The server command can be:
//      get
//      save
//      delete
//------------------------------------------------------------------------------
func SvcHandlerStateInfo(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerStateInfo, d.ID = %d\n", d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchStateInfo(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcErrorReturn(w, fmt.Errorf("StateInfoID is required but was not specified"))
				return
			}
			getStateInfo(w, r, d)
		}
		break
	case "save":
		saveStateInfo(w, r, d)
		break
	case "delete":
		deleteStateInfo(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// SvcSearchStateInfo generates a report of all StateInfo defined business d.BID
// wsdoc {
//  @Title  Search StateInfo
//	@URL /v1/StateInfo/[:GID]
//  @Method  POST
//	@Synopsis Search StateInfo
//  @Descr  Search all StateInfo and return those that match the Search Logic.
//  @Descr  The search criteria includes start and stop dates of interest.
//	@Input WebGridSearchRequest
//  @Response StateInfoSearchResponse
// wsdoc }
//-----------------------------------------------------------------------------
func SvcSearchStateInfo(w http.ResponseWriter, r *http.Request, d *ServiceData) {
}

// deleteStateInfo deletes a payment type from the database
// wsdoc {
//  @Title  Delete StateInfo
//	@URL /v1/StateInfo/PID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a StateInfo.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func deleteStateInfo(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteStateInfo"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	SvcWriteSuccessResponse(w)
}

// SaveStateInfo expects as input a full state definition (an array of StateInfo
//     structs that describes the state of a property).  It will efficiently
//     add / update / delete StateInfo records so that the database reflects the
//     array supplied
//
//	@URL /v1/StateInfo/PRID
//
///-----------------------------------------------------------------------------
func saveStateInfo(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveStateInfo"
	util.Console("Entered %s\n", funcname)

	util.Console("record data = %s\n", d.data)

	var foo SaveStateInfo
	data := []byte(d.data)
	err := json.Unmarshal(data, &foo)

	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}

	//---------------------------------------
	// start transaction
	//---------------------------------------
	tx, ctx, err := db.NewTransactionWithContext(r.Context())
	if err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
	}

	util.Console("read %d StateInfo\n", len(foo.Records))

	//------------------------------------------------------------------------
	// Read the existing StateInfo list first. We will compare what's being
	// saved to what's already in the database and only perform the changes
	// needed.  In this url, d.ID is the PRID
	//------------------------------------------------------------------------
	var a []db.StateInfo
	if d.ID > 0 {
		a, err = db.GetAllStateInfoItems(ctx, d.ID)
		if err != nil {
			util.Console("%s: B\n", funcname)
			tx.Rollback()
			SvcErrorReturn(w, err)
			return
		}
	}

	util.Console("%s: PRID = %d  num items = %d\n", funcname, d.ID, len(a))

	//------------------------------------------------------------------------
	// Now loop through and compare what was sent to what is in the db...
	//------------------------------------------------------------------------
	for i := 0; i < len(foo.Records); i++ {
		if foo.Records[i].SIID < 0 { // SIID < 0 means it's a new entry
			var x db.StateInfo
			util.Console("ADD: foo.Records[i].SIID = %d\n", foo.Records[i].SIID)
			util.MigrateStructVals(&foo.Records[i], &x)
			x.PRID = d.ID
			util.Console("\tx = %#v\n", x)
			if _, err = db.InsertStateInfo(ctx, &x); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
		} else {
			for j := 0; j < len(a); j++ {
				if foo.Records[i].SIID == a[j].SIID { // if the SIIDs match, compare the values...
					//----------------------------------------
					// Compare the two StateInfo items...
					//----------------------------------------
					if (foo.Records[i].InitiatorUID != a[j].InitiatorUID) ||
						(foo.Records[i].ApproverUID != a[j].ApproverUID) ||
						(foo.Records[i].FLAGS != a[j].FLAGS) ||
						util.EqualDtToJSONDateTime(&a[j].InitiatorDt, &foo.Records[i].InitiatorDt) ||
						util.EqualDtToJSONDateTime(&a[j].ApproverDt, &foo.Records[i].ApproverDt) { // if anything relevant changed
						//----------------------------------------
						// update a[j] to db...
						//----------------------------------------
						util.Console("MOD: foo.Records[i].SIID = %d\n", foo.Records[i].SIID)
						a[j].InitiatorUID = foo.Records[i].InitiatorUID
						a[j].ApproverUID = foo.Records[i].ApproverUID
						a[j].InitiatorDt = time.Time(foo.Records[i].InitiatorDt)
						a[j].ApproverDt = time.Time(foo.Records[i].ApproverDt)
						a[j].FLAGS = foo.Records[i].FLAGS
						if err = db.UpdateStateInfo(ctx, &a[j]); err != nil {
							tx.Rollback()
							SvcErrorReturn(w, err)
							return
						}
					}
				}
			}
		}
	}
	//-------------------------------------------------------------------
	// Now we need to check for any StateInfo that may have been removed
	//-------------------------------------------------------------------
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(foo.Records); j++ {
			if a[i].SIID == foo.Records[j].SIID {
				found = true
				break
			}
		}
		if !found {
			util.Console("DEL: a[i].SIID = %d\n", a[i].SIID)
			if err = db.DeleteStateInfo(ctx, a[i].SIID); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
		}
	}

	//---------------------------------------
	// commit
	//---------------------------------------
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
	}
	SvcWriteSuccessResponse(w)
}

// StateInfoUpdate updates the supplied StateInfo in the database with the supplied
// info. It only allows certain fields to be updated.
//-----------------------------------------------------------------------------
func StateInfoUpdate(p *StateInfo, d *ServiceData) error {
	util.Console("entered StateInfoUpdate\n")
	return nil
}

// GetStateInfo returns the list of StateInfo items associated with the supplied
// PRID.
//
// wsdoc {
//  @Title  Get StateInfo
//	@URL /v1/StateInfo/PRID
//  @Method  GET
//	@Synopsis Get the list of StateInfo
//  @Description  Return all StateInfo items for PRID
//	@Input WebGridSearchRequest
//  @Response GetStateInfo
// wsdoc }
//-----------------------------------------------------------------------------
func getStateInfo(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getStateInfo"
	util.Console("entered %s\n", funcname)
	var g GetStateInfo

	// util.Console("%s: A\n", funcname)
	a, err := db.GetAllStateInfoItems(r.Context(), d.ID)
	if err != nil {
		// util.Console("%s: B\n", funcname)
		SvcErrorReturn(w, err)
		return
	}
	// util.Console("%s: C.  num items = %d\n", funcname, len(a))
	for i := 0; i < len(a); i++ {
		var gg StateInfo
		// util.Console("%s: C.1  a[i] = %#v\n", funcname, a[i])
		util.MigrateStructVals(&a[i], &gg)
		// util.Console("%s: C.2  gg = %#v\n", funcname, gg)
		gg.Recid = gg.SIID
		g.Records = append(g.Records, gg)
		// util.Console("%s: C.3  len(g.Records) = %d\n", funcname, len(g.Records))
	}
	util.Console("%s: D\n", funcname)

	g.Status = "success"
	SvcWriteResponse(&g, w)
}
