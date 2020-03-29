package db

import (
	"context"
	"database/sql"
	"time"
)

// RentStep defines a date and a rent amount for a property. A RentStep record
// is part of a group or list. The group is defined by the RSLID
//-----------------------------------------------------------------------------
type RentStep struct {
	RSID        int64     // unique id for this record
	RSLID       int64     // id of RentStepList to which this record belongs
	Dt          time.Time // date for the rent amount
	Rent        float64   // amount of rent on the associated date
	LastModTime time.Time // when was the record last written
	LastModBy   int64     // id of user that did the modify
	CreateTS    time.Time // when was this record created
	CreateBy    int64     // id of user that created it
}

// DeleteRentStep deletes the RentStep with the specified id from the database
//
// INPUTS
// ctx - db context
// id - RSID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteRentStep(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Property", Wdb.Prepstmt.DeleteRentStep, id)
}

// GetRentStep reads and returns a RentStep structure
//
// INPUTS
// ctx - db context
// id - RSID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetRentStep(ctx context.Context, id int64) (RentStep, error) {
	var a RentStep
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}
	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Wdb.Prepstmt.GetRentStep, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadRentStep(row, &a)
}

// InsertRentStep writes a new RentStep record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertRentStep(ctx context.Context, a *RentStep) (int64, error) {
	// transaction... context
	fields := []interface{}{
		a.RSLID,
		a.Dt,
		a.Rent,
		a.CreateBy,
		a.LastModBy,
	}
	var err error
	a.CreateBy, a.LastModBy, a.RSID, err = genericInsert(ctx, "RentStep", Wdb.Prepstmt.InsertRentStep, fields, a)
	return a.RSID, err
}

// ReadRentStep reads a full RentStep structure of data from the database based
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
func ReadRentStep(row *sql.Row, a *RentStep) error {
	err := row.Scan(
		&a.RSID,
		&a.RSLID,
		&a.Dt,
		&a.Rent,
		&a.CreateTS,
		&a.CreateBy,
		&a.LastModTime,
		&a.LastModBy)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateRentStep updates an existing RentStep record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateRentStep(ctx context.Context, a *RentStep) error {
	fields := []interface{}{
		a.RSLID,
		a.Dt,
		a.Rent,
		a.LastModBy,
		a.RSID,
	}

	var err error
	a.LastModBy, err = genericUpdate(ctx, Wdb.Prepstmt.UpdateRentStep, fields)
	return updateError(err, "RentStep", *a)
}
