package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	db "wreis/db/lib"
	"wreis/session"
	util "wreis/util/lib"
)

// PropertyGrid contains the data from Property that is targeted to the UI Grid that displays
// a list of Property structs

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// PropertyGrid is the structure of data for a property we send to the UI
type PropertyGrid struct {
	Recid             int64  `json:"recid"`
	PRID              int64  // unique id
	Name              string // property name
	YearFounded       int64
	ParentCompany     string
	URL               string
	Symbol            string
	Price             float64
	DownPayment       float64
	RentableArea      int64
	RentableAreaUnits int64
	LotSize           float64
	LotSizeUnits      int64
	CapRate           float64
	AvgCap            float64
	BuildYear         int64
	RenovationYear    int64
	FlowState         int64
	//======================================================================
	// FLAGS
	//     1<<0  Drive Through?  0 = no, 1 = yes
	//	   1<<1  Roof & Structure Responsibility: 0 = Tenant, 1 = Landlord
	//	   1<<2  Right Of First Refusal: 0 = no, 1 = yes
	//======================================================================
	FLAGS              uint64
	OwnershipType      int
	TenantTradeName    string
	LeaseGuarantor     int64
	LeaseType          int64
	OriginalLeaseTerm  int64
	RentCommencementDt util.JSONDateTime
	LeaseExpirationDt  util.JSONDateTime
	ROLID              int64
	RSLID              int64
	Address            string
	Address2           string
	City               string
	State              string
	PostalCode         string
	Country            string
	LLResponsibilities string
	NOI                float64
	HQCity             string
	HQState            string
	Img1               string
	Img2               string
	Img3               string
	Img4               string
	Img5               string
	Img6               string
	Img7               string
	Img8               string
	Img9               string
	Img10              string
	Img11              string
	Img12              string
	CreateTime         util.JSONDateTime
	CreateBy           int64
	LastModTime        util.JSONDateTime
	LastModBy          int64
	//
	// RO db.RenewOptions // contains the list of RenewOptions and context
	// RS db.RentSteps    // contains the list of RentSteps and context
}

// which fields needs to be fetched for SQL query for property grid
var propFieldsMap = map[string][]string{
	"PRID":               {"Property.PRID"},
	"Name":               {"Property.Name"},
	"YearFounded":        {"Property.YearFounded"},
	"ParentCompany":      {"Property.ParentCompany"},
	"URL":                {"Property.URL"},
	"Symbol":             {"Property.Symbol"},
	"Price":              {"Property.Price"},
	"DownPayment":        {"Property.DownPayment"},
	"RentableArea":       {"Property.RentableArea"},
	"RentableAreaUnits":  {"Property.RentableAreaUnits"},
	"LotSize":            {"Property.LotSize"},
	"LotSizeUnits":       {"Property.LotSizeUnits"},
	"CapRate":            {"Property.CapRate"},
	"AvgCap":             {"Property.AvgCap"},
	"BuildYear":          {"Property.BuildYear"},
	"RenovationYear":     {"Property.RenovationYear"},
	"FlowState":          {"Property.FlowState"},
	"FLAGS":              {"Property.FLAGS"},
	"OwnershipType":      {"Property.OwnershipType"},
	"TenantTradeName":    {"Property.TenantTradeName"},
	"LeaseGuarantor":     {"Property.LeaseGuarantor"},
	"LeaseType":          {"Property.LeaseType"},
	"OriginalLeaseTerm":  {"Property.OriginalLeaseTerm"},
	"RentCommencementDt": {"Property.RentCommencementDt"},
	"LeaseExpirationDt":  {"Property.LeaseExpirationDt"},
	"ROLID":              {"Property.ROLID"},
	"RSLID":              {"Property.RSLID"},
	"Address":            {"Property.Address"},
	"Address2":           {"Property.Address2"},
	"City":               {"Property.City"},
	"State":              {"Property.State"},
	"PostalCode":         {"Property.PostalCode"},
	"Country":            {"Property.Country"},
	"LLResponsibilities": {"Property.LLResponsibilities"},
	"NOI":                {"Property.NOI"},
	"HQCity":             {"Property.HQCity"},
	"HQState":            {"Property.HQState"},
	"Img1":               {"Property.Img1"},
	"Img2":               {"Property.Img2"},
	"Img3":               {"Property.Img3"},
	"Img4":               {"Property.Img4"},
	"Img5":               {"Property.Img5"},
	"Img6":               {"Property.Img6"},
	"Img7":               {"Property.Img7"},
	"Img8":               {"Property.Img8"},
	"Img9":               {"Property.Img9"},
	"Img10":              {"Property.Img10"},
	"Img11":              {"Property.Img11"},
	"Img12":              {"Property.Img12"},
	"CreateTime":         {"Property.CreateTime"},
	"CreateBy":           {"Property.CreateBy"},
	"LastModTime":        {"Property.LastModTime"},
	"LastModBy":          {"Property.LastModBy"},
}

