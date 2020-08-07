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
		PRID:                 0,
		Name:                 "Bill's Boar Emporium",
		YearsInBusiness:      8,
		ParentCompany:        "",
		URL:                  "http://bbb.com/",
		Symbol:               "BBE",
		Price:                float64(12345.67),
		DownPayment:          float64(40000),
		RentableArea:         30000,
		RentableAreaUnits:    1,
		LotSize:              40000,
		LotSizeUnits:         1,
		CapRate:              float64(.7),
		AvgCap:               float64(.6),
		BuildDate:            dt,
		FLAGS:                0,
		Ownership:            0,
		TenantTradeName:      "Bill's Boar Emporium",
		LeaseGuarantor:       0,
		LeaseType:            0,
		DeliveryDt:           dt,
		OriginalLeaseTerm:    int64(dur),
		LeaseCommencementDt:  dt,
		LeaseExpirationDt:    dt,
		TermRemainingOnLease: int64(dur),
		ROLID:                0,
		Address:              "1234 Elm Street",
		Address2:             "",
		City:                 "Corn Bluff",
		State:                "AK",
		PostalCode:           "98765",
		Country:              "USA",
		LLResponsibilities:   "roof leaks",
		NOI:                  float64(30000),
		HQAddress:            "1234 Elm Street",
		HQAddress2:           "",
		HQCity:               "Corn Bluff",
		HQState:              "AK",
		HQPostalCode:         "98765",
		HQCountry:            "USA",
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
		Count: 1,
		Dt:    time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		Opt:   1,
		Rent:  float64(2500),
		FLAGS: 0,
	})
	// create some rent steps
	rsl.RS = append(rsl.RS, db.RentStep{
		RSLID: 0,
		Count: 2,
		Dt:    time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		Opt:   1,
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

	fmt.Printf("C. DeleteRentSteps where RSLID = %d\n", delid)
	// if err = db.DeleteRentSteps(ctx, delid); err != nil {
	// 	fmt.Printf("Error deleting RentSteps: %s\n", err)
	// 	os.Exit(1)
	// }
	err = db.DeleteRentSteps(ctx, delid)
	fmt.Printf("Returned from DeleteRentSteps, RSLID=%d\n", delid)
	if err != nil {
		fmt.Printf("Error deleting RentSteps: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("D. GetRentSteps\n")
	var rsl1 db.RentSteps
	if rsl1, err = db.GetRentSteps(ctx, id); err != nil {
		fmt.Printf("error in GetRentSteps: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Got %d rentsteps\n", len(rsl1.RS))

	fmt.Printf("E. UpdateRentSteps\n")
	fmt.Printf("rsl1 = %#v\n", rsl1)
	if err = db.UpdateRentSteps(ctx, &rsl1); err != nil {
		fmt.Printf("Error updating RentSteps: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update RentSteps\n")
}
