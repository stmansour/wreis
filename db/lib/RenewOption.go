package db

import (
	"context"
	"database/sql"
	"time"
)

// RenewOption defines a date and a rent amount for a property. A RenewOption record
// is part of a group or list. The group is defined by the ROLID
//-----------------------------------------------------------------------------
type RenewOption struct {
	ROID        int64     // unique id for this record
	ROLID       int64     // id of RenewOptionList to which this record belongs
	Dt          time.Time // date for the rent amount, valid when ROLID.FLAGS bit 0 = 1
	Opt         string    // option comment:  "years 1 - 3", etc.
	Rent        float64   // amount of rent on the associated date
	FLAGS       uint64    // 1<<0 :  0 -> count is valid, 1 -> Dt is valid
	LastModTime time.Time // when was the record last written
	LastModBy   int64     // id of user that did the modify
	CreateTS    time.Time // when was this record created
	CreateBy    int64     // id of user that created it
}

// DeleteRenewOption deletes the RenewOption with the specified id from the database
//
// INPUTS
// ctx - db context
// id - ROID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteRenewOption(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Property", Wdb.Prepstmt.DeleteRenewOption, id)
}

// GetRenewOption reads and returns a RenewOption structure
//
// INPUTS
// ctx - db context
// id - ROID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetRenewOption(ctx context.Context, id int64) (RenewOption, error) {
	var a RenewOption
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}
	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Wdb.Prepstmt.GetRenewOption, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadRenewOption(row, &a)
}

// InsertRenewOption writes a new RenewOption record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertRenewOption(ctx context.Context, a *RenewOption) (int64, error) {
	// transaction... context
	fields := []interface{}{
		a.ROLID,
		a.Dt,
		a.Opt,
		a.Rent,
		a.FLAGS,
		a.CreateBy,
		a.LastModBy,
	}
	var err error
	a.CreateBy, a.LastModBy, a.ROID, err = genericInsert(ctx, "RenewOption", Wdb.Prepstmt.InsertRenewOption, fields, a)
	return a.ROID, err
}

// ReadRenewOption reads a full RenewOption structure of data from the database based
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
func ReadRenewOption(row *sql.Row, a *RenewOption) error {
	err := row.Scan(
		&a.ROID,
		&a.ROLID,
		&a.Dt,
		&a.Opt,
		&a.Rent,
		&a.FLAGS,
		&a.CreateTS,
		&a.CreateBy,
		&a.LastModTime,
		&a.LastModBy)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateRenewOption updates an existing RenewOption record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateRenewOption(ctx context.Context, a *RenewOption) error {
	fields := []interface{}{
		a.ROLID,
		a.Dt,
		a.Opt,
		a.Rent,
		a.FLAGS,
		a.LastModBy,
		a.ROID,
	}

	var err error
	a.LastModBy, err = genericUpdate(ctx, Wdb.Prepstmt.UpdateRenewOption, fields)
	return updateError(err, "RenewOption", *a)
}
