package db

import (
	"context"
	"database/sql"
	"time"
	"wreis/session"
)

// StateInfo defines information about a particular state for a property.
//-----------------------------------------------------------------------------
type StateInfo struct {
	SIID         int64     // unique id for this record
	PRID         int64     // property to which this info belongs
	InitiatorUID int64     // date/time this state was initiated
	InitiatorDt  time.Time // date/time this state was initiated
	ApproverUID  int64     // date/time this state was approved
	ApproverDt   time.Time // date/time this state was approved
	FlowState    int64     // state being described
	FLAGS        uint64    // 1<<0 :  0 -> use Opt text, 1 -> use date
	LastModTime  time.Time // when was the record last written
	LastModBy    int64     // id of user that did the modify
	CreateTime   time.Time // when was this record created
	CreateBy     int64     // id of user that created it
}

// DeleteStateInfo deletes the StateInfo with the specified id from the database
//
// INPUTS
// ctx - db context
// id - SIID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteStateInfo(ctx context.Context, id int64) error {
	return genericDelete(ctx, "StateInfo", Wdb.Prepstmt.DeleteStateInfo, id)
}

// GetStateInfo reads and returns a StateInfo structure
//
// INPUTS
// ctx - db context
// id - SIID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetStateInfo(ctx context.Context, id int64) (StateInfo, error) {
	var a StateInfo
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}
	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Wdb.Prepstmt.GetStateInfo, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadStateInfo(row, &a)
}

// InsertStateInfo writes a new StateInfo record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertStateInfo(ctx context.Context, a *StateInfo) (int64, error) {
	// transaction... context
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.SIID, ErrSessionRequired
	}
	fields := []interface{}{
		a.PRID,
		a.InitiatorUID,
		a.InitiatorDt,
		a.ApproverUID,
		a.ApproverDt,
		a.FlowState,
		a.FLAGS,
		sess.UID,
		sess.UID,
	}
	var err error
	a.CreateBy, a.LastModBy, a.SIID, err = genericInsert(ctx, "StateInfo", Wdb.Prepstmt.InsertStateInfo, fields, a)
	return a.SIID, err
}

// ReadStateInfo reads a full StateInfo structure of data from the database based
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
func ReadStateInfo(row *sql.Row, a *StateInfo) error {
	err := row.Scan(
		&a.SIID,
		&a.PRID,
		&a.InitiatorUID,
		&a.InitiatorDt,
		&a.ApproverUID,
		&a.ApproverDt,
		&a.FlowState,
		&a.FLAGS,
		&a.CreateTime,
		&a.CreateBy,
		&a.LastModTime,
		&a.LastModBy)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadStateInfoItem reads a full StateInfo structure of data from the database based
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
func ReadStateInfoItem(rows *sql.Rows, a *StateInfo) error {
	err := rows.Scan(
		&a.SIID,
		&a.PRID,
		&a.InitiatorUID,
		&a.InitiatorDt,
		&a.ApproverUID,
		&a.ApproverDt,
		&a.FlowState,
		&a.FLAGS,
		&a.CreateTime,
		&a.CreateBy,
		&a.LastModTime,
		&a.LastModBy)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateStateInfo updates an existing StateInfo record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateStateInfo(ctx context.Context, a *StateInfo) error {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.PRID,
		a.InitiatorUID,
		a.InitiatorDt,
		a.ApproverUID,
		a.ApproverDt,
		a.FlowState,
		a.FLAGS,
		sess.UID,
		a.SIID,
	}

	var err error
	a.LastModBy, err = genericUpdate(ctx, Wdb.Prepstmt.UpdateStateInfo, fields)
	return updateError(err, "StateInfo", *a)
}
