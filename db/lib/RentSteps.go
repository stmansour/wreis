package db

import (
	"context"
	"database/sql"
	"mojo/util"
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

	if err = ValidateSessionForDBDelete(ctx); err != nil {
		return err
	}

	fields := []interface{}{id}
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(Wdb.Prepstmt.DeleteRentSteps)
		defer stmt.Close()
		_, err = stmt.Exec(fields...)
	} else {
		_, err = Wdb.Prepstmt.DeleteRentSteps.Exec(fields...)
	}
	if err != nil {
		util.Ulog("Error deleting RentSteps id=%d error: %v\n", id, err)
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
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(Wdb.Prepstmt.GetRentSteps)
		defer stmt.Close()
		row = stmt.QueryRow(fields...)
	} else {
		row = Wdb.Prepstmt.GetRentSteps.QueryRow(fields...)
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

	// transaction... context
	fields := []interface{}{
		a.CreateBy,
		a.LastModBy,
	}

	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(Wdb.Prepstmt.InsertRentSteps)
		defer stmt.Close()
		res, err = stmt.Exec(fields...)
	} else {
		res, err = Wdb.Prepstmt.InsertRentSteps.Exec(fields...)
	}

	// After getting result...
	if nil == err {
		x, err := res.LastInsertId()
		if err == nil {
			id = int64(x)
			a.RSLID = id
		}
	} else {
		err = insertError(err, "RentSteps", *a)
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
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(Wdb.Prepstmt.UpdateRentSteps)
		defer stmt.Close()
		_, err = stmt.Exec(fields...)
	} else {
		_, err = Wdb.Prepstmt.UpdateRentSteps.Exec(fields...)
	}
	return updateError(err, "RentSteps", *a)
}
