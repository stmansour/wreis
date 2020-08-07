package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"wreis/session"
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
	var err error

	if err = genericDelete(ctx, "RentStep", Wdb.Prepstmt.DeleteRentStepsMembers, id); err != nil {
		return err
	}
	err = genericDelete(ctx, "RentSteps", Wdb.Prepstmt.DeleteRentSteps, id)
	return err
}

// GetRentSteps reads and returns a RentSteps structure
//
// INPUTS
// ctx - db context
// id - RSLID of the record to read
// members - bool - if true, get all the items
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetRentSteps(ctx context.Context, id int64, members bool) (RentSteps, error) {
	var a RentSteps
	var err error

	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	var row *sql.Row
	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Wdb.Prepstmt.GetRentSteps, fields)
	if stmt != nil {
		defer stmt.Close() // it's a txn.  We need to close the statement
	}
	if err = ReadRentSteps(row, &a); err != nil {
		return a, err
	}

	if members {
		stmt2, rows, err := getRowsFromDB(ctx, Wdb.Prepstmt.GetRentStepsItems, fields)
		if err != nil {
			return a, err
		}
		if stmt2 != nil {
			defer stmt2.Close()
		}
		for i := 0; rows.Next(); i++ {
			var x RentStep
			if err = ReadRentStepItem(rows, &x); err != nil {
				return a, err
			}
			a.RS = append(a.RS, x)
		}
		if err = rows.Err(); err != nil {
			return a, err
		}
	}
	return a, nil
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
	var err error
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.RSLID, ErrSessionRequired
	}
	fields := []interface{}{
		a.FLAGS,
		a.CreateBy,
		sess.UID,
	}
	a.CreateBy, a.LastModBy, a.RSLID, err = genericInsert(ctx, "RentSteps", Wdb.Prepstmt.InsertRentSteps, fields, a)
	if err = insertRentStepsList(ctx, a); err != nil {
		return a.RSLID, err
	}
	return a.RSLID, err
}

func insertRentStepsList(ctx context.Context, a *RentSteps) error {
	var err error
	l := len(a.RS)
	for i := 0; i < l; i++ {
		a.RS[i].RSLID = a.RSLID
		if _, err = InsertRentStep(ctx, &a.RS[i]); err != nil {
			return err
		}
	}
	return nil
}

// InsertRentStepsWithList writes a new RentSteps record to the database.
// This function deletes all RentSteps first, then replaces the rentsteps
// with those supplied to this call.
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
	var id int64
	var err error
	if err = DeleteRentSteps(ctx, a.RSLID); err != nil {
		return id, err
	}
	if id, err = InsertRentSteps(ctx, a); err != nil {
		return id, err
	}
	return id, insertRentStepsList(ctx, a)
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
	fmt.Printf("ReadRentSteps: A\n")
	err := row.Scan(
		&a.RSLID,
		&a.FLAGS,
		&a.CreateTS,
		&a.CreateBy,
		&a.LastModTime,
		&a.LastModBy)
	fmt.Printf("ReadRentSteps: B\n")
	if err != nil {
		fmt.Printf("Error from row.Scan: %s\n", err.Error())
	}
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
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.FLAGS,
		sess.UID,
		a.RSLID,
	}
	if a.RSLID > 0 {
		if err = DeleteRentSteps(ctx, a.RSLID); err != nil {
			return err
		}
	}
	l := len(a.RS)
	for i := 0; i < l; i++ {
		a.RS[i].RSLID = a.RSLID // ensure it's the correct list
		if _, err = InsertRentStep(ctx, &a.RS[i]); err != nil {
			return err
		}
	}

	a.LastModBy, err = genericUpdate(ctx, Wdb.Prepstmt.UpdateRentSteps, fields)
	return updateError(err, "RentSteps", *a)
}