// which fields needs to be fetched for SQL query for property grid
var propQuerySelectFields = []string{
	"Property.PRID",
	"Property.Name",
	"Property.YearFounded",
	"Property.ParentCompany",
	"Property.URL",
	"Property.Symbol",
	"Property.Price",
	"Property.DownPayment",
	"Property.RentableArea",
	"Property.RentableAreaUnits",
	"Property.LotSize",
	"Property.LotSizeUnits",
	"Property.CapRate",
	"Property.AvgCap",
	"Property.BuildYear",
	"Property.RenovationYear",
	"Property.FlowState",
	"Property.FLAGS",
	"Property.OwnershipType",
	"Property.TenantTradeName",
	"Property.LeaseGuarantor",
	"Property.LeaseType",
	"Property.OriginalLeaseTerm",
	"Property.RentCommencementDt",
	"Property.LeaseExpirationDt",
	"Property.ROLID",
	"Property.RSLID",
	"Property.Address",
	"Property.Address2",
	"Property.City",
	"Property.State",
	"Property.PostalCode",
	"Property.Country",
	"Property.LLResponsibilities",
	"Property.NOI",
	"Property.HQCity",
	"Property.HQState",
	"Property.Img1",
	"Property.Img2",
	"Property.Img3",
	"Property.Img4",
	"Property.Img5",
	"Property.Img6",
	"Property.Img7",
	"Property.Img8",
	"Property.Img9",
	"Property.Img10",
	"Property.Img11",
	"Property.Img12",
	"Property.CreateTime",
	"Property.CreateBy",
	"Property.LastModTime",
	"Property.LastModBy",
}

// this is the list of fields to search for a string if the field name is blank
var propDefaultFields = []string{
	"Name",
	"City",
	"State",
	"PostalCode",
}

