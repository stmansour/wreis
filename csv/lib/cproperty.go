package wcsv

import (
	"fmt"
	db "wreis/db/lib"
	util "wreis/util/lib"
)

// cproperty.go defines structures needed to parse a Property csv file

// PRName, et al, are indeces into the canonical PropertyList
//------------------------------------------------------------------------------
const (
	PRName                 = 0
	PRYearsInBusiness      = iota
	PRParentCompany        = iota
	PRURL                  = iota
	PRSymbol               = iota
	PRPrice                = iota
	PRDownPayment          = iota
	PRRentableArea         = iota
	PRRentableAreaUnits    = iota
	PRLotSize              = iota
	PRLotSizeUnits         = iota
	PRCapRate              = iota
	PRAvgCap               = iota
	PRBuildDate            = iota
	PROwnership            = iota
	PRTenantTradeName      = iota
	PRLeaseGuarantor       = iota
	PRLeaseType            = iota
	PRDeliveryDt           = iota
	PROriginalLeaseTerm    = iota
	PRLeaseCommencementDt  = iota
	PRLeaseExpirationDt    = iota
	PRTermRemainingOnLease = iota
	PRROID                 = iota
	PRAddress              = iota
	PRAddress2             = iota
	PRCity                 = iota
	PRState                = iota
	PRPostalCode           = iota
	PRCountry              = iota
	PRLLResponsibilities   = iota
	PRNOI                  = iota
	PRHQAddress            = iota
	PRHQAddress2           = iota
	PRHQCity               = iota
	PRHQState              = iota
	PRHQPostalCode         = iota
	PRHQCountry            = iota
	PRDriveThrough         = iota
	PRRoofStructResp       = iota
	PRFirstRightofRefusal  = iota
)

