package main

import (
	"context"
	"fmt"
	"os"
	"time"
	db "wreis/db/lib"
)

// TestProperty checks the basic db functions for the Property struct
//-----------------------------------------------------------------------------
func TestProperty(ctx context.Context) {
	var err error
	dt := time.Date(2020, time.March, 23, 0, 0, 0, 0, time.UTC)
	dur := time.Duration((20 * 365 * 24)) * time.Hour
	rs := db.Property{
		PRID:                      0,
		Name:                      "Bill's Boar Emporium",
		YearsInBusiness:           8,
		ParentCompany:             "",
		URL:                       "http://bbb.com/",
		Symbol:                    "BBE",
		Price:                     float64(12345.67),
		DownPayment:               float64(40000),
		RentableArea:              30000,
		RentableAreaUnits:         1,
		LotSize:                   40000,
		LotSizeUnits:              1,
		CapRate:                   float64(.7),
		AvgCap:                    float64(.6),
		BuildDate:                 dt,
		FLAGS:                     0,
		Ownership:                 0,
		TenantTradeName:           "Bill's Boar Emporium",
		LeaseGuarantor:            0,
		LeaseType:                 0,
		DeliveryDt:                dt,
		OriginalLeaseTerm:         int64(dur),
		RentCommencementDt:        dt,
		LeaseExpirationDt:         dt,
		TermRemainingOnLease:      int64(dur),
		TermRemainingOnLeaseUnits: int64(1),
		ROLID:                     0,
		Address:                   "1234 Elm Street",
		Address2:                  "",
		City:                      "Corn Bluff",
		State:                     "AK",
		PostalCode:                "98765",
		Country:                   "USA",
		LLResponsibilities:        "roof leaks",
		NOI:                       float64(30000),
		HQAddress:                 "1234 Elm Street",
		HQAddress2:                "",
		HQCity:                    "Corn Bluff",
		HQState:                   "AK",
		HQPostalCode:              "98765",
		HQCountry:                 "USA",
	}
	var delid, id int64
	if id, err = db.InsertProperty(ctx, &rs); err != nil {
		fmt.Printf("Error inserting Property: %s\n", err)
		os.Exit(1)
	}

	// Insert another for delete...
	if delid, err = db.InsertProperty(ctx, &rs); err != nil {
		fmt.Printf("Error inserting Property: %s\n", err)
		os.Exit(1)
	}
	if err = db.DeleteProperty(ctx, delid); err != nil {
		fmt.Printf("Error deleting Property: %s\n", err)
		os.Exit(1)
	}

	var rs1 db.Property
	if rs1, err = db.GetProperty(ctx, id); err != nil {
		fmt.Printf("error in GetProperty: %s\n", err.Error())
		os.Exit(1)
	}
	rs1.ROLID += 4
	if err = db.UpdateProperty(ctx, &rs1); err != nil {
		fmt.Printf("Error updating Property: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update Property\n")
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

// TestRentSteps checks the basic db functions for the RentSteps struct
//-----------------------------------------------------------------------------
func TestRentSteps(ctx context.Context) {
	var err error
	var rsl db.RentSteps
	var delid, id int64

	// create some rent steps
	rsl.RS = append(rsl.RS, db.RentStep{
		RSLID: 0,
		Dt:    time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		Opt:   "Year 1",
		Rent:  float64(2500),
		FLAGS: 0,
	})
	// create some rent steps
	rsl.RS = append(rsl.RS, db.RentStep{
		RSLID: 0,
		Dt:    time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		Opt:   "Year 1",
		Rent:  float64(2750),
		FLAGS: 0,
	})

	fmt.Printf("A. InsertRentSteps\n")
	if id, err = db.InsertRentSteps(ctx, &rsl); err != nil {
		fmt.Printf("Error inserting RentSteps: %s\n", err)
		os.Exit(1)
	}

	// Insert another for delete...
	fmt.Printf("B. InsertRentSteps\n")
	if delid, err = db.InsertRentSteps(ctx, &rsl); err != nil {
		fmt.Printf("Error inserting RentSteps: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("InsertRentSteps returned RSLID = %d\n", delid)

	// Read the rentsteps and make sure there are 2...
	fmt.Printf("B.1 GetRentSteps\n")
	var rs db.RentSteps
	rs, err = db.GetRentSteps(ctx, delid, true /*include items*/)
	if err != nil {
		fmt.Printf("Error getting RentSteps: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("GetRentSteps returned RSLID = %d, with %d items\n", rs.RSLID, len(rs.RS))
	for i := 0; i < len(rs.RS); i++ {
		fmt.Printf("\t%d. RSID = %d\n", i, rs.RS[i].RSID)
	}
	if len(rs.RS) != 2 {
		fmt.Printf("It should have loaded 2 RentStep items\n")
		os.Exit(1)
	}

	// Get the rent step items by themselves...
	fmt.Printf("B.2 GetRentStepsItems\n")
	var rsi []db.RentStep
	if rsi, err = db.GetRentStepsItems(ctx, rs.RSLID); err != nil {
		fmt.Printf("Error loading RentStep items for RSLID = %d: %s\n", rs.RSLID, err.Error())
		os.Exit(1)
	}
	if len(rsi) != 2 {
		fmt.Printf("It should have loaded 2 RentStep items\n")
		os.Exit(1)
	}

	fmt.Printf("C. DeleteRentSteps where RSLID = %d\n", delid)
	err = db.DeleteRentSteps(ctx, delid)
	fmt.Printf("Returned from DeleteRentSteps, RSLID=%d\n", delid)
	if err != nil {
		fmt.Printf("Error deleting RentSteps: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("D. GetRentSteps\n")
	var rsl1 db.RentSteps
	if rsl1, err = db.GetRentSteps(ctx, id, false); err != nil {
		fmt.Printf("error in GetRentSteps: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Got %d rentsteps\n", len(rsl1.RS))

	fmt.Printf("E. UpdateRentSteps\n")
	fmt.Printf("rsl1 = %#v\n", rsl1)
	if err = db.UpdateRentSteps(ctx, &rsl1); err != nil {
		fmt.Printf("Error updating RentSteps: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update RentSteps\n")
}

// TestTraffic checks the basic db functions for the Traffic struct
//-----------------------------------------------------------------------------
func TestTraffic(ctx context.Context) {
	var err error
	var t = db.Traffic{
		PRID:        1,
		Description: "Vehicles per day on Main street",
		Count:       725,
		FLAGS:       0,
	}

	if _, err = db.InsertTraffic(ctx, &t); err != nil {
		fmt.Printf("Error inserting Traffic: %s\n", err.Error())
		os.Exit(1)
	}

	var t1 = db.Traffic{
		PRID:        1,
		Description: "Elm Street",
		Count:       1400,
		FLAGS:       0,
	}
	if _, err = db.InsertTraffic(ctx, &t1); err != nil {
		fmt.Printf("Error inserting Traffic: %s\n", err.Error())
		os.Exit(1)
	}

	var a []db.Traffic
	if a, err = db.GetTrafficItems(ctx, 1); err != nil {
		fmt.Printf("Error getting Traffic items: %s\n", err.Error())
		os.Exit(1)
	}

	if len(a) != 2 {
		fmt.Printf("Error: wrong number of Traffic items. Expecting 2, got: %d\n", len(a))
		os.Exit(1)
	}

	fmt.Printf("Traffic tested: Success!\n")
}
