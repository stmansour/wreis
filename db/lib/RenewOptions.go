package db

import (
	"context"
	"database/sql"
	"time"
)

// RenewOptions defines a date and a rent amount for a property. A RenewOptions record
// is part of a group or list. The group is defined by the RSLID
//-----------------------------------------------------------------------------
type RenewOptions struct {
	RSLID       int64         // id of RenewOptionsList to which this record belongs
	FLAGS       uint64        // 1<<0
	LastModTime time.Time     // when was the record last written
	LastModBy   int64         // id of user that did the modify
	CreateTS    time.Time     // when was this record created
	CreateBy    int64         // id of user that created it
	RS          []RenewOption // associated slice of RenewOption records
}

// DeleteRenewOptions deletes the RenewOptions with the specified id from the database
//
// INPUTS
// ctx - db context
// id - RSLID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteRenewOptions(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Property", Wdb.Prepstmt.DeleteRenewOptions, id)
}

// GetRenewOptions reads and returns a RenewOptions structure
//
// INPUTS
// ctx - db context
// id - RSLID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetRenewOptions(ctx context.Context, id int64) (RenewOptions, error) {
	var a RenewOptions
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	var row *sql.Row
	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Wdb.Prepstmt.GetRenewOptions, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadRenewOptions(row, &a)
}

// InsertRenewOptions writes a new RenewOptions record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertRenewOptions(ctx context.Context, a *RenewOptions) (int64, error) {
	fields := []interface{}{
		a.FLAGS,
		a.CreateBy,
		a.LastModBy,
	}

	var err error
	a.CreateBy, a.LastModBy, a.RSLID, err = genericInsert(ctx, "RenewOptions", Wdb.Prepstmt.InsertRenewOptions, fields, a)
	return a.RSLID, err
}

// ReadRenewOptions reads a full RenewOptions structure of data from the database based
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
func ReadRenewOptions(row *sql.Row, a *RenewOptions) error {
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

// UpdateRenewOptions updates an existing RenewOptions record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateRenewOptions(ctx context.Context, a *RenewOptions) error {
	fields := []interface{}{
		a.FLAGS,
		a.LastModBy,
		a.RSLID,
	}

	var err error
	a.LastModBy, err = genericUpdate(ctx, Wdb.Prepstmt.UpdateRenewOptions, fields)
	return updateError(err, "RenewOptions", *a)
}
