package db

import (
	"context"
	"database/sql"
	"mojo/util"
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
	var err error

	if err = ValidateSessionForDBDelete(ctx); err != nil {
		return err
	}

	fields := []interface{}{id}
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(Wdb.Prepstmt.DeleteRentStep)
		defer stmt.Close()
		_, err = stmt.Exec(fields...)
	} else {
		_, err = Wdb.Prepstmt.DeleteRentStep.Exec(fields...)
	}
	if err != nil {
		util.Ulog("Error deleting RentStep id=%d error: %v\n", id, err)
	}
	return err
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

	var row *sql.Row
	fields := []interface{}{id}
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(Wdb.Prepstmt.GetRentStep)
		defer stmt.Close()
		row = stmt.QueryRow(fields...)
	} else {
		row = Wdb.Prepstmt.GetRentStep.QueryRow(fields...)
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
	var id = int64(0)
	var err error
	var res sql.Result

	if err = ValidateSessionForDBInsert(ctx, &a.CreateBy, &a.LastModBy); err != nil {
		return id, err
	}

	// transaction... context
	fields := []interface{}{
		a.RSLID,
		a.Dt,
		a.Rent,
		a.CreateBy,
		a.LastModBy,
	}

	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(Wdb.Prepstmt.InsertRentStep)
		defer stmt.Close()
		res, err = stmt.Exec(fields...)
	} else {
		res, err = Wdb.Prepstmt.InsertRentStep.Exec(fields...)
	}

	// After getting result...
	if nil == err {
		x, err := res.LastInsertId()
		if err == nil {
			id = int64(x)
			a.RSID = id
		}
	} else {
		err = insertError(err, "RentStep", *a)
	}
	return id, err
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
	var err error

	if err = ValidateSessionForDBUpdate(ctx, &a.LastModBy); err != nil {
		return err
	}

	fields := []interface{}{
		a.RSLID,
		a.Dt,
		a.Rent,
		a.LastModBy,
		a.RSID,
	}
	if tx, ok := TxFromContext(ctx); ok { // if transaction is supplied
		stmt := tx.Stmt(Wdb.Prepstmt.UpdateRentStep)
		defer stmt.Close()
		_, err = stmt.Exec(fields...)
	} else {
		_, err = Wdb.Prepstmt.UpdateRentStep.Exec(fields...)
	}
	return updateError(err, "RentStep", *a)
}
