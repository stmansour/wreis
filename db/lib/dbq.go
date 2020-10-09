package db

import "database/sql"

// PrepSQL is the structure containing all the prepared statements
type PrepSQL struct {
	GetProperty    *sql.Stmt
	InsertProperty *sql.Stmt
	UpdateProperty *sql.Stmt
	DeleteProperty *sql.Stmt

	GetRentStep    *sql.Stmt
	InsertRentStep *sql.Stmt
	UpdateRentStep *sql.Stmt
	DeleteRentStep *sql.Stmt

	GetRentSteps           *sql.Stmt
	GetRentStepsItems      *sql.Stmt
	InsertRentSteps        *sql.Stmt
	UpdateRentSteps        *sql.Stmt
	DeleteRentStepsMembers *sql.Stmt
	DeleteRentSteps        *sql.Stmt

	GetRenewOption    *sql.Stmt
	InsertRenewOption *sql.Stmt
	UpdateRenewOption *sql.Stmt
	DeleteRenewOption *sql.Stmt

	GetStateInfo         *sql.Stmt
	GetAllStateInfoItems *sql.Stmt
	InsertStateInfo      *sql.Stmt
	UpdateStateInfo      *sql.Stmt
	DeleteStateInfo      *sql.Stmt

	GetRenewOptions           *sql.Stmt
	GetRenewOptionsItems      *sql.Stmt
	InsertRenewOptions        *sql.Stmt
	UpdateRenewOptions        *sql.Stmt
	DeleteRenewOptions        *sql.Stmt
	DeleteRenewOptionsMembers *sql.Stmt

	GetTraffic         *sql.Stmt
	GetTrafficItems    *sql.Stmt
	InsertTraffic      *sql.Stmt
	UpdateTraffic      *sql.Stmt
	DeleteTraffic      *sql.Stmt
	DeleteTrafficItems *sql.Stmt
}
