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
	if errlist = wcsv.ReadPropertyFile(App.fname, PropertyHandler); err != nil {
		if len(errlist) > 0 {
			fmt.Printf("csv.ReadPropertyFile returned %d errors\n", len(errlist))
			for i := 0; i < len(errlist); i++ {
				fmt.Printf("%d. %s\n", i, errlist[i].Error())
			}
			os.Exit(1)
		}
	}
}

// PropertyHandler is called for each record of a Property csv file.
//
// INPUTS
// csvctx - context for this csv file, used to determine which column contains
//          what information.
// ss     - array of strings, one for each column in the csv file
// linno  - line number in the csvfile
//-----------------------------------------------------------------------------
func PropertyHandler(csvctx wcsv.Context, ss []string, lineno int) []error {
	var p db.Property
	var errlist []error

	for i := 0; i < len(csvctx.Order); i++ {
		switch i {
		case wcsv.PRName:
			p.Name = ss[csvctx.Order[wcsv.PRName]]
		case wcsv.PRYearsInBusiness:
			p.YearsInBusiness, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRParentCompany:
			p.ParentCompany = ss[csvctx.Order[i]]
		case wcsv.PRURL:
			p.URL = ss[csvctx.Order[i]]
		case wcsv.PRSymbol:
			p.Symbol = ss[csvctx.Order[i]]
		case wcsv.PRPrice:
			p.Price, errlist = wcsv.ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRDownPayment:
			p.DownPayment, errlist = wcsv.ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRRentableArea:
			p.RentableArea, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRRentableAreaUnits:
			p.RentableAreaUnits, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRLotSize:
			p.LotSize, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRLotSizeUnits:
			p.LotSizeUnits, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRCapRate:
			p.CapRate, errlist = wcsv.ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRAvgCap:
			p.AvgCap, errlist = wcsv.ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRBuildDate:
			p.BuildDate, errlist = wcsv.ParseDate(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PROwnership:
			p.Ownership, errlist = wcsv.ParseInt(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRTenantTradeName:
			p.TenantTradeName = ss[csvctx.Order[i]]
		case wcsv.PRLeaseGuarantor:
			p.LeaseGuarantor, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRLeaseType:
			p.LeaseType, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRDeliveryDt:
			p.DeliveryDt, errlist = wcsv.ParseDate(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PROriginalLeaseTerm:
			p.OriginalLeaseTerm, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRLeaseCommencementDt:
			p.LeaseCommencementDt, errlist = wcsv.ParseDate(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRLeaseExpirationDt:
			p.LeaseExpirationDt, errlist = wcsv.ParseDate(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRTermRemainingOnLease:
			p.TermRemainingOnLease, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRROID:
			p.ROID, errlist = wcsv.ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRAddress:
			p.Address = ss[csvctx.Order[i]]
		case wcsv.PRAddress2:
			p.Address2 = ss[csvctx.Order[i]]
		case wcsv.PRCity:
			p.City = ss[csvctx.Order[i]]
		case wcsv.PRState:
			p.State = ss[csvctx.Order[i]]
		case wcsv.PRPostalCode:
			p.PostalCode = ss[csvctx.Order[i]]
		case wcsv.PRCountry:
			p.Country = ss[csvctx.Order[i]]
		case wcsv.PRLLResponsibilities:
			p.LLResponsibilities = ss[csvctx.Order[i]]
		case wcsv.PRNOI:
			p.NOI, errlist = wcsv.ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case wcsv.PRHQAddress:
			p.HQAddress = ss[csvctx.Order[i]]
		case wcsv.PRHQAddress2:
			p.HQAddress2 = ss[csvctx.Order[i]]
		case wcsv.PRHQCity:
			p.HQCity = ss[csvctx.Order[i]]
		case wcsv.PRHQState:
			p.HQState = ss[csvctx.Order[i]]
		case wcsv.PRHQPostalCode:
			p.HQPostalCode = ss[csvctx.Order[i]]
		case wcsv.PRHQCountry:
			p.HQCountry = ss[csvctx.Order[i]]
		}
	}

	fmt.Printf("Line: %d p = %#v\n", lineno, p)

	return errlist
}
