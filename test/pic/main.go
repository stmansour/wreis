package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	db "wreis/db/lib"
	util "wreis/util/lib"
	"wreis/ws"
)

// App is the global application structure.
// We need this because the db reads in all the credentials we need to access
// AWS S3.
//-----------------------------------------------------------------------------
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
	portPtr := flag.Int("p", 8276, "port on which WREIS server listens")
	noauth := flag.Bool("noauth", false, "if specified, inhibit authentication")

	flag.Parse()

	App.dbUser = *dbuPtr
	App.dbPort = *portPtr
	App.dbName = *dbrrPtr
	App.NoAuth = *noauth
}

func main() {
	readCommandLineArgs()

	//----------------------------
	// Open RentRoll database
	//----------------------------
	if err := db.ReadConfig(); err != nil {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.Wdb.Config.WREISDbname, db.Wdb.Config.WREISDbuser, err)
		os.Exit(1)
	}

	// s := extres.GetSQLOpenString(App.dbName, &db.Wdb.Config)
	// App.db, err = sql.Open("mysql", s)
	// if nil != err {
	// 	fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", App.dbName, App.dbUser, err)
	// 	os.Exit(1)
	// }
	// defer App.db.Close()
	// err = App.db.Ping()
	// if nil != err {
	// 	fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", App.dbName, App.dbUser, err)
	// 	os.Exit(1)
	// }
	// db.Init(App.db)

	filename := "roller-32.png"
	// open the file for use
	usrfile, err := os.Open(filename)
	if err != nil {
		util.Console("Error opening %s: %s\n", filename, err.Error())
		os.Exit(1)
	}
	defer usrfile.Close()

	PRID := int64(4)
	idx := int64(3)

	var path, url string
	if path, url, err = ws.UploadImageFileToS3(filename, usrfile, PRID, idx); err != nil {
		util.Console("UploadImageFileToS3 error: %s\n", err.Error())
		os.Exit(1)
	}

	util.Console("success! CREATED: path = %s, URL = %s\n", path, url)

	//---------------------------------------------------------------
	// now read in a different image and overwrite the old image...
	//---------------------------------------------------------------
	fname2 := "receipts-32.png"
	file2, err := os.Open(fname2)
	if err != nil {
		util.Console("Error opening %s: %s\n", fname2, err.Error())
		os.Exit(1)
	}
	defer file2.Close()

	//---------------------------------------------------------------
	// use the same filename and url, but change the image...
	//---------------------------------------------------------------
	if path, url, err = ws.UploadImageFileToS3(filename, file2, PRID, idx); err != nil {
		util.Console("UploadImageFileToS3 error: %s\n", err.Error())
		os.Exit(1)
	}

	util.Console("GetImageFilenameFromURL reports filename is: %s\n", ws.GetImageFilenameFromURL(url))

	util.Console("success! UPDATED: path = %s, URL = %s\n", path, url)

	//---------------------------------------------------------------
	// now, delete that object...
	//---------------------------------------------------------------
	if err = ws.DeleteS3ImageFile(filename, PRID, idx); err != nil {
		util.Console("DeleteS3ImageFile error: %s\n", err.Error())
	}

	util.Console("success! DELETED: path = %s, URL = %s\n", path, url)

}