// CanonicalPropertyList defines the cannonical array of ColumnDefs for the CSV
// import of Properties.
//------------------------------------------------------------------------------
var CanonicalPropertyList = []ColumnDef{
	{Name: []string{"Name"}, Required: true, CaseSensitive: false, CanonicalIndex: PRName, Index: -1, FlagBit: 0},
	{Name: []string{"YearsInBusiness"}, Required: false, CaseSensitive: false, CanonicalIndex: PRYearsInBusiness, Index: -1, FlagBit: 0},
	{Name: []string{"ParentCompany"}, Required: false, CaseSensitive: false, CanonicalIndex: PRParentCompany, Index: -1, FlagBit: 0},
	{Name: []string{"URL"}, Required: false, CaseSensitive: false, CanonicalIndex: PRURL, Index: -1, FlagBit: 0},
	{Name: []string{"Symbol"}, Required: false, CaseSensitive: false, CanonicalIndex: PRSymbol, Index: -1, FlagBit: 0},
	{Name: []string{"Price"}, Required: false, CaseSensitive: false, CanonicalIndex: PRPrice, Index: -1, FlagBit: 0},
	{Name: []string{"DownPayment"}, Required: false, CaseSensitive: false, CanonicalIndex: PRDownPayment, Index: -1, FlagBit: 0},
	{Name: []string{"RentableArea"}, Required: false, CaseSensitive: false, CanonicalIndex: PRRentableArea, Index: -1, FlagBit: 0},
	{Name: []string{"RentableAreaUnits"}, Required: false, CaseSensitive: false, CanonicalIndex: PRRentableAreaUnits, Index: -1, FlagBit: 0},
	{Name: []string{"LotSize"}, Required: false, CaseSensitive: false, CanonicalIndex: PRLotSize, Index: -1, FlagBit: 0},
	{Name: []string{"LotSizeUnits"}, Required: false, CaseSensitive: false, CanonicalIndex: PRLotSizeUnits, Index: -1, FlagBit: 0},
	{Name: []string{"CapRate"}, Required: false, CaseSensitive: false, CanonicalIndex: PRCapRate, Index: -1, FlagBit: 0},
	{Name: []string{"AvgCap"}, Required: false, CaseSensitive: false, CanonicalIndex: PRAvgCap, Index: -1, FlagBit: 0},
	{Name: []string{"BuildDate"}, Required: false, CaseSensitive: false, CanonicalIndex: PRBuildDate, Index: -1, FlagBit: 0},
	{Name: []string{"Ownership"}, Required: false, CaseSensitive: false, CanonicalIndex: PROwnership, Index: -1, FlagBit: 0},
	{Name: []string{"TenantTradeName"}, Required: false, CaseSensitive: false, CanonicalIndex: PRTenantTradeName, Index: -1, FlagBit: 0},
	{Name: []string{"LeaseGuarantor"}, Required: false, CaseSensitive: false, CanonicalIndex: PRLeaseGuarantor, Index: -1, FlagBit: 0},
	{Name: []string{"LeaseType"}, Required: false, CaseSensitive: false, CanonicalIndex: PRLeaseType, Index: -1, FlagBit: 0},
	{Name: []string{"DeliveryDt"}, Required: false, CaseSensitive: false, CanonicalIndex: PRDeliveryDt, Index: -1, FlagBit: 0},
	{Name: []string{"OriginalLeaseTerm"}, Required: false, CaseSensitive: false, CanonicalIndex: PROriginalLeaseTerm, Index: -1, FlagBit: 0},
	{Name: []string{"LeaseCommencementDt"}, Required: false, CaseSensitive: false, CanonicalIndex: PRLeaseCommencementDt, Index: -1, FlagBit: 0},
	{Name: []string{"LeaseExpirationDt"}, Required: false, CaseSensitive: false, CanonicalIndex: PRLeaseExpirationDt, Index: -1, FlagBit: 0},
	{Name: []string{"TermRemainingOnLease"}, Required: false, CaseSensitive: false, CanonicalIndex: PRTermRemainingOnLease, Index: -1, FlagBit: 0},
	{Name: []string{"ROID"}, Required: false, CaseSensitive: false, CanonicalIndex: PRROID, Index: -1, FlagBit: 0},
	{Name: []string{"Address"}, Required: false, CaseSensitive: false, CanonicalIndex: PRAddress, Index: -1, FlagBit: 0},
	{Name: []string{"Address2"}, Required: false, CaseSensitive: false, CanonicalIndex: PRAddress2, Index: -1, FlagBit: 0},
	{Name: []string{"City"}, Required: false, CaseSensitive: false, CanonicalIndex: PRCity, Index: -1, FlagBit: 0},
	{Name: []string{"State"}, Required: false, CaseSensitive: false, CanonicalIndex: PRState, Index: -1, FlagBit: 0},
	{Name: []string{"PostalCode"}, Required: false, CaseSensitive: false, CanonicalIndex: PRPostalCode, Index: -1, FlagBit: 0},
	{Name: []string{"Country"}, Required: false, CaseSensitive: false, CanonicalIndex: PRCountry, Index: -1, FlagBit: 0},
	{Name: []string{"LLResponsibilities"}, Required: false, CaseSensitive: false, CanonicalIndex: PRLLResponsibilities, Index: -1, FlagBit: 0},
	{Name: []string{"NOI"}, Required: false, CaseSensitive: false, CanonicalIndex: PRNOI, Index: -1, FlagBit: 0},
	{Name: []string{"HQAddress"}, Required: false, CaseSensitive: false, CanonicalIndex: PRHQAddress, Index: -1, FlagBit: 0},
	{Name: []string{"HQAddress2"}, Required: false, CaseSensitive: false, CanonicalIndex: PRHQAddress2, Index: -1, FlagBit: 0},
	{Name: []string{"HQCity"}, Required: false, CaseSensitive: false, CanonicalIndex: PRHQCity, Index: -1, FlagBit: 0},
	{Name: []string{"HQState"}, Required: false, CaseSensitive: false, CanonicalIndex: PRHQState, Index: -1, FlagBit: 0},
	{Name: []string{"HQPostalCode"}, Required: false, CaseSensitive: false, CanonicalIndex: PRHQPostalCode, Index: -1, FlagBit: 0},
	{Name: []string{"HQCountry"}, Required: false, CaseSensitive: false, CanonicalIndex: PRHQCountry, Index: -1, FlagBit: 0},
	{Name: []string{"DriveThrough"}, Required: false, CaseSensitive: false, CanonicalIndex: PRDriveThrough, Index: -1, FlagBit: uint64(1 << 0)},
	{Name: []string{"RoofStructResp"}, Required: false, CaseSensitive: false, CanonicalIndex: PRRoofStructResp, Index: -1, FlagBit: uint64(1 << 1)},
	{Name: []string{"FirstRightofRefusal"}, Required: false, CaseSensitive: false, CanonicalIndex: PRFirstRightofRefusal, Index: -1, FlagBit: uint64(1 << 2)},
}

