// faa  a program to scrape the FAA directory site.
package main

import (
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"phonebook/lib"
	"strings"
	db "wreis/db/lib"
	"wreis/session"
	util "wreis/util/lib"
	"wreis/ws"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global data structure for this app
var App struct {
	db        *sql.DB
	DBName    string
	DBUser    string
	Port      int      // port on which wreis listens
	LogFile   *os.File // where to log messages
	fname     string
	startName string
}

// HomeHandler serves static http content such as the .css files
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		Chttp.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/home/", http.StatusFound)
	}
}

// Chttp is a server mux for handling unprocessed html page requests.
// For example, a .css file or an image file.
var Chttp = http.NewServeMux()

func initHTTP() {
	Chttp.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/home/", ws.HomeUIHandler)
	http.HandleFunc("/v1/", ws.V1ServiceHandler)
}

func readCommandLineArgs() {
	portPtr := flag.Int("p", 8276, "port on which wreis server listens")
	vptr := flag.Bool("v", false, "Show version, then exit")
	flag.Parse()
	if *vptr {
		fmt.Printf("Version:   %s\n", ws.GetVersionNo())
		os.Exit(0)
	}
	App.Port = *portPtr
}

func main() {
	var err error
	readCommandLineArgs()
	err = db.ReadConfig()
	if err != nil {
		fmt.Printf("Error in db.ReadConfig: %s\n", err.Error())
		os.Exit(1)
	}

	//==============================================
	// Open the logfile and begin logging...
	//==============================================
	App.LogFile, err = os.OpenFile("wreis.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	lib.Errcheck(err)
	defer App.LogFile.Close()
	log.SetOutput(App.LogFile)
	util.Ulog("*** WREIS WERTZ REAL ESTATE INVESTMENT SERVICES ***\n")
	util.Ulog("Using database: %s , host = %s, port = %d\n", db.Wdb.Config.WREISDbname, db.Wdb.Config.WREISDbhost, db.Wdb.Config.WREISDbport)

	// Get the database...
	// s := "<awsdbusername>:<password>@tcp(<rdsinstancename>:3306)/accord"
	s := extres.GetSQLOpenString(db.Wdb.Config.WREISDbname, &db.Wdb.Config)
	App.db, err = sql.Open("mysql", s)
	if nil != err {

		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.Wdb.Config.WREISDbname, db.Wdb.Config.WREISDbuser, err)
		os.Exit(1)
	}
	util.Ulog("successfully opened database %q as user %q on %s\n", db.Wdb.Config.WREISDbname, db.Wdb.Config.WREISDbuser, db.Wdb.Config.WREISDbhost)
	defer App.db.Close()

	err = App.db.Ping()
	if nil != err {
		util.Ulog("could not ping database %q as user %q on %s\n", db.Wdb.Config.WREISDbname, db.Wdb.Config.WREISDbuser, db.Wdb.Config.WREISDbhost)
		util.Ulog("error: %s\n", err.Error())
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", db.Wdb.Config.WREISDbname, db.Wdb.Config.WREISDbuser, err)
		os.Exit(1)
	}
	db.Init(App.db)                 // initializes database
	session.Init(10, db.Wdb.Config) // we must have login sessions
	db.BuildPreparedStatements()    // the prepared statement for db access
	initHTTP()
	util.Ulog("wreis initiating HTTP service on port %d\n", App.Port)
	fmt.Printf("Using database: %s , host = %s, port = %d\n", db.Wdb.Config.WREISDbname, db.Wdb.Config.WREISDbhost, db.Wdb.Config.WREISDbport)

	//go http.ListenAndServeTLS(fmt.Sprintf(":%d", App.Port+1), App.CertFile, App.KeyFile, nil)
	err = http.ListenAndServe(fmt.Sprintf(":%d", App.Port), nil)
	if nil != err {
		fmt.Printf("*** Error on http.ListenAndServe: %v\n", err)
		util.Ulog("*** Error on http.ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