// SearchPropertyResponse is the response data for a Rental Agreement Search
type SearchPropertyResponse struct {
	Status  string         `json:"status"`
	Total   int64          `json:"total"`
	Records []PropertyGrid `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveProperty is sent to save one of open time slots as a reservation
type SaveProperty struct {
	Cmd    string       `json:"cmd"`
	Record PropertyGrid `json:"record"`
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetProperty is the struct returned on a request for a reservation.
type GetProperty struct {
	Status string       `json:"status"`
	Record PropertyGrid `json:"record"`
}

// StateFilter captures the filter data from the propertyGrid toolbar that
// indicates how the properties should be filtered by State
//--------------------------------------------------------------------------
type StateFilter struct {
	States         []int64 `json:"statefilter"`
	ShowTerminated int64   `json:"showTerminated"`
	MyQueue        int64   `json:"myQueue"`
}

//-----------------------------------------------------------------------------
//##########################################################################################################################################################
//-----------------------------------------------------------------------------

// SvcHandlerProperty formats a complete data record for an property for use
// with the w2ui Form
//
// The server command can be:
//      get
//      save
//      delete
//------------------------------------------------------------------------------
func SvcHandlerProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerProperty, d.ID = %d\n", d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchProperty(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcErrorReturn(w, fmt.Errorf("field PropertyID is required"))
				return
			}
			getProperty(w, r, d)
		}
	case "save":
		saveProperty(w, r, d)
	case "delete":
		deleteProperty(w, r, d)
	default:
		err := fmt.Errorf("unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// SvcSearchProperty generates a report of all Property records matching the
// search criteria.
//
//	@URL /v1/Property/
//
//-----------------------------------------------------------------------------
func SvcSearchProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchProperty"
	util.Console("Entered %s\n", funcname)
	var g SearchPropertyResponse
	var err error
	var sf StateFilter
	var statefltr string
	var joinwhere string
	var joins string
	var whereClause, orderClause string

	sess, ok := session.GetSessionFromContext(r.Context())
	if !ok {
		SvcErrorReturn(w, db.ErrSessionRequired)
		return
	}

	if strings.Contains(d.data, "statefilter") {
		// util.Console("Unmarshal statefilter:  d.data = %s\n", d.data)
		err = json.Unmarshal([]byte(d.data), &sf)
		if err != nil {
			e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
			SvcErrorReturn(w, e)
			return
		}
		util.Console("sf = %#v\n", sf)
		if sf.MyQueue > 0 {
			joins = " LEFT JOIN StateInfo ON (StateInfo.PRID=Property.PRID AND StateInfo.FlowState = Property.FlowState)"
			var findterminated = 0
			if sf.ShowTerminated == 1 {
				findterminated = 64
			}
			joinwhere = fmt.Sprintf(`
				WHERE
					(StateInfo.FLAGS & 0x4)=0 AND (StateInfo.FLAGS & 64)=%d
					AND
					(
						(%d = StateInfo.OwnerUID AND (StateInfo.FLAGS & 2)=0)
				  		OR
				 		(%d = StateInfo.ApproverUID AND (StateInfo.FLAGS & 2)=2)
					)
				`, findterminated, sess.UID, sess.UID)
		} else {
			if len(sf.States) > 0 {
				statefltr = " Property.FlowState IN ("
				for j := 0; j < len(sf.States); j++ {
					statefltr += fmt.Sprintf("%d", sf.States[j])
					if j+1 < len(sf.States) {
						statefltr += ","
					}
				}
				statefltr += ")"
			}
			// util.Console("len(sf.States) = %d\n", len(sf.States))
			// The terminate bit is 1<<6  == 64.  If that bit is set then the property is terminated
			switch sf.ShowTerminated {
			case 0:
				statefltr += " AND ((Property.FLAGS & 64)=0)"
			case 1:
				statefltr += " AND ((Property.FLAGS & 64)!=0)"
			}
		}
	}

	whr := ""
	order := `Property.Name ASC` // default ORDER

	// get where clause and order clause for sql query
	// util.Console("len(d.wsSearchReq.Search) = %d\n", len(d.wsSearchReq.Search))
	HandleBlankSearchField(d, propDefaultFields)
	// util.Console("AFTER HandleBlankSearchField:  len(d.wsSearchReq.Search) = %d\n", len(d.wsSearchReq.Search))

	//------------------------------------------------------
	// use MyQueue if present, otherwise use generic...
	//------------------------------------------------------
	whereClause, orderClause = GetSearchAndSortSQL(d, propFieldsMap)
	if sf.MyQueue > 0 {
		whr = joins + " " + joinwhere
	} else {
		if len(statefltr) > 0 {
			if len(whereClause) > 0 {
				whereClause += " AND "
			}
			whereClause += statefltr
		}
		if len(whereClause) > 0 {
			whr += "WHERE " + whereClause
		}
	}
	if len(orderClause) > 0 {
		order = orderClause
	}
	util.Console("whr = %s\n", whr)
	util.Console("order = %s\n", order)

	query := `
	SELECT DISTINCT {{.SelectClause}}
	FROM Property {{.WhereClause}}
	ORDER BY {{.OrderClause}}`

	qc := db.QueryClause{
		"SelectClause": strings.Join(propQuerySelectFields, ","),
		"WhereClause":  whr,
		"OrderClause":  order,
	}

	countQuery := db.RenderSQLQuery(query, qc)

	util.Console("countQuery = %s\n", countQuery)
	g.Total, err = db.GetQueryCount(countQuery)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	// util.Console("g.Total = %d\n", g.Total)

	// FETCH the records WITH LIMIT AND OFFSET
	// limit the records to fetch from server, page by page
	limitAndOffsetClause := `
	LIMIT {{.LimitClause}}
	OFFSET {{.OffsetClause}};`

	// build query with limit and offset clause
	// if query ends with ';' then remove it
	queryWithLimit := query + limitAndOffsetClause

	// Add limit and offset value
	qc["LimitClause"] = strconv.Itoa(d.wsSearchReq.Limit)
	qc["OffsetClause"] = strconv.Itoa(d.wsSearchReq.Offset)

	// get formatted query with substitution of select, where, order clause
	qry := db.RenderSQLQuery(queryWithLimit, qc)
	util.Console("SvcSearchProperty: db query = %s\n", qry)

	// execute the query
	rows, err := db.Wdb.DB.Query(qry)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	defer rows.Close()

	i := int64(d.wsSearchReq.Offset)
	count := 0
	for rows.Next() {
		q, err := PropertyRowScan(rows)
		if err != nil {
			SvcErrorReturn(w, err)
			return
		}
		q.Recid = i

		g.Records = append(g.Records, q)
		count++ // update the count only after adding the record
		if count >= d.wsSearchReq.Limit {
			break // if we've added the max number requested, then exit
		}
		i++
	}

	err = rows.Err()
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}

	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// PropertyRowScan scans a result from sql row and dump it in a
// PropertyGrid struct
//
// RETURNS
//  Property
//-----------------------------------------------------------------------------
func PropertyRowScan(rows *sql.Rows) (PropertyGrid, error) {
	var q PropertyGrid
	err := rows.Scan(
		&q.PRID,
		&q.Name,
		&q.YearFounded,
		&q.ParentCompany,
		&q.URL,
		&q.Symbol,
		&q.Price,
		&q.DownPayment,
		&q.RentableArea,
		&q.RentableAreaUnits,
		&q.LotSize,
		&q.LotSizeUnits,
		&q.CapRate,
		&q.AvgCap,
		&q.BuildYear,
		&q.RenovationYear,
		&q.FlowState,
		&q.FLAGS,
		&q.OwnershipType,
		&q.TenantTradeName,
		&q.LeaseGuarantor,
		&q.LeaseType,
		&q.OriginalLeaseTerm,
		&q.RentCommencementDt,
		&q.LeaseExpirationDt,
		&q.ROLID,
		&q.RSLID,
		&q.Address,
		&q.Address2,
		&q.City,
		&q.State,
		&q.PostalCode,
		&q.Country,
		&q.LLResponsibilities,
		&q.NOI,
		&q.HQCity,
		&q.HQState,
		&q.Img1,
		&q.Img2,
		&q.Img3,
		&q.Img4,
		&q.Img5,
		&q.Img6,
		&q.Img7,
		&q.Img8,
		&q.Img9,
		&q.Img10,
		&q.Img11,
		&q.Img12,
		&q.CreateTime,
		&q.CreateBy,
		&q.LastModTime,
		&q.LastModBy,
	)
	return q, err
}

// deleteProperty deletes a payment type from the database
// wsdoc {
//  @Title  Delete Property
//	@URL /v1/Property/PRID
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

// SaveProperty returns the requested property
// wsdoc {
//  @Title  Save Property
//	@URL /v1/Property/PRID
//  @Method  GET
//	@Synopsis Update the information on a Property with the supplied data, create if necessary.
//  @Description  This service creates a Property if PRID == 0 or updates a Property if PRID > 0 with
//  @Description  the information supplied. All fields must be supplied.
//	@Input SaveProperty
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func saveProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveProperty"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	util.Console("PRID = %d\n", d.ID)

	var foo SaveProperty

	data := []byte(d.data)
	err := json.Unmarshal(data, &foo)

	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}

	// util.Console("read foo.  foo.Record.PRID = %d, foo.Record.Name = %s\n", foo.Record.PRID, foo.Record.Name)
	var p db.Property
	if err = util.MigrateStructVals(&foo.Record, &p); err != nil {
		e := fmt.Errorf("%s: Error with MigrateStructVals:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}
	util.Console("After Migrate:  p.LeaseGuarantor = %d\n", p.LeaseGuarantor)
	if p.PRID < 1 {
		if _, err = db.InsertProperty(r.Context(), &p); err != nil {
			e := fmt.Errorf("%s: Error with db.CreateProperty:  %s", funcname, err.Error())
			SvcErrorReturn(w, e)
			return
		}
		//----------------------------------
		// Now, save the initial state...
		//----------------------------------
		sess, ok := session.GetSessionFromContext(r.Context())
		if !ok {
			SvcErrorReturn(w, db.ErrSessionRequired)
			return
		}

		now := time.Now()
		var s = db.StateInfo{
			PRID:        p.PRID,
			OwnerUID:    sess.UID,
			ApproverUID: sess.UID,
			OwnerDt:     now,
			ApproverDt:  util.TIME0,
			FlowState:   1,
			FLAGS:       uint64(0),
		}
		if _, err := db.InsertStateInfo(r.Context(), &s); err != nil {
			SvcErrorReturn(w, err)
			return
		}
	} else {
		if err = db.UpdateProperty(r.Context(), &p); err != nil {
			e := fmt.Errorf("%s: error with db.UpdateProperty:  %s", funcname, err.Error())
			SvcErrorReturn(w, e)
			return
		}
	}
	// util.Console("UpdateProperty completed successfully\n")
	SvcWriteSuccessResponseWithID(w, p.PRID)
}

// PropertyUpdate updates the supplied Property in the database with the supplied
// info. It only allows certain fields to be updated.
//-----------------------------------------------------------------------------
func PropertyUpdate(p *PropertyGrid, d *ServiceData) error {
	util.Console("entered PropertyUpdate\n")
	return nil
}

// GetProperty returns the requested property
// wsdoc {
//  @Title  Get Property
//	@URL /v1/property/:PRID
//  @Method  GET
//	@Synopsis Get information on a Property
//  @Description  Return all fields for property :PRID
//	@Input WebGridSearchRequest
//  @Response GetProperty
// wsdoc }
//-----------------------------------------------------------------------------
func getProperty(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getProperty"
	util.Console("entered %s\n", funcname)
	var g GetProperty
	a, err := db.GetProperty(r.Context(), d.ID)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	if a.PRID == d.ID {
		var gg PropertyGrid
		util.MigrateStructVals(&a, &gg)
		gg.Recid = gg.PRID
		g.Record = gg
	} else {
		err = fmt.Errorf("could not find property with PRID = %d", d.ID)
		SvcErrorReturn(w, err)
		return
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
