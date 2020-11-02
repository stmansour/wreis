package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	db "wreis/db/lib"
	"wreis/session"
	util "wreis/util/lib"
)

const (
	fn = "?? unknown"
	ln = "user ??"
)

// StateInfo deals with state changes to the StateInfo structs. The states are:
//
//		1. New Record Entry
//		2. First Task List
//		3. Marketing Package
//		4. Ready To List
//		5. Listed
//		6. Under Contract
//		7. Closed
//
// The following actions apply to states:
//
// 		Approve          Approve this state and advance to the next state
//      Reject           Do not advance the state, add a reason for rejection
//      Ready            Lets the approver know the owner believes the property
//                       is ready to progress to the next state
//      Assign Approver  Set the approver
//      Assign Owner     Set the owner
//      Revert           Set the current state of the property back one
//
//

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// StateInfo is an individual member of the rent step list.
type StateInfo struct {
	Recid         int64             `json:"recid"`
	SIID          int64             // unique id for this record
	PRID          int64             // id of property to which this record belongs
	OwnerUID      int64             // uid of Owner
	OwnerDt       util.JSONDateTime // date/time this state was initiated
	OwnerName     string            //
	ApproverUID   int64             // date/time this state was approved
	ApproverDt    util.JSONDateTime // date/time this state was approved
	ApproverName  string            //
	FlowState     int64             // state being described
	Reason        string            // if rejected, why
	FLAGS         uint64            // 1<a :  0 -> Opt is valid, 1 -> Dt is valid
	LastModTime   util.JSONDateTime // when was the record last written
	LastModBy     int64             // id of user that did the modify
	CreateTime    util.JSONDateTime // when was this record created
	CreateBy      int64             // id of user that created it
	CreateByName  string            // creator name
	LastModByName string            // modifier name
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
	case "reject":
		saveStateReject(w, r, d)
		break
	case "revert":
		saveStateRevert(w, r, d)
		break
	case "approve":
		saveStateApprove(w, r, d)
		break
	case "ready":
		saveStateReady(w, r, d)
		break
	case "save":
		saveStateInfo(w, r, d)
		break
	case "setowner":
		saveStateOwner(w, r, d)
		break
	case "setapprover":
		saveStateApprover(w, r, d)
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

func stateInfoHelper(w http.ResponseWriter, r *http.Request, d *ServiceData) (SaveStateInfo, db.StateInfo, *session.Session, error) {
	var foo SaveStateInfo
	var si db.StateInfo
	var sess *session.Session
	var ok bool

	data := []byte(d.data)
	err := json.Unmarshal(data, &foo)

	if err != nil {
		return foo, si, sess, fmt.Errorf("Error with json.Unmarshal:  %s", err.Error())
	}

	if len(foo.Records) != 1 {
		return foo, si, sess, fmt.Errorf("Error with json.Unmarshal:  %s", err.Error())
	}

	//---------------------------------------------------------------------
	// if the id making this save is NOT the approver, then send back an error
	//---------------------------------------------------------------------
	if sess, ok = session.GetSessionFromContext(r.Context()); !ok {
		return foo, si, sess, db.ErrSessionRequired
	}

	//---------------------------------------------------------------------
	// get the SIID so we know we start with the last saved version...
	//---------------------------------------------------------------------
	if si, err = db.GetStateInfo(r.Context(), foo.Records[0].SIID); err != nil {
		return foo, si, sess, err
	}

	//---------------------------------------------------------------------
	// make sure we're not working on something that's already completed
	//---------------------------------------------------------------------
	if si.FLAGS&0x4 != 0 {
		return foo, si, sess, fmt.Errorf("This is not the latest state information")
	}

	return foo, si, sess, nil
}

// saveStateApprover sets approver of the state.  Anyone can do it. The person making
// the change will be noted.
//
//  ANYONE CAN CHANGE THE APPROVER:  But the person who made the change will be
//         kept in the audit trail. Changer will be the UID
//         of LastModBy on this StateInfo, and creator of the StateInfo
//         with new owner
//
//	@URL /v1/StateInfo/PRID
//
//-----------------------------------------------------------------------------
func saveStateApprover(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	foo, si, _, err := stateInfoHelper(w, r, d)
	if err != nil {
		SvcErrorReturn(w, err)
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

	//--------------------------------------------------------------------------
	// mark this version as finished due to approver change...
	// 		0  valid only when ApproverUID > 0, 0 = State Approved, 1 = not approved
	// 	*	1  0 = no request is being made,       1 = request approval for this state
	// 		2  0 = this state is work in progress, 1 = work is concluded on this StateInfo
	//      3  0 = not reverted, 1 = reverted
	//      4  0 = no owner change, 1 = owner change, changer will be the UID
	//         of LastModBy on this StateInfo, and creator of the StateInfo
	//         with new owner
	//  *   5  0 = no approver change, 1 = approver changed.  changer will be the UID
	//         of LastModBy on this StateInfo, and creator of the StateInfo
	//         with new approver
	// we'll set the lower byte to (not approved, no request being made, work concluded )
	//--------------------------------------------------------------------------
	a := si
	a.FLAGS &= si.FLAGS & 0xefffffffffffffc0
	util.Console("before: a.FLAGS = %x\n", a.FLAGS)
	a.FLAGS |= 0x24
	a.Reason = foo.Records[0].Reason // save the reason
	util.Console("after: a.FLAGS = %x\n", a.FLAGS)
	if err = db.UpdateStateInfo(ctx, &a); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
	}

	//--------------------------------------------------------------------------
	// now create the new owner state info...
	//--------------------------------------------------------------------------
	a.SIID = 0
	a.FLAGS = 0
	a.ApproverUID = foo.Records[0].ApproverUID // this has the new owner
	a.Reason = ""
	if _, err = db.InsertStateInfo(ctx, &a); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
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

// saveStateOwner sets owner of the state.  Anyone can do it. The person making
// the change will be noted.
//
//  ANYONE CAN CHANGE THE OWNER:  But the person who made the change will be
//         kept in the audit trail. Changer will be the UID
//         of LastModBy on this StateInfo, and creator of the StateInfo
//         with new owner
//
//	@URL /v1/StateInfo/PRID
//
//-----------------------------------------------------------------------------
func saveStateOwner(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	foo, si, _, err := stateInfoHelper(w, r, d)
	if err != nil {
		SvcErrorReturn(w, err)
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

	//--------------------------------------------------------------------------
	// mark this version as finished due to owner change...
	// 		0  valid only when ApproverUID > 0, 0 = State Approved, 1 = not approved
	// 	*	1  0 = no request is being made,       1 = request approval for this state
	// 		2  0 = this state is work in progress, 1 = work is concluded on this StateInfo
	//      3  0 = not reverted, 1 = reverted
	//  *   4  0 = no owner change, 1 = owner change, changer will be the UID
	//         of LastModBy on this StateInfo, and creator of the StateInfo
	//         with new owner
	// we'll set the lower byte to (not approved, no request being made, work concluded )
	//--------------------------------------------------------------------------
	a := si
	a.FLAGS &= si.FLAGS & 0xefffffffffffffe0
	util.Console("before: a.FLAGS = %x\n", a.FLAGS)
	a.FLAGS |= 0x14
	a.Reason = foo.Records[0].Reason // save the reason
	util.Console("after: a.FLAGS = %x\n", a.FLAGS)
	if err = db.UpdateStateInfo(ctx, &a); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
	}

	//--------------------------------------------------------------------------
	// now create the new owner state info...
	//--------------------------------------------------------------------------
	a.SIID = 0
	a.FLAGS = 0
	a.OwnerDt = time.Now()
	a.OwnerUID = foo.Records[0].OwnerUID // this has the new owner
	a.Reason = ""
	if _, err = db.InsertStateInfo(ctx, &a); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
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

// saveStateReady sets the flag to indicate the owner wishes the
// to approve the state and advance to the next state
//
// Here we expect to read the command "requestapprove" and an array of records
// containing precisely one StateInfo struct, the latest one for
// the state in question.
//
//	@URL /v1/StateInfo/PRID
//
//-----------------------------------------------------------------------------
func saveStateReady(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	foo, si, sess, err := stateInfoHelper(w, r, d)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}

	//--------------------------------------------------------------------------
	// Only the owner can ask for approval
	//--------------------------------------------------------------------------
	util.Console("sess.UID = %d, OwnerUID = %d\n", sess.UID, foo.Records[0].OwnerUID)
	if si.OwnerUID != sess.UID {
		e := fmt.Errorf("Only the owner can request approval")
		SvcErrorReturn(w, e)
		return
	}

	//--------------------------------------------------------------------------
	// before we save it, make sure that there is an approver
	//--------------------------------------------------------------------------
	if si.ApproverUID < 1 {
		e := fmt.Errorf("An approver must be assigned first")
		SvcErrorReturn(w, e)
		return
	}

	//--------------------------------------------------------------------------
	// mark this version as ready to be approved...  The request approval flag
	// (that is bit 1).
	// we'll set the lower byte to (not approved, no request being made, work concluded )
	//--------------------------------------------------------------------------
	si.FLAGS |= 0x2
	if err = db.UpdateStateInfo(r.Context(), &si); err != nil {
		SvcErrorReturn(w, err)
		return
	}

	SvcWriteSuccessResponse(w)
}

// saveStateApprove approve the proposal to advance the property to the next
// state.
//
// Here we expect to read the command "approve" and an array of records
// containing precisely one StateInfo struct, the latest one for
// the state in question.
//
// We need to update this SIID with the reason and create a new one with
// all the same info as in the current SIID, but not yet accepted/rejected.
//
//	@URL /v1/StateInfo/PRID
//
//-----------------------------------------------------------------------------
func saveStateApprove(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveStateApprove"
	util.Console("Entered %s\n", funcname)

	_, si, sess, err := stateInfoHelper(w, r, d)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}

	if sess.UID != si.ApproverUID {
		err = fmt.Errorf("You are not the current Approver for this state")
		SvcErrorReturn(w, err)
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

	a := si
	a.FLAGS &= si.FLAGS & 0xeffffffffffffff0
	a.FLAGS |= 0x4
	a.ApproverDt = time.Now()

	if err = db.UpdateStateInfo(ctx, &a); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
	}

	// now create the new next StateInfo
	if a.FlowState < 7 {
		a.FlowState++
		a.SIID = 0
		a.FLAGS = 0
		a.OwnerDt = time.Now()
		a.ApproverDt = util.TIME0
		if _, err = db.InsertStateInfo(ctx, &a); err != nil {
			tx.Rollback()
			SvcErrorReturn(w, err)
			return
		}

		util.Console("New stateinfo created: SIID = %d\n", a.SIID)

		// Now we need to update the property's state...
		prop, err := db.GetProperty(ctx, a.PRID)
		if err != nil {
			tx.Rollback()
			SvcErrorReturn(w, err)
			return
		}

		prop.FlowState = a.FlowState
		if err = db.UpdateProperty(ctx, &prop); err != nil {
			tx.Rollback()
			SvcErrorReturn(w, err)
			return
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

// saveStateReject reject the proposal to advance the property to the next
// and save the reason for the rejection.
//
// Here we expect to read the command "reject" and an array of records
// containing precisely one StateInfo struct, the latest one for
// the state in question.  It should contain the Reason why it was rejected.
//
// We need to update this SIID with the reason and create a new one with
// all the same info as in the current SIID, but not yet accepted/rejected.
//
//	@URL /v1/StateInfo/PRID
//
//-----------------------------------------------------------------------------
func saveStateReject(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveStateReject"

	foo, si, sess, err := stateInfoHelper(w, r, d)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}

	if sess.UID != si.ApproverUID {
		err = fmt.Errorf("You are not the current Approver for this state")
		SvcErrorReturn(w, err)
		return
	}

	//--------------------------------------------------------------------------
	// before we save it, make sure that there is a reason...
	//--------------------------------------------------------------------------
	if len(foo.Records[0].Reason) < 1 {
		e := fmt.Errorf("%s: You must supply the reason for the reject", funcname)
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
	//--------------------------------------------------------------------------
	// mark this version as rejected...  The reject is to set flags bit 0 to 1
	// 		0  valid only when ApproverUID > 0, 0 = State Approved, 1 = not approved
	// 		1  0 = no request is being made,       1 = request approval for this state
	// 		2  0 = this state is work in progress, 1 = work is concluded on this StateInfo
	// we'll set the lower byte to (not approved, no request being made, work concluded )
	//--------------------------------------------------------------------------
	reject := si
	reject.FLAGS &= si.FLAGS & 0xeffffffffffffff0
	reject.FLAGS |= 0x5
	reject.Reason = foo.Records[0].Reason
	reject.ApproverDt = time.Now()
	if err = db.UpdateStateInfo(ctx, &reject); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
	}

	//--------------------------------------------------------------------------
	// create the new "current" version.
	//--------------------------------------------------------------------------
	reject.SIID = 0
	reject.FLAGS = 0
	reject.OwnerDt = time.Now()
	reject.ApproverDt = util.TIME0
	reject.Reason = ""

	if _, err = db.InsertStateInfo(ctx, &reject); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
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

// saveStateRevert moves the property to its previous state
//
// Here we expect to read the command "revert" and an array of records
// containing precisely one StateInfo struct, the latest one for
// the state in question.  It should contain the Reason why it
// it is being reverted.
//
// We need to update this SIID with the reason and create a new one with
// all the same info as in the current SIID, but not yet accepted/rejected.
//
//	@URL /v1/StateInfo/PRID
//
//-----------------------------------------------------------------------------
func saveStateRevert(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveStateRevert"

	foo, si, sess, err := stateInfoHelper(w, r, d)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}

	if sess.UID != si.ApproverUID {
		err = fmt.Errorf("You are not the current Approver for this state")
		SvcErrorReturn(w, err)
		return
	}

	if si.FlowState == 1 {
		err = fmt.Errorf("You cannot revert from this state")
		SvcErrorReturn(w, err)
		return
	}

	//--------------------------------------------------------------------------
	// before we save it, make sure that there is a reason...
	//--------------------------------------------------------------------------
	if len(foo.Records[0].Reason) < 1 {
		e := fmt.Errorf("%s: You must supply the reason", funcname)
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

	//--------------------------------------------------------------------------
	// mark this version as reverted...  The revert sets flags bits 0:3 to 0xC
	//    bit  description
	//    ---  ---------------------------------------------------------------------------
	// 		2  0 = this state is work in progress, 1 = work is concluded on this StateInfo
	//      3  0 = not reverted                    1 = reverted
	// we'll set the lower byte to (work concluded, reverted)
	//--------------------------------------------------------------------------
	revert := si
	revert.FLAGS &= si.FLAGS & 0xeffffffffffffff0
	revert.FLAGS |= 0xc
	revert.Reason = foo.Records[0].Reason
	revert.ApproverDt = time.Now()
	if err = db.UpdateStateInfo(ctx, &revert); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
	}

	//--------------------------------------------------------------------------
	// create the new "current" version.
	//--------------------------------------------------------------------------
	revert.SIID = 0
	revert.FLAGS = 0
	revert.OwnerDt = time.Now()
	revert.ApproverDt = util.TIME0
	revert.Reason = ""
	revert.FlowState = si.FlowState - 1

	if _, err = db.InsertStateInfo(ctx, &revert); err != nil {
		tx.Rollback()
		SvcErrorReturn(w, err)
		return
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

// saveStateInfo expects as input a full state definition (an array of StateInfo
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
	// util.Console("record data = %s\n", d.data)

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
					if (foo.Records[i].OwnerUID != a[j].OwnerUID) ||
						(foo.Records[i].ApproverUID != a[j].ApproverUID) ||
						(foo.Records[i].FLAGS != a[j].FLAGS) ||
						util.EqualDtToJSONDateTime(&a[j].OwnerDt, &foo.Records[i].OwnerDt) ||
						util.EqualDtToJSONDateTime(&a[j].ApproverDt, &foo.Records[i].ApproverDt) { // if anything relevant changed
						//----------------------------------------
						// update a[j] to db...
						//----------------------------------------
						util.Console("MOD: foo.Records[i].SIID = %d\n", foo.Records[i].SIID)
						a[j].OwnerUID = foo.Records[i].OwnerUID
						a[j].ApproverUID = foo.Records[i].ApproverUID
						a[j].OwnerDt = time.Time(foo.Records[i].OwnerDt)
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
	var mm = map[int64]int{}
	var mmm = map[int64]UserInfo{}

	// util.Console("Getting all state info items for property: %d\n", d.ID)
	a, err := db.GetAllStateInfoItems(r.Context(), d.ID)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	// util.Console("Number of state info items found: %d\n", len(a))
	// util.Console("a = %#v\n", a)

	for i := 0; i < len(a); i++ {
		var gg StateInfo
		util.MigrateStructVals(&a[i], &gg)
		gg.Recid = gg.SIID
		g.Records = append(g.Records, gg)
		// util.Console("after migrate: gg.OwnerUID = %d, gg.ApproverUID = %d, gg.CreateBy = %d, gg.LastModBy = %d\n", gg.OwnerUID, gg.ApproverUID, gg.CreateBy, gg.LastModBy)

		// Keep track of the users we need, we'll pull them down after the
		// loop completes...
		//-------------------------------------------------------------------
		if gg.OwnerUID > 0 {
			mm[gg.OwnerUID] = 1
			// util.Console("+ mm[%d]\n", mm[gg.OwnerUID])
		}
		if gg.ApproverUID > 0 {
			mm[gg.ApproverUID] = 1
			// util.Console("+ mm[%d]\n", mm[gg.ApproverUID])
		}
		if gg.CreateBy > 0 {
			mm[gg.CreateBy] = 1
			// util.Console("+ mm[%d]\n", mm[gg.CreateBy])
		}
		if gg.LastModBy > 0 {
			mm[gg.LastModBy] = 1
			// util.Console("+ mm[%d]\n", mm[gg.LastModBy])
		}
	}

	// util.Console("LOOPING THROUGH mm\n")
	// for k := range mm {
	// 	var p UserInfo
	// 	// util.Console("getting user info for: k = %d, v = %d\n", k, v)
	// 	if p, err = GetUserInfo(k); err != nil {
	// 		SvcErrorReturn(w, err)
	// 		return
	// 	}
	// 	mmm[k] = p
	// }

	var uids []int64
	var p []UserInfo
	for k := range mm {
		uids = append(uids, k)
	}
	if p, err = GetUserListInfo(uids); err != nil {
		SvcErrorReturn(w, err)
		return
	}
	for i := 0; i < len(p); i++ {
		mmm[p[i].UID] = p[i]
	}

	for i := 0; i < len(g.Records); i++ {
		j := g.Records[i].ApproverUID
		var ui UserInfo
		var ok bool
		if ui, ok = mmm[j]; !ok {
			ui.FirstName = fn
			ui.LastName = ln
			mmm[j] = ui
			// util.Console("NO USER INFO FOR UID = %d\n", j)
		}
		g.Records[i].ApproverName = fmt.Sprintf("%s %s", ui.FirstName, ui.LastName)
		j = g.Records[i].OwnerUID
		if ui, ok = mmm[j]; !ok {
			ui.FirstName = fn
			ui.LastName = ln
			mmm[j] = ui
			// util.Console("NO USER INFO FOR UID = %d\n", j)
		}
		g.Records[i].OwnerName = fmt.Sprintf("%s %s", ui.FirstName, ui.LastName)
		if ui, ok = mmm[g.Records[i].CreateBy]; !ok {
			ui.FirstName = fn
			ui.LastName = ln
			mmm[j] = ui
			// util.Console("NO USER INFO FOR UID = %d\n", j)
		}
		g.Records[i].CreateByName = fmt.Sprintf("%s %s", ui.FirstName, ui.LastName)
		if ui, ok = mmm[g.Records[i].LastModBy]; !ok {
			ui.FirstName = fn
			ui.LastName = ln
			mmm[j] = ui
			// util.Console("NO USER INFO FOR UID = %d\n", j)
		}
		g.Records[i].LastModByName = fmt.Sprintf("%s %s", ui.FirstName, ui.LastName)
	}

	g.Status = "success"
	SvcWriteResponse(&g, w)
}
