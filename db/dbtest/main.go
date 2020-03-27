package main

import (
	"context"
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"os"
	"time"
	db "wreis/db/lib"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global application structure
var App struct {
	db     *sql.DB // wreis db
	dbUser string  // db user name
	dbName string  // db name
	dbPort int     // db port
	NoAuth bool
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbrrPtr := flag.String("M", "wreis", "database name (wreis)")
	portPtr := flag.Int("p", 8275, "port on which WREIS server listens")
	noauth := flag.Bool("noauth", false, "if specified, inhibit authentication")

	flag.Parse()

	App.dbUser = *dbuPtr
	App.dbPort = *portPtr
	App.dbName = *dbrrPtr
	App.NoAuth = *noauth
}

func main() {
	var err error
	readCommandLineArgs()

	//----------------------------
	// Open RentRoll database
	//----------------------------
	if err = db.ReadConfig(); err != nil {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.Wdb.Config.WREISDbname, db.Wdb.Config.WREISDbuser, err)
		os.Exit(1)
	}

	s := extres.GetSQLOpenString(App.dbName, &db.Wdb.Config)
	App.db, err = sql.Open("mysql", s)
	if nil != err {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", App.dbName, App.dbUser, err)
		os.Exit(1)
	}
	defer App.db.Close()
	err = App.db.Ping()
	if nil != err {
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", App.dbName, App.dbUser, err)
		os.Exit(1)
	}
	db.Init(App.db)
	db.SessionInit(10)

	//------------------------------------------------------------------------
	// Create a session that this process can use for accessing the database
	//------------------------------------------------------------------------
	now := time.Now()
	ctx := context.Background()
	expire := now.Add(10 * time.Minute)
	sess := db.SessionNew(
		"dbtest-app"+fmt.Sprintf("%010x", expire.Unix()), // token
		"dbtest",      // username
		"dbtest-app",  // name string
		int64(-99998), // uid
		"",            // image url
		-1,            // security role id
		&expire)       // expiredt
	ctx = db.SetSessionContextKey(ctx, sess)

	TestRentStep(ctx)
}

// TestRentStep checks the basic db functions for the RentStep struct
//-----------------------------------------------------------------------------
func TestRentStep(ctx context.Context) {
	var err error
	rs := db.RentStep{
		RSID:  0,
		RSLID: 1,
		Dt:    time.Date(2020, time.March, 23, 0, 0, 0, 0, time.UTC),
		Rent:  float64(2750.00),
	}
	var delid, id int64
	if id, err = db.InsertRentStep(ctx, &rs); err != nil {
		fmt.Printf("Error inserting RentStep: %s\n", err)
		os.Exit(1)
	}

	// Insert another for delete...
	if delid, err = db.InsertRentStep(ctx, &rs); err != nil {
		fmt.Printf("Error inserting RentStep: %s\n", err)
		os.Exit(1)
	}
	if err = db.DeleteRentStep(ctx, delid); err != nil {
		fmt.Printf("Error deleting RentStep: %s\n", err)
		os.Exit(1)
	}

	var rs1 db.RentStep
	if rs1, err = db.GetRentStep(ctx, id); err != nil {
		fmt.Printf("error in GetRentStep: %s\n", err.Error())
		os.Exit(1)
	}
	rs1.Rent += float64(10)
	if err = db.UpdateRentStep(ctx, &rs1); err != nil {
		fmt.Printf("Error updating RentStep: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update RentStep\n")
}
