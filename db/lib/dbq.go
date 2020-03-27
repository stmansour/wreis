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

	GetRentSteps    *sql.Stmt
	InsertRentSteps *sql.Stmt
	UpdateRentSteps *sql.Stmt
	DeleteRentSteps *sql.Stmt
}
