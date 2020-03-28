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
	RSLID       int64     // id of RentStepsList to which this record belongs
	LastModTime time.Time // when was the record last written
	LastModBy   int64     // id of user that did the modify
	CreateTS    time.Time // when was this record created
	CreateBy    int64     // id of user that created it
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
	var err error
	var stmt *sql.Stmt

	if err = ValidateSessionForDBDelete(ctx); err != nil {
		return err
	}

	fields := []interface{}{id}
	if stmt, err = deleteDBRow(ctx, "RentSteps", Wdb.Prepstmt.DeleteRentSteps, fields); stmt != nil {
		defer stmt.Close()
	}

	return err
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
	var id = int64(0)
	var err error
	var res sql.Result

	if err = ValidateSessionForDBInsert(ctx, &a.CreateBy, &a.LastModBy); err != nil {
		return id, err
	}

	fields := []interface{}{
		a.CreateBy,
		a.LastModBy,
	}

	stmt, res, err := insertRowToDB(ctx, Wdb.Prepstmt.InsertRentSteps, fields)
	if stmt != nil {
		defer stmt.Close()
	}

	a.RSLID, err = getIDFromResult("RentSteps", res, a, err)
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
	var err error

	if err = ValidateSessionForDBUpdate(ctx, &a.LastModBy); err != nil {
		return err
	}

	fields := []interface{}{
		a.LastModBy,
		a.RSLID,
	}
	err = updateDBRow(ctx, Wdb.Prepstmt.UpdateRentSteps, fields)

	return updateError(err, "RentSteps", *a)
}
