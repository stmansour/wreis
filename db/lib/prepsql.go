package db

import (
	"strings"
)

// June 3, 2016 -- As the params change, it's easy to forget to update all the statements with the correct
// field names and the proper number of replacement characters.  I'm starting a convention where the SELECT
// fields are set into a variable and used on all the SELECT statements for that table.  The fields and
// replacement variables for INSERT and UPDATE are derived from the SELECT string.

var mySQLRpl = string("?")
var myRpl = mySQLRpl

// TRNSfields defined fields for Transactant, used in at least one other function
var TRNSfields = string("TCID,BID,NLID,FirstName,MiddleName,LastName,PreferredName,CompanyName,IsCompany,PrimaryEmail,SecondaryEmail,WorkPhone,CellPhone,Address,Address2,City,State,PostalCode,Country,Website,Comment,FLAGS,CreateTS,CreateBy,LastModTime,LastModBy")

// GenSQLInsertAndUpdateStrings generates a string suitable for SQL INSERT and UPDATE statements given the fields as used in SELECT statements.
//
//  example:
//	given this string:      "LID,BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModTime,LastModBy"
//  we return these five strings:
//  1)  "BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModBy"                 -- use for SELECT
//  2)  "?,?,?,?,?,?,?,?"  														-- use for INSERT
//  3)  "BID=?RAID=?,GLNumber=?,Status=?,Type=?,Name=?,AcctType=?,LastModBy=?"  -- use for UPDATE
//  4)  "LID,BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModBy", 			-- use for INSERT (no PRIMARYKEY), add "WHERE LID=?"
//  5)  "?,?,?,?,?,?,?,?,?"  													-- use for INSERT (no PRIMARYKEY)
//
// Note that in this convention, we remove LastModTime from insert and update statements (the db is set up to update them by default) and
// we remove the initial ID as that number is AUTOINCREMENT on INSERTs and is not updated on UPDATE.
func GenSQLInsertAndUpdateStrings(s string) (string, string, string, string, string) {
	fields := strings.Split(s, ",")

	// mostly 0th element is ID, but it is not necessary
	s0 := fields[0]
	s2 := fields[1:] // skip the ID

	insertFields := []string{} // fields which are allowed while INSERT
	updateFields := []string{} // fields which are allowed while while UPDATE

	// remove fields which value automatically handled by database while insert and update op.
	for _, fld := range s2 {
		fld = strings.TrimSpace(fld)
		if fld == "" { // if nothing then continue
			continue
		}
		// INSERT FIELDS Inclusion
		if fld != "LastModTime" && fld != "CreateTS" { // remove these fields for INSERT
			insertFields = append(insertFields, fld)
		}
		// UPDATE FIELDS Inclusion
		if fld != "LastModTime" && fld != "CreateTS" && fld != "CreateBy" { // remove these fields for UPDATE
			updateFields = append(updateFields, fld)
		}
	}

	var s3, s4 string
	for i := range insertFields {
		if i == len(insertFields)-1 {
			s3 += myRpl
		} else {
			s3 += myRpl + ","
		}
	}

	for i, uFld := range updateFields {
		if i == len(updateFields)-1 {
			s4 += uFld + "=" + myRpl
		} else {
			s4 += uFld + "=" + myRpl + ","
		}
	}

	// list down insert fields with comma separation
	s = strings.Join(insertFields, ",")

	s5 := s0 + "," + s     // for INSERT where first val is not AUTOINCREMENT
	s6 := s3 + "," + myRpl // for INSERT where first val is not AUTOINCREMENT
	return s, s3, s4, s5, s6
}

