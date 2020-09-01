package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	db "wreis/db/lib"
	util "wreis/util/lib"
)

// TrafficGrid contains the data from Traffic that is targeted to the UI Grid

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// Traffic is an individual member of the Count list for a property.
type Traffic struct {
	Recid       int64  `json:"recid"`
	TID         int64  // unique id for this record
	PRID        int64  // unique id for this record
	Description string // Descriptionion string, "First year", or "years 1-3"
	Count       int64  // amount of Count on the associated date
	FLAGS       uint64 // 1<<0 :  0 -> Description is valid, 1 -> Dt is valid
}

// SearchTrafficResponse is the response data for a Countal Agreement Search
type SearchTrafficResponse struct {
	Status  string    `json:"status"`
	Total   int64     `json:"total"`
	Records []Traffic `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveTraffic saves all the traffic records
type SaveTraffic struct {
	Cmd     string    `json:"cmd"`
	Records []Traffic `json:"records"`
}

// UpdateTraffic updates an individual Traffic record
type UpdateTraffic struct {
	Cmd    string  `json:"cmd"`
	Record Traffic `json:"record"`
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetTraffic is the struct returned on a request for a reservation.
type GetTraffic struct {
	Status  string    `json:"status"`
	Records []Traffic `json:"records"`
}

//-----------------------------------------------------------------------------
//##########################################################################################################################################################
//-----------------------------------------------------------------------------

// SvcHandlerTraffic formats a complete data record for an assessment for use
// with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
//    /v1/Traffic/PRID
//
// The server command can be:
//      get
//      save
//      delete
//------------------------------------------------------------------------------
func SvcHandlerTraffic(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerTraffic, d.ID = %d\n", d.ID)

	if d.Service == "trafficitems" {
		// util.Console("trafficitems:  cmd = %s\n", d.wsSearchReq.Cmd)
		switch d.wsSearchReq.Cmd {
		case "get":
			getTrafficItems(w, r, d)
		case "save":
			saveTrafficItems(w, r, d)
		}
		return
	}

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchTraffic(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcErrorReturn(w, fmt.Errorf("TrafficID is required but was not specified"))
				return
			}
			// getTraffic(w, r, d)
		}
	case "save":
		// saveTraffic(w, r, d)
	case "delete":
		deleteTraffic(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
	}
}

// SvcSearchTraffic generates a report of all Traffic defined business d.BID
// wsdoc {
//  @Title  Search Traffic
//	@URL /v1/Traffic/[:GID]
//  @Method  POST
//	@Synopsis Search Traffic
//  @Descr  Search all Traffic and return those that match the Search Logic.
//  @Descr  The search criteria includes start and stop dates of interest.
//	@Input WebGridSearchRequest
//  @Response TrafficSearchResponse
// wsdoc }
//-----------------------------------------------------------------------------
func SvcSearchTraffic(w http.ResponseWriter, r *http.Request, d *ServiceData) {
}

// deleteTraffic deletes a payment type from the database
// wsdoc {
//  @Title  Delete Traffic
//	@URL /v1/Traffic/TID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a Traffic.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func deleteTraffic(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteTraffic"
	util.Console("Entered %s\n", funcname)

	var err error
	if err = db.DeleteTraffic(r.Context(), d.ID); err != nil {
		SvcErrorReturn(w, err)
		return
	}

	SvcWriteSuccessResponse(w)
}

// saveTrafficItems returns the requested assessment
// wsdoc {
//  @Title  Save Traffic
//	@URL /v1/Traffice/PRID
//  @Method  GET
//	@Synopsis Update the information on a Traffic with the supplied data, create if necessary.
//  @Description  Create or update a Traffic List
//  @Description  the information supplied. All fields must be supplied.
//	@Input TrafficGridSave
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func saveTrafficItems(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveTrafficItems"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)

	var foo SaveTraffic
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

	util.Console("read %d Traffic\n", len(foo.Records))

	//------------------------------------------------------------------------
	// if TID < 0 it is newly added.  For the rest, check to see if they've
	// changed and update those that have.
	//------------------------------------------------------------------------
	var a []db.Traffic

	if d.ID > 0 {
		a, err = db.GetTrafficItems(ctx, d.ID)
		if err != nil {
			// util.Console("%s: B\n", funcname)
			tx.Rollback()
			SvcErrorReturn(w, err)
			return
		}
	}

	util.Console("%s: PRID = %d  num items = %d\n", funcname, d.ID, len(a))
	for i := 0; i < len(foo.Records); i++ {
		if foo.Records[i].TID < 0 { // These are the new ones
			var x db.Traffic
			util.Console("ADD: foo.Records[i].TID = %d\n", foo.Records[i].TID)
			util.MigrateStructVals(&foo.Records[i], &x)
			util.Console("\tx = %#v\n", x)
			if _, err = db.InsertTraffic(ctx, &x); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
		} else { // These need to be checked
			for j := 0; j < len(a); j++ {
				if foo.Records[i].TID == a[j].TID {
					t := foo.Records[i].Count != a[j].Count                   // check for Count change
					t = t || foo.Records[i].FLAGS != a[j].FLAGS               // check for FLAGS change
					t = t || (foo.Records[i].Description != a[j].Description) // check for Description change
					if t {                                                    // if anything relevant changed
						util.Console("MOD: foo.Records[i].TID = %d\n", foo.Records[i].TID)
						a[j].Count = foo.Records[i].Count
						a[j].Description = foo.Records[i].Description
						a[j].FLAGS = foo.Records[i].FLAGS
						if err = db.UpdateTraffic(ctx, &a[j]); err != nil {
							tx.Rollback()
							SvcErrorReturn(w, err)
							return
						}
						util.Console("\tx = %#v\n", foo.Records[i])
					}
				}
			}
		}
	}
	//-------------------------------------------------------------------
	// Now we need to check for any Traffic that may have been removed
	//-------------------------------------------------------------------
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(foo.Records); j++ {
			if a[i].TID == foo.Records[j].TID {
				found = true
				break
			}
		}
		if !found {
			util.Console("DEL: a[i].TID = %d\n", a[i].TID)
			if err = db.DeleteTraffic(ctx, a[i].TID); err != nil {
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

// // TrafficUpdate updates the supplied Traffic in the database with the supplied
// // info. It only allows certain fields to be updated.
// //-----------------------------------------------------------------------------
// func TrafficUpdate(p *Traffic, d *ServiceData) error {
// 	funcname := "TrafficUpdate"
// 	util.Console("entered TrafficUpdate\n")
// 	var foo UpdateTraffic
// 	data := []byte(d.data)
// 	err := json.Unmarshal(data, &foo)
//
// 	if err != nil {
// 		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
// 		SvcErrorReturn(w, e)
// 		return
// 	}
//
// 	util.MigrateStructVals(&foo.Record, &x)
// 	// util.Console("\tx = %#v\n", x)
// 	if _, err = db.UpdaTraffic(ctx, &x); err != nil {
// 		tx.Rollback()
// 		SvcErrorReturn(w, err)
// 		return
// 	}
//
// 	return nil
// }

// GetTraffic returns the list of Traffic items associated with the supplied
// PRID.
//
// wsdoc {
//  @Title  Get Traffic Items
//	@URL /v1/trafficitems/PRID
//  @Method  GET
//	@Synopsis Get the list of Traffic
//  @Description  Return all Traffic items for PRID
//	@Input WebGridSearchRequest
//  @Response GetTraffic
// wsdoc }
//-----------------------------------------------------------------------------
func getTrafficItems(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getTrafficItems"
	util.Console("entered %s\n", funcname)
	var g GetTraffic
	// util.Console("%s: A\n", funcname)
	a, err := db.GetTrafficItems(r.Context(), d.ID)
	if err != nil {
		// util.Console("%s: B\n", funcname)
		SvcErrorReturn(w, err)
		return
	}
	// util.Console("%s: C.  num items = %d\n", funcname, len(a))
	for i := 0; i < len(a); i++ {
		var gg Traffic
		// util.Console("%s: C.1  a[i] = %#v\n", funcname, a[i])
		util.MigrateStructVals(&a[i], &gg)
		// util.Console("%s: C.2  gg = %#v\n", funcname, gg)
		gg.Recid = gg.TID
		g.Records = append(g.Records, gg)
		// util.Console("%s: C.3  len(g.Records) = %d\n", funcname, len(g.Records))
	}
	// util.Console("%s: D\n", funcname)
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
