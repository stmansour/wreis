package main

import (
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"os"
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
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.DBConfig.WREISDbname, db.DBConfig.WREISDbuser, err)
		os.Exit(1)
	}

	s := extres.GetSQLOpenString(App.dbName, &db.DBConfig)
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

	fmt.Printf("Success!\n")

}
