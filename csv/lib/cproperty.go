package wcsv

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
)

// CanonicalPropertyList defines the cannonical array of ColumnDefs for the CSV
// import of Properties.
//------------------------------------------------------------------------------
var CanonicalPropertyList = []ColumnDef{
	{Name: []string{"Name"}, Required: true, CaseSensitive: false, CanonicalIndex: PRName, Index: -1},
	{Name: []string{"YearsInBusiness"}, Required: true, CaseSensitive: false, CanonicalIndex: PRYearsInBusiness, Index: -1},
	{Name: []string{"ParentCompany"}, Required: true, CaseSensitive: false, CanonicalIndex: PRParentCompany, Index: -1},
	{Name: []string{"URL"}, Required: true, CaseSensitive: false, CanonicalIndex: PRURL, Index: -1},
	{Name: []string{"Symbol"}, Required: true, CaseSensitive: false, CanonicalIndex: PRSymbol, Index: -1},
	{Name: []string{"Price"}, Required: true, CaseSensitive: false, CanonicalIndex: PRPrice, Index: -1},
	{Name: []string{"DownPayment"}, Required: true, CaseSensitive: false, CanonicalIndex: PRDownPayment, Index: -1},
	{Name: []string{"RentableArea"}, Required: true, CaseSensitive: false, CanonicalIndex: PRRentableArea, Index: -1},
	{Name: []string{"RentableAreaUnits"}, Required: true, CaseSensitive: false, CanonicalIndex: PRRentableAreaUnits, Index: -1},
	{Name: []string{"LotSize"}, Required: true, CaseSensitive: false, CanonicalIndex: PRLotSize, Index: -1},
	{Name: []string{"LotSizeUnits"}, Required: true, CaseSensitive: false, CanonicalIndex: PRLotSizeUnits, Index: -1},
	{Name: []string{"CapRate"}, Required: true, CaseSensitive: false, CanonicalIndex: PRCapRate, Index: -1},
	{Name: []string{"AvgCap"}, Required: true, CaseSensitive: false, CanonicalIndex: PRAvgCap, Index: -1},
	{Name: []string{"BuildDate"}, Required: true, CaseSensitive: false, CanonicalIndex: PRBuildDate, Index: -1},
	{Name: []string{"Ownership"}, Required: true, CaseSensitive: false, CanonicalIndex: PROwnership, Index: -1},
	{Name: []string{"TenantTradeName"}, Required: true, CaseSensitive: false, CanonicalIndex: PRTenantTradeName, Index: -1},
	{Name: []string{"LeaseGuarantor"}, Required: true, CaseSensitive: false, CanonicalIndex: PRLeaseGuarantor, Index: -1},
	{Name: []string{"LeaseType"}, Required: true, CaseSensitive: false, CanonicalIndex: PRLeaseType, Index: -1},
	{Name: []string{"DeliveryDt"}, Required: true, CaseSensitive: false, CanonicalIndex: PRDeliveryDt, Index: -1},
	{Name: []string{"OriginalLeaseTerm"}, Required: true, CaseSensitive: false, CanonicalIndex: PROriginalLeaseTerm, Index: -1},
	{Name: []string{"LeaseCommencementDt"}, Required: true, CaseSensitive: false, CanonicalIndex: PRLeaseCommencementDt, Index: -1},
	{Name: []string{"LeaseExpirationDt"}, Required: true, CaseSensitive: false, CanonicalIndex: PRLeaseExpirationDt, Index: -1},
	{Name: []string{"TermRemainingOnLease"}, Required: true, CaseSensitive: false, CanonicalIndex: PRTermRemainingOnLease, Index: -1},
	{Name: []string{"ROID"}, Required: true, CaseSensitive: false, CanonicalIndex: PRROID, Index: -1},
	{Name: []string{"Address"}, Required: true, CaseSensitive: false, CanonicalIndex: PRAddress, Index: -1},
	{Name: []string{"Address2"}, Required: true, CaseSensitive: false, CanonicalIndex: PRAddress2, Index: -1},
	{Name: []string{"City"}, Required: true, CaseSensitive: false, CanonicalIndex: PRCity, Index: -1},
	{Name: []string{"State"}, Required: true, CaseSensitive: false, CanonicalIndex: PRState, Index: -1},
	{Name: []string{"PostalCode"}, Required: true, CaseSensitive: false, CanonicalIndex: PRPostalCode, Index: -1},
	{Name: []string{"Country"}, Required: true, CaseSensitive: false, CanonicalIndex: PRCountry, Index: -1},
	{Name: []string{"LLResponsibilities"}, Required: true, CaseSensitive: false, CanonicalIndex: PRLLResponsibilities, Index: -1},
	{Name: []string{"NOI"}, Required: true, CaseSensitive: false, CanonicalIndex: PRNOI, Index: -1},
	{Name: []string{"HQAddress"}, Required: true, CaseSensitive: false, CanonicalIndex: PRHQAddress, Index: -1},
	{Name: []string{"HQAddress2"}, Required: true, CaseSensitive: false, CanonicalIndex: PRHQAddress2, Index: -1},
	{Name: []string{"HQCity"}, Required: true, CaseSensitive: false, CanonicalIndex: PRHQCity, Index: -1},
	{Name: []string{"HQState"}, Required: true, CaseSensitive: false, CanonicalIndex: PRHQState, Index: -1},
	{Name: []string{"HQPostalCode"}, Required: true, CaseSensitive: false, CanonicalIndex: PRHQPostalCode, Index: -1},
	{Name: []string{"HQCountry"}, Required: true, CaseSensitive: false, CanonicalIndex: PRHQCountry, Index: -1},
}
