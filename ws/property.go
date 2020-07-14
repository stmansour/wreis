package ws

import (
	"fmt"
	"net/http"
	"time"
	db "wreis/db/lib"
	util "wreis/util/lib"
)

// PropertyGrid contains the data from Property that is targeted to the UI Grid that displays
// a list of Property structs
type PropertyGrid struct {
	Recid             int64 `json:"recid"`
	PID               int64
	PRID              int64  // unique id
	Name              string // property name
	YearsInBusiness   int64
	ParentCompany     string
	URL               string
	Symbol            string
	Price             float64
	DownPayment       float64
	RentableArea      int64
	RentableAreaUnits int64
	LotSize           int64
	LotSizeUnits      int64
	CapRate           float64
	AvgCap            float64
	BuildDate         time.Time
	// FLAGS
	//     1<<0  Drive Through?  0 = no, 1 = yes
	//	   1<<1  Roof & Structure Responsibility: 0 = Tenant, 1 = Landlord
	//	   1<<2  Right Of First Refusal: 0 = no, 1 = yes
	FLAGS                uint64
	Ownership            int
	TenantTradeName      string
	LeaseGuarantor       int64
	LeaseType            int64
	DeliveryDt           time.Time
	OriginalLeaseTerm    int64
	LeaseCommencementDt  time.Time
	LeaseExpirationDt    time.Time
	TermRemainingOnLease int64
	ROLID                int64
	RSLID                int64
	Address              string
	Address2             string
	City                 string
	State                string
	PostalCode           string
	Country              string
	LLResponsibilities   string
	NOI                  float64
	HQAddress            string
	HQAddress2           string
	HQCity               string
	HQState              string
	HQPostalCode         string
	HQCountry            string
	RO                   db.RenewOptions // contains the list of RenewOptions and context
	RS                   db.RentSteps    // contains the list of RentSteps and context
}

// SvcHandlerProperty formats a complete data record for an assessment for use with the w2ui Form
// For this call, we expect the URI to contain the BID and the PID as follows:
//
// The server command can be:
//      get
//      save
//      delete
//-----------------------------------------------------------------------------------
func SvcHandlerProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerProperty\n")

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchProperty(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcGridErrorReturn(w, fmt.Errorf("PropertyID is required but was not specified"))
				return
			}
			getProperty(w, r, d)
		}
		break
	case "save":
		saveProperty(w, r, d)
		break
	case "delete":
		deleteProperty(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcGridErrorReturn(w, err)
		return
	}
}

// SvcSearchProperty generates a report of all Property defined business d.BID
// wsdoc {
//  @Title  Search Property
//	@URL /v1/Property/[:GID]
//  @Method  POST
//	@Synopsis Search Property
//  @Descr  Search all Property and return those that match the Search Logic.
//  @Descr  The search criteria includes start and stop dates of interest.
//	@Input WebGridSearchRequest
//  @Response PropertySearchResponse
// wsdoc }
//-----------------------------------------------------------------------------
func SvcSearchProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchProperty"
	util.Console("Entered %s\n", funcname)

	// var g PropertySearchResponse
	//
	// util.Console("g.Total = %d\n", g.Total)
	// w.Header().Set("Content-Type", "application/json")
	// g.Status = "success"
	// SvcWriteResponse(&g, w)

}

// deleteProperty deletes a payment type from the database
// wsdoc {
//  @Title  Delete Property
//	@URL /v1/Property/PID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a Property.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func deleteProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteProperty"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	SvcWriteSuccessResponse(w)
}

// SaveProperty returns the requested assessment
// wsdoc {
//  @Title  Save Property
//	@URL /v1/Propertye/PID
//  @Method  GET
//	@Synopsis Update the information on a Property with the supplied data, create if necessary.
//  @Description  This service creates a Property if PID == 0 or updates a Property if PID > 0 with
//  @Description  the information supplied. All fields must be supplied.
//	@Input PropertyGridSave
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func saveProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveProperty"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)

	// var foo PropertyGridSave
	// data := []byte(d.data)
	// err := json.Unmarshal(data, &foo)
	//
	// if err != nil {
	// 	e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
	// 	SvcGridErrorReturn(w, e)
	// 	return
	// }

	SvcWriteSuccessResponse(w)
}

// PropertyUpdate updates the supplied Property in the database with the supplied
// info. It only allows certain fields to be updated.
//-----------------------------------------------------------------------------
func PropertyUpdate(p *PropertyGrid, d *ServiceData) error {
	util.Console("entered PropertyUpdate\n")
	return nil
}

// GetProperty returns the requested assessment
// wsdoc {
//  @Title  Get Property
//	@URL /v1/dep/:BUI/:PID
//  @Method  GET
//	@Synopsis Get information on a Property
//  @Description  Return all fields for assessment :PID
//	@Input WebGridSearchRequest
//  @Response PropertyGetResponse
// wsdoc }
//-----------------------------------------------------------------------------
func getProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getProperty"
	util.Console("entered %s\n", funcname)
	// var g PropertyGetResponse
	// a, err := db.GetProperty(d.ID)
	// if err != nil {
	// 	SvcGridErrorReturn(w, err)
	// 	return
	// }
	// if a.PID > 0 {
	// 	var gg PropertyGrid
	// 	util.MigrateStructVals(&a, &gg)
	// 	gg.Recid = gg.PID
	// 	g.Record = gg
	// }
	// g.Status = "success"
	// SvcWriteResponse(&g, w)
}
