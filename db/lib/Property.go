package db

import (
	"context"
	"database/sql"
	"time"
	"wreis/session"
)

// Property defines a date and a rent amount for a property. A Property record
// is part of a group or list. The group is defined by the RSLID
//-----------------------------------------------------------------------------
type Property struct {
	PRID              int64 // unique id
	ROLID             int64
	RSLID             int64
	FlowState         int64
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
	BuildDate         time.Time
	// FLAGS
	//     1<<0  Drive Through?  0 = no, 1 = yes
	//	   1<<1  Roof & Structure Responsibility: 0 = Tenant, 1 = Landlord
	//	   1<<2  Right Of First Refusal: 0 = no, 1 = yes
	//     1<<3  0 = processing/processed. 1 = Terminated
	FLAGS                     uint64
	Ownership                 int
	TenantTradeName           string
	LeaseGuarantor            int64
	LeaseType                 int64
	OriginalLeaseTerm         int64
	RentCommencementDt        time.Time
	LeaseExpirationDt         time.Time
	TermRemainingOnLease      int64
	TermRemainingOnLeaseUnits int64
	Address                   string
	Address2                  string
	City                      string
	State                     string
	PostalCode                string
	Country                   string
	LLResponsibilities        string
	NOI                       float64
	HQCity                    string
	HQState                   string
	Img1                      string // url to image
	Img2                      string
	Img3                      string
	Img4                      string
	Img5                      string
	Img6                      string
	Img7                      string
	Img8                      string
	LastModTime               time.Time    // when was the record last written
	LastModBy                 int64        // id of user that did the modify
	CreateTime                time.Time    // when was this record created
	CreateBy                  int64        // id of user that created it
	RO                        RenewOptions // contains the list of RenewOptions and context
	RS                        RentSteps    // contains the list of RentSteps and context
}

// DeleteProperty deletes the Property with the specified id from the database
//
// INPUTS
// ctx - db context
// id - PRID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteProperty(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Property", Wdb.Prepstmt.DeleteProperty, id)
}

// GetProperty reads and returns a Property structure
//
// INPUTS
// ctx - db context
// id - PRID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetProperty(ctx context.Context, id int64) (Property, error) {
	var a Property
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Wdb.Prepstmt.GetProperty, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadProperty(row, &a)
}

// InsertProperty writes a new Property record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertProperty(ctx context.Context, a *Property) (int64, error) {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.PRID, ErrSessionRequired
	}
	fields := []interface{}{
		a.ROLID,
		a.RSLID,
		a.FlowState,
		a.Name,
		a.YearFounded,
		a.ParentCompany,
		a.URL,
		a.Symbol,
		a.Price,
		a.DownPayment,
		a.RentableArea,
		a.RentableAreaUnits,
		a.LotSize,
		a.LotSizeUnits,
		a.CapRate,
		a.AvgCap,
		a.BuildDate,
		a.FLAGS,
		a.Ownership,
		a.TenantTradeName,
		a.LeaseGuarantor,
		a.LeaseType,
		a.OriginalLeaseTerm,
		a.RentCommencementDt,
		a.LeaseExpirationDt,
		a.TermRemainingOnLease,
		a.TermRemainingOnLeaseUnits,
		a.Address,
		a.Address2,
		a.City,
		a.State,
		a.PostalCode,
		a.Country,
		a.LLResponsibilities,
		a.NOI,
		a.HQCity,
		a.HQState,
		a.Img1,
		a.Img2,
		a.Img3,
		a.Img4,
		a.Img5,
		a.Img6,
		a.Img7,
		a.Img8,
		sess.UID,
		sess.UID,
	}

	var err error
	a.CreateBy, a.LastModBy, a.PRID, err = genericInsert(ctx, "Property", Wdb.Prepstmt.InsertProperty, fields, a)
	return a.PRID, err
}

// InsertPropertyWithLists writes a new Property record to the database and
//     creates the RentSteps and RenewOptions if needed
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertPropertyWithLists(ctx context.Context, a *Property) (int64, error) {
	var id int64
	var err error
	if len(a.RO.RO) > 0 {
		if a.ROLID, err = InsertRenewOptionsWithList(ctx, &a.RO); err != nil {
			return id, err
		}
	}
	if len(a.RS.RS) > 0 {
		if a.RSLID, err = InsertRentStepsWithList(ctx, &a.RS); err != nil {
			return id, err
		}
	}
	return InsertProperty(ctx, a)
}

