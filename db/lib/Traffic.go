package db

import (
	"context"
	"database/sql"
	"time"
	"wreis/session"
)

// Traffic defines a date and a Count amount for a property. A Traffic record
// is associated with a property
//-----------------------------------------------------------------------------
type Traffic struct {
	TID         int64     // unique id for this record
	PRID        int64     // Property with which this record is associated
	FLAGS       uint64    // 1<<0 :  0 -> use Description text, 1 -> use date
	Count       int64     // amount of Count on the associated date
	Description string    // Descriptionion comment:  "years 1 - 2" etc.
	LastModTime time.Time // when was the record last written
	LastModBy   int64     // id of user that did the modify
	CreateTS    time.Time // when was this record created
	CreateBy    int64     // id of user that created it
}

// DeleteTraffic deletes the Traffic with the specified id from the database
//
// INPUTS
// ctx - db context
// id - TID of the record to delete
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteTraffic(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Traffic", Wdb.Prepstmt.DeleteTraffic, id)
}

// DeleteTrafficItems deletes the Traffic items associated with PRID = id
//
// INPUTS
// ctx - db context
// id - PRID of the record to read - this will delete
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteTrafficItems(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Traffic", Wdb.Prepstmt.DeleteTrafficItems, id)
}

// GetTraffic reads and returns a Traffic structure
//
// INPUTS
// ctx - db context
// id - TID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetTraffic(ctx context.Context, id int64) (Traffic, error) {
	var a Traffic
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}
	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Wdb.Prepstmt.GetTraffic, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadTraffic(row, &a)
}

// GetTrafficItems reads only the array of Traffic items associated with
// the supplied PRID
//
// INPUTS
// ctx - db context
// id - PRID to which all the traffic items belong
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetTrafficItems(ctx context.Context, id int64) ([]Traffic, error) {
	var err error
	var a []Traffic
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt2, rows, err := getRowsFromDB(ctx, Wdb.Prepstmt.GetTrafficItems, fields)
	if err != nil {
		return a, err
	}
	if stmt2 != nil {
		defer stmt2.Close()
	}
	for i := 0; rows.Next(); i++ {
		var x Traffic
		if err = ReadTrafficItem(rows, &x); err != nil {
			return a, err
		}
		a = append(a, x)
	}
	if err = rows.Err(); err != nil {
		return a, err
	}
	return a, nil
}

// InsertTraffic writes a new Traffic record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertTraffic(ctx context.Context, a *Traffic) (int64, error) {
	// transaction... context
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.TID, ErrSessionRequired
	}
	fields := []interface{}{
		a.PRID,
		a.FLAGS,
		a.Count,
		a.Description,
		a.CreateBy,
		sess.UID,
	}
	var err error
	a.CreateBy, a.LastModBy, a.TID, err = genericInsert(ctx, "Traffic", Wdb.Prepstmt.InsertTraffic, fields, a)
	return a.TID, err
}

// ReadTraffic reads a full Traffic structure of data from the database based
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
func ReadTraffic(row *sql.Row, a *Traffic) error {
	err := row.Scan(
		&a.TID,
		&a.PRID,
		&a.FLAGS,
		&a.Count,
		&a.Description,
		&a.CreateTS,
		&a.CreateBy,
		&a.LastModTime,
		&a.LastModBy)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadTrafficItem reads a full Traffic structure of data from the database based
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
func ReadTrafficItem(rows *sql.Rows, a *Traffic) error {
	err := rows.Scan(
		&a.TID,
		&a.PRID,
		&a.FLAGS,
		&a.Count,
		&a.Description,
		&a.CreateTS,
		&a.CreateBy,
		&a.LastModTime,
		&a.LastModBy)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateTraffic updates an existing Traffic record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateTraffic(ctx context.Context, a *Traffic) error {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.PRID,
		a.FLAGS,
		a.Count,
		a.Description,
		sess.UID,
		a.TID,
	}

	var err error
	a.LastModBy, err = genericUpdate(ctx, Wdb.Prepstmt.UpdateTraffic, fields)
	return updateError(err, "Traffic", *a)
}