// PropertyHandler is called for each record of a Property csv file.
//
// INPUTS
// csvctx - context for this csv file, used to determine which column contains
//          what information.
// ss     - array of strings, one for each column in the csv file
// linno  - line number in the csvfile
//-----------------------------------------------------------------------------
func PropertyHandler(csvctx Context, ss []string, lineno int) []error {
	var p db.Property
	var errlist []error
	var elcount = int(0)
	var u uint64

	for i := 0; i < len(csvctx.Order); i++ {
		switch i {
		case PRName:
			p.Name = ss[csvctx.Order[PRName]]
		case PRYearsInBusiness:
			p.YearsInBusiness, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRParentCompany:
			p.ParentCompany = ss[csvctx.Order[i]]
		case PRURL:
			p.URL = ss[csvctx.Order[i]]
		case PRSymbol:
			p.Symbol = ss[csvctx.Order[i]]
		case PRPrice:
			p.Price, errlist = ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case PRDownPayment:
			p.DownPayment, errlist = ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case PRRentableArea:
			p.RentableArea, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRRentableAreaUnits:
			p.RentableAreaUnits, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRLotSize:
			p.LotSize, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRLotSizeUnits:
			p.LotSizeUnits, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRCapRate:
			p.CapRate, errlist = ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case PRAvgCap:
			p.AvgCap, errlist = ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case PRBuildDate:
			p.BuildDate, errlist = ParseDate(ss[csvctx.Order[i]], lineno, errlist)
		case PROwnership:
			p.Ownership, errlist = ParseInt(ss[csvctx.Order[i]], lineno, errlist)
		case PRTenantTradeName:
			p.TenantTradeName = ss[csvctx.Order[i]]
		case PRLeaseGuarantor:
			p.LeaseGuarantor, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRLeaseType:
			p.LeaseType, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRDeliveryDt:
			p.DeliveryDt, errlist = ParseDate(ss[csvctx.Order[i]], lineno, errlist)
		case PROriginalLeaseTerm:
			p.OriginalLeaseTerm, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRLeaseCommencementDt:
			p.LeaseCommencementDt, errlist = ParseDate(ss[csvctx.Order[i]], lineno, errlist)
		case PRLeaseExpirationDt:
			p.LeaseExpirationDt, errlist = ParseDate(ss[csvctx.Order[i]], lineno, errlist)
		case PRTermRemainingOnLease:
			p.TermRemainingOnLease, errlist = ParseInt64(ss[csvctx.Order[i]], lineno, errlist)
		case PRROID:
			p.ROID, errlist = ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case PRAddress:
			p.Address = ss[csvctx.Order[i]]
		case PRAddress2:
			p.Address2 = ss[csvctx.Order[i]]
		case PRCity:
			p.City = ss[csvctx.Order[i]]
		case PRState:
			p.State = ss[csvctx.Order[i]]
		case PRPostalCode:
			p.PostalCode = ss[csvctx.Order[i]]
		case PRCountry:
			p.Country = ss[csvctx.Order[i]]
		case PRLLResponsibilities:
			p.LLResponsibilities = ss[csvctx.Order[i]]
		case PRNOI:
			p.NOI, errlist = ParseFloat64(ss[csvctx.Order[i]], lineno, errlist)
		case PRHQAddress:
			p.HQAddress = ss[csvctx.Order[i]]
		case PRHQAddress2:
			p.HQAddress2 = ss[csvctx.Order[i]]
		case PRHQCity:
			p.HQCity = ss[csvctx.Order[i]]
		case PRHQState:
			p.HQState = ss[csvctx.Order[i]]
		case PRHQPostalCode:
			p.HQPostalCode = ss[csvctx.Order[i]]
		case PRHQCountry:
			p.HQCountry = ss[csvctx.Order[i]]
		case PRDriveThrough:
			u, errlist = GetBitFlagValue(ss[csvctx.Order[i]], 1<<0, errlist)
			p.FLAGS |= u
		case PRRoofStructResp:
			u, errlist = GetBitFlagValue(ss[csvctx.Order[i]], 1<<1, errlist)
			p.FLAGS |= u
		case PRFirstRightofRefusal:
			u, errlist = GetBitFlagValue(ss[csvctx.Order[i]], 1<<2, errlist)
			p.FLAGS |= u
		}
		if len(errlist) > elcount {
			errlist = append(errlist, fmt.Errorf("PropertyHandler: last error was detected on value for: %s = %s", CanonicalPropertyList[i].Name, ss[csvctx.Order[i]]))
			elcount = len(errlist)
			break
		}
	}

	util.Console("Line: %d p = %#v\n", lineno, p)

	return errlist
}

// ImportPropertyFile reads properties into the database from eht supplied file
//
// INPUTS
// fname = name of the csv file with Property data
//
// RETURNS
// a list of errors encountered. If there were no errors the length of the list
// will be 0.
//------------------------------------------------------------------------------
func ImportPropertyFile(fname string) []error {
	return ReadPropertyFile(fname, PropertyHandler)
}