// ReadProperty reads a full Property structure of data from the database based
// on the supplied Rows pointer.
//
// INPUTS
// row - db Row pointer
// a   - pointer to struct to fill
//
// RETURNS
//
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func ReadProperty(row *sql.Row, a *Property) error {
	err := row.Scan(
		&a.PRID,
		&a.ROLID,
		&a.RSLID,
		&a.FlowState,
		&a.Name,
		&a.YearFounded,
		&a.ParentCompany,
		&a.URL,
		&a.Symbol,
		&a.Price,
		&a.DownPayment,
		&a.RentableArea,
		&a.RentableAreaUnits,
		&a.LotSize,
		&a.LotSizeUnits,
		&a.CapRate,
		&a.AvgCap,
		&a.BuildDate,
		&a.FLAGS,
		&a.Ownership,
		&a.TenantTradeName,
		&a.LeaseGuarantor,
		&a.LeaseType,
		&a.OriginalLeaseTerm,
		&a.RentCommencementDt,
		&a.LeaseExpirationDt,
		&a.TermRemainingOnLease,
		&a.TermRemainingOnLeaseUnits,
		&a.Address,
		&a.Address2,
		&a.City,
		&a.State,
		&a.PostalCode,
		&a.Country,
		&a.LLResponsibilities,
		&a.NOI,
		&a.HQCity,
		&a.HQState,
		&a.Img1,
		&a.Img2,
		&a.Img3,
		&a.Img4,
		&a.Img5,
		&a.Img6,
		&a.Img7,
		&a.Img8,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadProperties reads a full Property structure of data from the database based
// on the supplied Rows pointer.
//
// INPUTS
// row - db Row pointer
// a   - pointer to struct to fill
//
// RETURNS
//
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func ReadProperties(rows *sql.Rows, a *Property) error {
	err := rows.Scan(
		&a.PRID,
		&a.ROLID,
		&a.RSLID,
		&a.FlowState,
		&a.Name,
		&a.YearFounded,
		&a.ParentCompany,
		&a.URL,
		&a.Symbol,
		&a.Price,
		&a.DownPayment,
		&a.RentableArea,
		&a.RentableAreaUnits,
		&a.LotSize,
		&a.LotSizeUnits,
		&a.CapRate,
		&a.AvgCap,
		&a.BuildDate,
		&a.FLAGS,
		&a.Ownership,
		&a.TenantTradeName,
		&a.LeaseGuarantor,
		&a.LeaseType,
		&a.OriginalLeaseTerm,
		&a.RentCommencementDt,
		&a.LeaseExpirationDt,
		&a.TermRemainingOnLease,
		&a.TermRemainingOnLeaseUnits,
		&a.Address,
		&a.Address2,
		&a.City,
		&a.State,
		&a.PostalCode,
		&a.Country,
		&a.LLResponsibilities,
		&a.NOI,
		&a.HQCity,
		&a.HQState,
		&a.Img1,
		&a.Img2,
		&a.Img3,
		&a.Img4,
		&a.Img5,
		&a.Img6,
		&a.Img7,
		&a.Img8,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateProperty updates an existing Property record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateProperty(ctx context.Context, a *Property) error {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.ROLID,
		a.RSLID,
		a.FlowState,
		a.Name,
		a.YearFounded,
		a.ParentCompany,
		a.URL,
		a.Symbol,
		a.Price,
		a.DownPayment,
		a.RentableArea,
		a.RentableAreaUnits,
		a.LotSize,
		a.LotSizeUnits,
		a.CapRate,
		a.AvgCap,
		a.BuildDate,
		a.FLAGS,
		a.Ownership,
		a.TenantTradeName,
		a.LeaseGuarantor,
		a.LeaseType,
		a.OriginalLeaseTerm,
		a.RentCommencementDt,
		a.LeaseExpirationDt,
		a.TermRemainingOnLease,
		a.TermRemainingOnLeaseUnits,
		a.Address,
		a.Address2,
		a.City,
		a.State,
		a.PostalCode,
		a.Country,
		a.LLResponsibilities,
		a.NOI,
		a.HQCity,
		a.HQState,
		a.Img1,
		a.Img2,
		a.Img3,
		a.Img4,
		a.Img5,
		a.Img6,
		a.Img7,
		a.Img8,
		sess.UID,
		a.PRID,
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, Wdb.Prepstmt.UpdateProperty, fields)
	return updateError(err, "Property", *a)
}
