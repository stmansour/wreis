package ws

import (
	"fmt"
	"net/http"
	db "wreis/db/lib"
	util "wreis/util/lib"
)

// RentStepsGrid contains the data from RentSteps that is targeted to the UI Grid that displays
// a list of RentSteps structs

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// RentStep is an individual member of the rent step list.
type RentStep struct {
	Recid int64         `json:"recid"`
	RSID  int64         // unique id for this record
	RSLID int64         // id of RentStepList to which this record belongs
	Dt    util.JSONDate // date for the rent amount; valid when RSLID.FLAGS bit 0 = 1
	Opt   string        // option string, "First year", or "years 1-3"
	Rent  float64       // amount of rent on the associated date
	FLAGS uint64        // 1<<0 :  0 -> Opt is valid, 1 -> Dt is valid
}

// SearchRentStepsResponse is the response data for a Rental Agreement Search
type SearchRentStepsResponse struct {
	Status  string     `json:"status"`
	Total   int64      `json:"total"`
	Records []RentStep `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveRentSteps is sent to save one of open time slots as a reservation
type SaveRentSteps struct {
	Cmd    string   `json:"cmd"`
	Record RentStep `json:"record"`
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetRentSteps is the struct returned on a request for a reservation.
type GetRentSteps struct {
	Status  string     `json:"status"`
	Records []RentStep `json:"records"`
}

//-----------------------------------------------------------------------------
//##########################################################################################################################################################
//-----------------------------------------------------------------------------

// SvcHandlerRentSteps formats a complete data record for an assessment for use
// with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
// The server command can be:
//      get
//      save
//      delete
//------------------------------------------------------------------------------
func SvcHandlerRentSteps(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerRentSteps, d.ID = %d\n", d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchRentSteps(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcErrorReturn(w, fmt.Errorf("RentStepsID is required but was not specified"))
				return
			}
			getRentSteps(w, r, d)
		}
		break
	case "save":
		saveRentSteps(w, r, d)
		break
	case "delete":
		deleteRentSteps(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// SvcSearchRentSteps generates a report of all RentSteps defined business d.BID
// wsdoc {
//  @Title  Search RentSteps
//	@URL /v1/RentSteps/[:GID]
//  @Method  POST
//	@Synopsis Search RentSteps
//  @Descr  Search all RentSteps and return those that match the Search Logic.
//  @Descr  The search criteria includes start and stop dates of interest.
//	@Input WebGridSearchRequest
//  @Response RentStepsSearchResponse
// wsdoc }
//-----------------------------------------------------------------------------
func SvcSearchRentSteps(w http.ResponseWriter, r *http.Request, d *ServiceData) {
}

// deleteRentSteps deletes a payment type from the database
// wsdoc {
//  @Title  Delete RentSteps
//	@URL /v1/RentSteps/PID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a RentSteps.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func deleteRentSteps(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteRentSteps"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	SvcWriteSuccessResponse(w)
}

// SaveRentSteps returns the requested assessment
// wsdoc {
//  @Title  Save RentSteps
//	@URL /v1/RentStepse/PID
//  @Method  GET
//	@Synopsis Update the information on a RentSteps with the supplied data, create if necessary.
//  @Description  This service creates a RentSteps if PID == 0 or updates a RentSteps if PID > 0 with
//  @Description  the information supplied. All fields must be supplied.
//	@Input RentStepsGridSave
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func saveRentSteps(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveRentSteps"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)

	// var foo RentStepsGridSave
	// data := []byte(d.data)
	// err := json.Unmarshal(data, &foo)
	//
	// if err != nil {
	// 	e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
	// 	SvcErrorReturn(w, e)
	// 	return
	// }

	SvcWriteSuccessResponse(w)
}

// RentStepsUpdate updates the supplied RentSteps in the database with the supplied
// info. It only allows certain fields to be updated.
//-----------------------------------------------------------------------------
func RentStepsUpdate(p *RentStep, d *ServiceData) error {
	util.Console("entered RentStepsUpdate\n")
	return nil
}

// GetRentSteps returns the list of RentStep items associated with the supplied
// RSLID.
//
// wsdoc {
//  @Title  Get RentSteps
//	@URL /v1/rentsteps/RSLID
//  @Method  GET
//	@Synopsis Get the list of RentSteps
//  @Description  Return all RentStep items for RSLID
//	@Input WebGridSearchRequest
//  @Response GetRentSteps
// wsdoc }
//-----------------------------------------------------------------------------
func getRentSteps(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getRentSteps"
	util.Console("entered %s\n", funcname)
	var g GetRentSteps
	util.Console("%s: A\n", funcname)
	a, err := db.GetRentStepsItems(r.Context(), d.ID)
	if err != nil {
		util.Console("%s: B\n", funcname)
		SvcErrorReturn(w, err)
		return
	}
	util.Console("%s: C.  num items = %d\n", funcname, len(a))
	for i := 0; i < len(a); i++ {
		var gg RentStep
		util.Console("%s: C.1  a[i] = %#v\n", funcname, a[i])
		util.MigrateStructVals(&a[i], &gg)
		util.Console("%s: C.2  gg = %#v\n", funcname, gg)
		gg.Recid = gg.RSID
		g.Records = append(g.Records, gg)
		util.Console("%s: C.3  len(g.Records) = %d\n", funcname, len(g.Records))
	}
	util.Console("%s: D\n", funcname)
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
