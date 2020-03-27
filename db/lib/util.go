package db

import (
	"context"
	"database/sql"
	"extres"
	"fmt"
	"mojo/util"
	"runtime/debug"
)

// SkipSQLNoRowsError assing nil to original err variable
// if its kind of no rows in result error from sql package
func SkipSQLNoRowsError(err *error) {
	if IsSQLNoResultsError(*err) {
		*err = nil
	}
}

// IsSQLNoResultsError returns true if the error provided is a sql err indicating no rows in the solution set.
func IsSQLNoResultsError(err error) bool {
	return err == sql.ErrNoRows
}

// Errcheck - saves a bunch of typing, prints error if it exists
//            and provides a traceback as well
// Note that the error is printed only if the environment is NOT production.
func Errcheck(err error) {
	if err != nil {
		if IsSQLNoResultsError(err) {
			return
		}
		if extres.APPENVPROD != Wdb.Config.Env {
			fmt.Printf("error = %v\n", err)
		}
		debug.PrintStack()
		// log.Fatal(err)
	}
}

// ValidateSession validates that we either have a valid session or that a
// session is not required.  If there's a problem it returns an error.
//
// INPUTS
// context - session Context
//
// RETURNS
// true if the session is valid
// false otherwise
//------------------------------------------------------------------------------
func ValidateSession(ctx context.Context) bool {
	if !(Wdb.noAuth && Wdb.Config.Env != extres.APPENVPROD) {
		_, ok := SessionFromContext(ctx)
		if !ok {
			return false
		}
	}
	return true
}

// ValidateSessionForDBDelete validates that we either have a valid session or that a
// session is not required.  If there's a problem it returns an error.
//
// INPUTS
// context - session Context
//
// RETURNS
// true if the session is valid
// false otherwise
//------------------------------------------------------------------------------
func ValidateSessionForDBDelete(ctx context.Context) error {
	if !(Wdb.noAuth && Wdb.Config.Env != extres.APPENVPROD) {
		_, ok := SessionFromContext(ctx)
		if !ok {
			return ErrSessionRequired
		}
	}
	return nil
}

// ValidateSessionForDBInsert is a convenience function that replaces 8 lines
// of code with about 3. Since these lines are needed for every insert call
// it saves a lot of lines.  Added this routine at the time Task,TaskList,
// TaskDescriptor and  TaskListDefinition were added.
//
// INPUTS
// ctx - session Context
// id1 - pointer to creator UID, will be returned
// id2 - pointer to lastmod UID, will be returned
//
// RETURNS
// any error encountered, nil otherwise
//-----------------------------------------------------------------------------
func ValidateSessionForDBInsert(ctx context.Context, id1, id2 *int64) error {
	if !(Wdb.noAuth && Wdb.Config.Env != extres.APPENVPROD) {
		sess, ok := SessionFromContext(ctx)
		if !ok {
			return ErrSessionRequired
		}
		(*id1) = sess.UID
		(*id2) = sess.UID
		return nil
	}
	return nil
}

// ValidateSessionForDBUpdate is a convenience function that replaces 8 lines
// of code with about 4. Since these lines are needed for every update call
// it saves a lot of lines.  Added this routine at the time Task,TaskList,
// TaskDescriptor and  TaskListDefinition were added.
//
// INPUTS
// ctx - session Context
// id2 - pointer to lastmod UID, will be returned
//
// RETURNS
// any error encountered, nil otherwise
//-----------------------------------------------------------------------------
func ValidateSessionForDBUpdate(ctx context.Context, id2 *int64) error {
	if !(Wdb.noAuth && Wdb.Config.Env != extres.APPENVPROD) {
		sess, ok := SessionFromContext(ctx)
		if !ok {
			return ErrSessionRequired
		}
		(*id2) = sess.UID
		return nil
	}
	return nil
}

func updateError(err error, n string, a interface{}) error {
	if nil != err {
		util.Ulog("Update%s: error updating %s:  %v\n", n, n, err)
		util.Ulog("%s = %#v\n", n, a)
	}
	return err
}

func insertError(err error, n string, a interface{}) error {
	if nil != err {
		util.Ulog("Insert%s: error inserting %s:  %v\n", n, n, err)
		util.Ulog("%s = %#v\n", n, a)
	}
	return err
}
