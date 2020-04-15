package db

import (
	"context"
	"database/sql"
	"time"
)

// RentSteps defines a date and a rent amount for a property. A RentSteps record
// is part of a group or list. The group is defined by the RSLID
//-----------------------------------------------------------------------------
type RentSteps struct {
	RSLID       int64      // id of RentStepsList to which this record belongs
	FLAGS       uint64     // 1<<0 = 0 means count based, 1 means date based
	LastModTime time.Time  // when was the record last written
	LastModBy   int64      // id of user that did the modify
	CreateTS    time.Time  // when was this record created
	CreateBy    int64      // id of user that created it
	RS          []RentStep // associated slice of RentStep records
}

// DeleteRentSteps deletes the RentSteps with the specified id from the database
//
// INPUTS
// ctx - db context
// id - RSLID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteRentSteps(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Property", Wdb.Prepstmt.DeleteRentSteps, id)
}

// GetRentSteps reads and returns a RentSteps structure
//
// INPUTS
// ctx - db context
// id - RSLID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetRentSteps(ctx context.Context, id int64) (RentSteps, error) {
	var a RentSteps
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	var row *sql.Row
	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Wdb.Prepstmt.GetRentSteps, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadRentSteps(row, &a)
}

// InsertRentSteps writes a new RentSteps record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertRentSteps(ctx context.Context, a *RentSteps) (int64, error) {
	fields := []interface{}{
		a.FLAGS,
		a.CreateBy,
		a.LastModBy,
	}

	var err error
	a.CreateBy, a.LastModBy, a.RSLID, err = genericInsert(ctx, "RentSteps", Wdb.Prepstmt.InsertRentSteps, fields, a)
	return a.RSLID, err
}

// InsertRentStepsWithList writes a new RentSteps record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertRentStepsWithList(ctx context.Context, a *RentSteps) (int64, error) {
	id, err := InsertRentSteps(ctx, a)
	if err != nil {
		return id, err
	}
	l := len(a.RS)
	for i := 0; i < l; i++ {
		a.RS[i].RSLID = id
		if _, err = InsertRentStep(ctx, &a.RS[i]); err != nil {
			return id, err
		}
	}
	return id, err
}

// ReadRentSteps reads a full RentSteps structure of data from the database based
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
func ReadRentSteps(row *sql.Row, a *RentSteps) error {
	err := row.Scan(
		&a.RSLID,
		&a.FLAGS,
		&a.CreateTS,
		&a.CreateBy,
		&a.LastModTime,
		&a.LastModBy)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateRentSteps updates an existing RentSteps record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateRentSteps(ctx context.Context, a *RentSteps) error {
	fields := []interface{}{
		a.FLAGS,
		a.LastModBy,
		a.RSLID,
	}

	var err error
	a.LastModBy, err = genericUpdate(ctx, Wdb.Prepstmt.UpdateRentSteps, fields)
	return updateError(err, "RentSteps", *a)
}
