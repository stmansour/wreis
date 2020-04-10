package main

import (
	"context"
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"os"
	"time"
	wcsv "wreis/csv/lib"
	db "wreis/db/lib"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global application structure
var App struct {
	db     *sql.DB // wreis db
	dbUser string  // db user name
	dbName string  // db name
	dbPort int     // db port
	NoAuth bool    // skip authorization checks (not implemented at this time)
	fname  string  // csv filename
}

func readCommandLineArgs() {
	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbrrPtr := flag.String("M", "wreis", "database name (wreis)")
	fPtr := flag.String("f", "test.csv", "csv file to import")
	portPtr := flag.Int("p", 8275, "port on which WREIS server listens")
	noauth := flag.Bool("noauth", false, "if specified, inhibit authentication")

	flag.Parse()

	App.dbUser = *dbuPtr
	App.dbPort = *portPtr
	App.dbName = *dbrrPtr
	App.NoAuth = *noauth
	App.fname = *fPtr
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

	//----------------------------
	// process the csv file...
	//----------------------------
	var errlist []error
	if errlist = wcsv.ImportPropertyFile(App.fname); len(errlist) > 0 {
		fmt.Printf("csv.ImportPropertyFile returned %d errors\n", len(errlist))
		for i := 0; i < len(errlist); i++ {
			fmt.Printf("%d. %s\n", i, errlist[i].Error())
		}
		os.Exit(1)
	}
}
