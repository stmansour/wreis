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
	RSID  int64         // unique id for this record
	RSLID int64         // id of RentStepList to which this record belongs
	Count int64         // FIX THIS,  not sure what this is
	Dt    util.JSONDate // date for the rent amount; valid when RSLID.FLAGS bit 0 = 1
	Opt   int64         // option number, 1 .. n
	Rent  float64       // amount of rent on the associated date
	FLAGS uint64        // 1<<0 :  0 -> count is valid, 1 -> Dt is valid
}

// RentStepsGrid is the structure of data for a RentSteps we send to the UI
type RentStepsGrid struct {
	Recid          int64 `json:"recid"`
	PRID           int64
	RSLID          int64 // unique id
	RS             []RentStep
	CreateTime     util.JSONDateTime
	CreatedBy      int64
	LastModifyTime util.JSONDateTime
	LastModifyBy   int64
	//
	// RO db.RenewOptions // contains the list of RenewOptions and context
	// RS db.RentSteps    // contains the list of RentSteps and context
}

// SearchRentStepsResponse is the response data for a Rental Agreement Search
type SearchRentStepsResponse struct {
	Status  string          `json:"status"`
	Total   int64           `json:"total"`
	Records []RentStepsGrid `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveRentSteps is sent to save one of open time slots as a reservation
type SaveRentSteps struct {
	Cmd    string        `json:"cmd"`
	Record RentStepsGrid `json:"record"`
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetRentSteps is the struct returned on a request for a reservation.
type GetRentSteps struct {
	Status string        `json:"status"`
	Record RentStepsGrid `json:"record"`
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
	// funcname := "SvcSearchRentSteps"
	// util.Console("Entered %s\n", funcname)
	//
	// var g SearchRentStepsResponse
	// var err error
	//
	// //---------------------------------------------
	// // We'll grab all fields for the properties
	// //---------------------------------------------
	// q := fmt.Sprintf("SELECT %s FROM RentSteps ", db.Wdb.DBFields["RentSteps"]) // the fields we want
	//
	// // any WHERE clause work store in qw
	// qw := "" // for now, no WHERE clause
	// q += " ORDER BY "
	// order := "PRID ASC" // default ORDER
	// if len(d.wsSearchReq.Sort) > 0 {
	// 	for i := 0; i < len(d.wsSearchReq.Sort); i++ {
	// 		if i > 0 {
	// 			q += ","
	// 		}
	// 		q += d.wsSearchReq.Sort[i].Field + " " + d.wsSearchReq.Sort[i].Direction
	// 	}
	// } else {
	// 	q += order
	// }
	// // now set up the offset and limit
	// q += fmt.Sprintf(" LIMIT %d OFFSET %d", d.wsSearchReq.Limit, d.wsSearchReq.Offset)
	// g.Total, err = db.GetRowCountRaw("RentSteps", "", qw)
	// if err != nil {
	// 	util.Console("Error from db.GetRowCountRaw: %s\n", err.Error())
	// 	SvcErrorReturn(w, err)
	// 	return
	// }
	//
	// util.Console("\nQuery = %s\n\n", q)
	// rows, err := db.Wdb.DB.Query(q)
	// if err != nil {
	// 	util.Console("Error from DB Query: %s\n", err.Error())
	// 	SvcErrorReturn(w, err)
	// 	return
	// }
	// defer rows.Close()
	//
	// i := int64(d.wsSearchReq.Offset)
	// count := 0
	// for rows.Next() {
	// 	var q RentStepsGrid
	// 	var p db.RentSteps
	// 	if err = db.ReadProperties(rows, &p); err != nil {
	// 		util.Console("%s.  Error reading Person: %s\n", funcname, err.Error())
	// 	}
	// 	util.MigrateStructVals(&p, &q)
	// 	q.Recid = p.PRID
	// 	g.Records = append(g.Records, q)
	// 	count++ // update the count only after adding the record
	// 	if count >= d.wsSearchReq.Limit {
	// 		break // if we've added the max number requested, then exit
	// 	}
	// 	i++
	// }
	//
	// util.Console("g.Total = %d\n", g.Total)
	// w.Header().Set("Content-Type", "application/json")
	// g.Status = "success"
	// SvcWriteResponse(&g, w)
	//
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
func RentStepsUpdate(p *RentStepsGrid, d *ServiceData) error {
	util.Console("entered RentStepsUpdate\n")
	return nil
}

// GetRentSteps returns the requested assessment
// wsdoc {
//  @Title  Get RentSteps
//	@URL /v1/RentSteps/:PID
//  @Method  GET
//	@Synopsis Get information on a RentSteps
//  @Description  Return all fields for assessment :PID
//	@Input WebGridSearchRequest
//  @Response GetRentSteps
// wsdoc }
//-----------------------------------------------------------------------------
func getRentSteps(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getRentSteps"
	util.Console("entered %s\n", funcname)
	var g GetRentSteps
	a, err := db.GetRentSteps(r.Context(), d.ID, false)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	if a.RSLID > 0 {
		var gg RentStepsGrid
		util.MigrateStructVals(&a, &gg)
		gg.Recid = gg.RSLID
		g.Record = gg
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
