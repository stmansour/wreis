package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	db "wreis/db/lib"
	util "wreis/util/lib"
)

// RenewOptionsGrid contains the data from RenewOptions that is targeted to the UI Grid that displays
// a list of RenewOptions structs

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// RenewOption is an individual member of the rent step list.
type RenewOption struct {
	Recid int64         `json:"recid"`
	ROID  int64         // unique id for this record
	ROLID int64         // id of RenewOptionList to which this record belongs
	Dt    util.JSONDate // date for the rent amount; valid when ROLID.FLAGS bit 0 = 1
	Opt   string        // option string, "First year", or "years 1-3"
	Rent  float64       // amount of rent on the associated date
	FLAGS uint64        // 1<<0 :  0 -> Opt is valid, 1 -> Dt is valid
}

// SearchRenewOptionsResponse is the response data for a Rental Agreement Search
type SearchRenewOptionsResponse struct {
	Status  string        `json:"status"`
	Total   int64         `json:"total"`
	Records []RenewOption `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveRenewOptions is sent to save one of open time slots as a reservation
type SaveRenewOptions struct {
	Cmd     string        `json:"cmd"`
	PRID    int64         `json:"PRID"`
	Records []RenewOption `json:"records"`
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetRenewOptions is the struct returned on a request for a reservation.
type GetRenewOptions struct {
	Status  string        `json:"status"`
	Records []RenewOption `json:"records"`
}

//-----------------------------------------------------------------------------
//##########################################################################################################################################################
//-----------------------------------------------------------------------------

// SvcHandlerRenewOptions formats a complete data record for an assessment for use
// with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
//    /v1/renewoptions/ROLID
//
// The server command can be:
//      get
//      save
//      delete
//------------------------------------------------------------------------------
func SvcHandlerRenewOptions(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerRenewOptions, d.ID = %d\n", d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchRenewOptions(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcErrorReturn(w, fmt.Errorf("RenewOptionsID is required but was not specified"))
				return
			}
			getRenewOptions(w, r, d)
		}
	case "save":
		saveRenewOptions(w, r, d)
	case "delete":
		deleteRenewOptions(w, r, d)
	default:
		err := fmt.Errorf("unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// SvcSearchRenewOptions generates a report of all RenewOptions defined business d.BID
// wsdoc {
//  @Title  Search RenewOptions
//	@URL /v1/RenewOptions/[:GID]
//  @Method  POST
//	@Synopsis Search RenewOptions
//  @Descr  Search all RenewOptions and return those that match the Search Logic.
//  @Descr  The search criteria includes start and stop dates of interest.
//	@Input WebGridSearchRequest
//  @Response RenewOptionsSearchResponse
// wsdoc }
//-----------------------------------------------------------------------------
func SvcSearchRenewOptions(w http.ResponseWriter, r *http.Request, d *ServiceData) {
}

// deleteRenewOptions deletes a payment type from the database
// wsdoc {
//  @Title  Delete RenewOptions
//	@URL /v1/RenewOptions/PID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a RenewOptions.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func deleteRenewOptions(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteRenewOptions"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	SvcWriteSuccessResponse(w)
}

// SaveRenewOptions - adds or updates the list of renew options
//-----------------------------------------------------------------------------
func saveRenewOptions(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveRenewOptions"
	util.Console("Entered %s\n", funcname)
	// util.Console("record data = %s\n", d.data)

	var foo SaveRenewOptions
	data := []byte(d.data)
	// util.Console("data to unmarshal:  %s\n", string(data))
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

	var ReturnedROLID = int64(0)

	//------------------------------------------------------------------------
	// Make sure we have a list for this set of RenewOptions
	//------------------------------------------------------------------------
	// util.Console("len(foo.Records) = %d\n", len(foo.Records))
	if len(foo.Records) > 0 {
		// util.Console("foo.Records[0].ROLID = %d\n", foo.Records[0].ROLID)

		if foo.Records[0].ROLID < 1 {
			var ROLID int64
			var list db.RenewOptions
			// util.Console("Create new RenewOptions\n")
			if ROLID, err = db.InsertRenewOptions(ctx, &list); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
			// util.Console("New ROLID = %d\n", ROLID)
			// util.Console("foo.PRID = %d\n", foo.PRID)
			if ReturnedROLID < 1 {
				ReturnedROLID = ROLID
			}

			//---------------------------------------------------
			// now update this property to point to the list...
			//---------------------------------------------------
			var prop db.Property
			if prop, err = db.GetProperty(ctx, foo.PRID); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
			// util.Console("successfully got PRID = %d\n", prop.PRID)

			prop.ROLID = ROLID
			// util.Console("set prop.ROLID = %d\n", prop.ROLID)
			if err = db.UpdateProperty(ctx, &prop); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
			// util.Console("Updated! Will now save all RenewOption records\n")

			//---------------------------------------------------
			// now update each renew option...
			//---------------------------------------------------
			for i := 0; i < len(foo.Records); i++ {
				foo.Records[i].ROLID = ROLID
				var x db.RenewOption
				// util.Console("ADD: foo.Records[i].ROID = %d\n", foo.Records[i].ROID)
				util.MigrateStructVals(&foo.Records[i], &x)
				if _, err = db.InsertRenewOption(ctx, &x); err != nil {
					tx.Rollback()
					SvcErrorReturn(w, err)
					return
				}

			}
			// util.Console("Done.  We can commit the txn and exit at this point\n")
			//---------------------------------------
			// commit
			//---------------------------------------
			if err := tx.Commit(); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
			SvcWriteSuccessResponseWithID(w, ReturnedROLID)
			return
		}
	}

	// util.Console("read %d RenewOptions\n", len(foo.Records))

	//------------------------------------------------------------------------
	// if ROID < 0 it is newly added.  For the rest, check to see if they've
	// changed and update those that have.
	//------------------------------------------------------------------------
	var a []db.RenewOption

	if d.ID > 0 {
		a, err = db.GetRenewOptionsItems(ctx, d.ID)
		if err != nil {
			util.Console("%s: B\n", funcname)
			tx.Rollback()
			SvcErrorReturn(w, err)
			return
		}
	}

	// util.Console("%s: ROLID = %d  num items = %d\n", funcname, d.ID, len(a))
	for i := 0; i < len(foo.Records); i++ {
		if ReturnedROLID < 1 {
			ReturnedROLID = d.ID
		}
		if foo.Records[i].ROID < 1 {
			var x db.RenewOption
			util.Console("ADD: foo.Records[i].ROID = %d\n", foo.Records[i].ROID)
			util.MigrateStructVals(&foo.Records[i], &x)
			x.ROLID = d.ID
			util.Console("\tx = %#v\n", x)
			if _, err = db.InsertRenewOption(ctx, &x); err != nil {
				tx.Rollback()
				SvcErrorReturn(w, err)
				return
			}
		} else {
			for j := 0; j < len(a); j++ {
				if foo.Records[i].ROID == a[j].ROID {
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
						util.Console("MOD: foo.Records[i].ROID = %d\n", foo.Records[i].ROID)
						a[j].Rent = foo.Records[i].Rent
						a[j].Opt = foo.Records[i].Opt
						a[j].Dt = time.Time(foo.Records[i].Dt)
						a[j].FLAGS = foo.Records[i].FLAGS
						if err = db.UpdateRenewOption(ctx, &a[j]); err != nil {
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
	// Now we need to check for any RenewOption that may have been removed
	//-------------------------------------------------------------------
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(foo.Records); j++ {
			if a[i].ROID == foo.Records[j].ROID {
				found = true
				break
			}
		}
		if !found {
			util.Console("DEL: a[i].ROID = %d\n", a[i].ROID)
			if err = db.DeleteRenewOption(ctx, a[i].ROID); err != nil {
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
	SvcWriteSuccessResponseWithID(w, ReturnedROLID)
}

// RenewOptionsUpdate updates the supplied RenewOptions in the database with the supplied
// info. It only allows certain fields to be updated.
//-----------------------------------------------------------------------------
func RenewOptionsUpdate(p *RenewOption, d *ServiceData) error {
	util.Console("entered RenewOptionsUpdate\n")
	return nil
}

// GetRenewOptions returns the list of RenewOption items associated with the supplied
// ROLID.
//
// wsdoc {
//  @Title  Get RenewOptions
//	@URL /v1/renewoptions/ROLID
//  @Method  GET
//	@Synopsis Get the list of RenewOptions
//  @Description  Return all RenewOption items for ROLID
//	@Input WebGridSearchRequest
//  @Response GetRenewOptions
// wsdoc }
//-----------------------------------------------------------------------------
func getRenewOptions(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getRenewOptions"
	util.Console("entered %s\n", funcname)
	var g GetRenewOptions
	// util.Console("%s: A\n", funcname)
	a, err := db.GetRenewOptionsItems(r.Context(), d.ID)
	if err != nil {
		// util.Console("%s: B\n", funcname)
		SvcErrorReturn(w, err)
		return
	}
	// util.Console("%s: C.  num items = %d\n", funcname, len(a))
	for i := 0; i < len(a); i++ {
		var gg RenewOption
		// util.Console("%s: C.1  a[i] = %#v\n", funcname, a[i])
		util.MigrateStructVals(&a[i], &gg)
		// util.Console("%s: C.2  gg = %#v\n", funcname, gg)
		gg.Recid = gg.ROID
		g.Records = append(g.Records, gg)
		// util.Console("%s: C.3  len(g.Records) = %d\n", funcname, len(g.Records))
	}
	// util.Console("%s: D\n", funcname)
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