// BuildPreparedStatements is where we build the DBFields map and create the
// prepared sql statements for queries
//
// INPUTS
//
// RETURNS
//
//------------------------------------------------------------------------------
func BuildPreparedStatements() {
	var err error
	var s1, s2, s3, flds string

	//==========================================
	// Property
	//==========================================
	flds = "PRID,Name,YearsInBusiness,ParentCompany,URL,Symbol,Price,DownPayment,RentableArea,RentableAreaUnits,LotSize,LotSizeUnits,CapRate,AvgCap,BuildDate,FLAGS,Ownership,TenantTradeName,LeaseGuarantor,LeaseType,DeliveryDt,OriginalLeaseTerm,LeaseCommencementDt,LeaseExpirationDt,TermRemainingOnLease,ROLID,RSLID,Address,Address2,City,State,PostalCode,Country,LLResponsibilities,NOI,HQAddress,HQAddress2,HQCity,HQState,HQPostalCode,HQCountry,CreateTS,CreateBy,LastModTime,LastModBy"
	Wdb.DBFields["Property"] = flds
	Wdb.Prepstmt.GetProperty, err = Wdb.DB.Prepare("SELECT " + flds + " FROM Property where PRID=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	Wdb.Prepstmt.InsertProperty, err = Wdb.DB.Prepare("INSERT INTO Property (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	Wdb.Prepstmt.UpdateProperty, err = Wdb.DB.Prepare("UPDATE Property SET " + s3 + " WHERE PRID=?")
	Errcheck(err)
	Wdb.Prepstmt.DeleteProperty, err = Wdb.DB.Prepare("DELETE from Property WHERE PRID=?")
	Errcheck(err)

	//==========================================
	// Rent Step
	//==========================================
	flds = "RSID,RSLID,Count,Dt,Rent,CreateTS,CreateBy,LastModTime,LastModBy"
	Wdb.DBFields["RentStep"] = flds
	Wdb.Prepstmt.GetRentStep, err = Wdb.DB.Prepare("SELECT " + flds + " FROM RentStep where RSID=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	Wdb.Prepstmt.InsertRentStep, err = Wdb.DB.Prepare("INSERT INTO RentStep (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	Wdb.Prepstmt.UpdateRentStep, err = Wdb.DB.Prepare("UPDATE RentStep SET " + s3 + " WHERE RSID=?")
	Errcheck(err)
	Wdb.Prepstmt.DeleteRentStep, err = Wdb.DB.Prepare("DELETE from RentStep WHERE RSID=?")
	Errcheck(err)

	//==========================================
	// Rent Steps
	//==========================================
	flds = "RSLID,FLAGS,CreateTS,CreateBy,LastModTime,LastModBy"
	Wdb.DBFields["RentSteps"] = flds
	Wdb.Prepstmt.GetRentSteps, err = Wdb.DB.Prepare("SELECT " + flds + " FROM RentSteps where RSLID=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	Wdb.Prepstmt.InsertRentSteps, err = Wdb.DB.Prepare("INSERT INTO RentSteps (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	Wdb.Prepstmt.UpdateRentSteps, err = Wdb.DB.Prepare("UPDATE RentSteps SET " + s3 + " WHERE RSLID=?")
	Errcheck(err)
	Wdb.Prepstmt.DeleteRentSteps, err = Wdb.DB.Prepare("DELETE from RentSteps WHERE RSLID=?")
	Errcheck(err)

	//==========================================
	// Renew Option
	//==========================================
	flds = "ROID,ROLID,Count,Dt,Rent,CreateTS,CreateBy,LastModTime,LastModBy"
	Wdb.DBFields["RenewOption"] = flds
	Wdb.Prepstmt.GetRenewOption, err = Wdb.DB.Prepare("SELECT " + flds + " FROM RenewOption where ROID=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	Wdb.Prepstmt.InsertRenewOption, err = Wdb.DB.Prepare("INSERT INTO RenewOption (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	Wdb.Prepstmt.UpdateRenewOption, err = Wdb.DB.Prepare("UPDATE RenewOption SET " + s3 + " WHERE ROID=?")
	Errcheck(err)
	Wdb.Prepstmt.DeleteRenewOption, err = Wdb.DB.Prepare("DELETE from RenewOption WHERE ROID=?")
	Errcheck(err)

	//==========================================
	// Renew Options
	//==========================================
	flds = "ROLID,FLAGS,CreateTS,CreateBy,LastModTime,LastModBy"
	Wdb.DBFields["RenewOptions"] = flds
	Wdb.Prepstmt.GetRenewOptions, err = Wdb.DB.Prepare("SELECT " + flds + " FROM RenewOptions where ROLID=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	Wdb.Prepstmt.InsertRenewOptions, err = Wdb.DB.Prepare("INSERT INTO RenewOptions (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	Wdb.Prepstmt.UpdateRenewOptions, err = Wdb.DB.Prepare("UPDATE RenewOptions SET " + s3 + " WHERE ROLID=?")
	Errcheck(err)
	Wdb.Prepstmt.DeleteRenewOptions, err = Wdb.DB.Prepare("DELETE from RenewOptions WHERE ROLID=?")
	Errcheck(err)

}
