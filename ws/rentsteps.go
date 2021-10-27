package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	Cmd     string     `json:"cmd"`
	PRID    int64      `json:"PRID"`
	Records []RentStep `json:"records"`
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
//    /v1/rentsteps/RSLID
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
//	@URL /v1/RentSteps/RSLID
//  @Method  GET
//	@Synopsis Update the information on a RentSteps with the supplied data, create if necessary.
//  @Description  Create or update a RentStep List
//  @Description  the information supplied. All fields must be supplied.
//	@Input RentStepsGridSave
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func saveRentSteps(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveRentSteps"
	util.Console("Entered %s\n", funcname)
	// util.Console("record data = %s\n", d.data)

	var ReturnRSLID = int64(0)
	var foo SaveRentSteps
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

	//------------------------------------------------------------------------
	// Make sure we have a list for this set of RentSteps
	//------------------------------------------------------------------------
	// util.Console("foo.Records contains %d elements\n", len(foo.Records))
	if len(foo.Records) > 0 {

		if foo.Records[0].RSLID < 1 {
			// util.Console("New RSLID being created\n")
			var RSLID int64
			var list db.RentSteps
			if RSLID, err = db.InsertRentSteps(ctx, &list); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
			if ReturnRSLID < 1 {
				ReturnRSLID = RSLID // save the return id
			}

			// util.Console("Created.  Now updating property with RSLID = %d\n", RSLID)
			//---------------------------------------------------
			// now update this property to point to the list...
			//---------------------------------------------------
			var prop db.Property
			if prop, err = db.GetProperty(ctx, foo.PRID); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}

			prop.RSLID = RSLID
			if err = db.UpdateProperty(ctx, &prop); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
			// util.Console("Property successfully updated, prop.RSLID = %d.  Now saving RentSteps: %d\n", prop.RSLID, len(foo.Records))

			//---------------------------------------------------
			// now update each renew option...
			//---------------------------------------------------
			for i := 0; i < len(foo.Records); i++ {
				foo.Records[i].RSLID = RSLID
				var x db.RentStep
				util.MigrateStructVals(&foo.Records[i], &x)
				// util.Console("foo.Records[%d] = %#v\n", i, foo.Records[i])
				// util.Console("x = %#v\n", x)
				if _, err = db.InsertRentStep(ctx, &x); err != nil {
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
			SvcWriteSuccessResponseWithID(w, ReturnRSLID)
			return
		}
	}

	//------------------------------------------------------------------------
	// if RSID < 0 it is newly added.  For the rest, check to see if they've
	// changed and update those that have.
	//------------------------------------------------------------------------
	var a []db.RentStep

	if d.ID > 0 {
		a, err = db.GetRentStepsItems(ctx, d.ID)
		if err != nil {
			// util.Console("%s: B\n", funcname)
			tx.Rollback()
			SvcErrorReturn(w, err)
			return
		}
	}

	// util.Console("%s: RSLID = %d, len(a) = %d, len(foo) = %d\n", funcname, d.ID, len(a), len(foo.Records))
	for i := 0; i < len(foo.Records); i++ {
		if ReturnRSLID < 1 {
			ReturnRSLID = d.ID // save the return id
		}

		// util.Console("i = %d, foo.Records[i].RSID = %d\n", i, foo.Records[i].RSID)
		if foo.Records[i].RSID < 1 {
			// util.Console("IN PATH A")
			var x db.RentStep
			// util.Console("ADD: foo.Records[i].RSID = %d\n", foo.Records[i].RSID)
			util.MigrateStructVals(&foo.Records[i], &x)
			x.RSLID = d.ID
			// util.Console("\tx = %#v\n", x)
			if _, err = db.InsertRentStep(ctx, &x); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
		} else {
			// util.Console("IN PATH B")
			for j := 0; j < len(a); j++ {
				if foo.Records[i].RSID == a[j].RSID {
					t := foo.Records[i].Rent != a[j].Rent       // check for rent change
					t = t || foo.Records[i].FLAGS != a[j].FLAGS // check for FLAGS change
					switch foo.Records[i].FLAGS & 0x1 {
					case 0: // Opt
						t = t || (foo.Records[i].Opt != a[j].Opt) // check for Opt change
					case 1: // Date
						d := time.Time(foo.Records[i].Dt) // check for Date change
						t = t || !d.Equal(a[j].Dt)
					}
					if t { // if anything relevant changed
						// util.Console("MOD: foo.Records[i].RSID = %d\n", foo.Records[i].RSID)
						a[j].Rent = foo.Records[i].Rent
						a[j].Opt = foo.Records[i].Opt
						a[j].Dt = time.Time(foo.Records[i].Dt)
						a[j].FLAGS = foo.Records[i].FLAGS
						if err = db.UpdateRentStep(ctx, &a[j]); err != nil {
							tx.Rollback()
							SvcErrorReturn(w, err)
							return
						}
						// util.Console("\tx = %#v\n", foo.Records[i])
					}
				}
			}
		}
	}
	//-------------------------------------------------------------------
	// Now we need to check for any RentStep that may have been removed
	//-------------------------------------------------------------------
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(foo.Records); j++ {
			if a[i].RSID == foo.Records[j].RSID {
				found = true
				break
			}
		}
		if !found {
			// util.Console("DEL: a[i].RSID = %d\n", a[i].RSID)
			if err = db.DeleteRentStep(ctx, a[i].RSID); err != nil {
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
	SvcWriteSuccessResponseWithID(w, ReturnRSLID)
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
	// util.Console("%s: A\n", funcname)
	a, err := db.GetRentStepsItems(r.Context(), d.ID)
	if err != nil {
		// util.Console("%s: B\n", funcname)
		SvcErrorReturn(w, err)
		return
	}
	// util.Console("%s: C.  num items = %d\n", funcname, len(a))
	for i := 0; i < len(a); i++ {
		var gg RentStep
		// util.Console("%s: C.1  a[i] = %#v\n", funcname, a[i])
		util.MigrateStructVals(&a[i], &gg)
		// util.Console("%s: C.2  gg = %#v\n", funcname, gg)
		gg.Recid = gg.RSID
		g.Records = append(g.Records, gg)
		// util.Console("%s: C.3  len(g.Records) = %d\n", funcname, len(g.Records))
	}
	// util.Console("%s: D\n", funcname)
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
