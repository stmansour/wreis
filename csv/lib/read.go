package wcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	util "wreis/util/lib"
)

// csv is a library that enables a very forgiving way to import csv files.
// We define a canonical list of columns and attributes, we then match the
// supplied csv file with the attributes as well as possible. If we have
// enough of the required columns then each line of the csv file is read,
// parsed, and the handler is called with an array of strings containing
// the values from the csv file record just read in the canonical order.
//------------------------------------------------------------------------------

// ColumnDef defines a column of the supplied csv file so that it can be handdled
// efficiently. Each column is compared to the ColumnDefs in the canonical list.
// If it is found, the index of the expected position (from the []CanonicalIndex
//  is stored in the Index value.
//------------------------------------------------------------------------------
type ColumnDef struct {
	Name           []string // list of acceptable column-heading names for this value
	Required       bool     // is this value required: false = no, true = yes
	CaseSensitive  bool     // is the column heading match case-sensitive: false = no
	CanonicalIndex int      // column index in the canonical definition
	Index          int      // index where this
	FlagBit        uint64   // 0 means not a flag bit, otherwise 1<<n where n is the bit number
}

// Context is the context structure for a specific csv file
//------------------------------------------------------------------------------
type Context struct {
	Handler    HandlerFunc // handler to call for each line
	ColumnDefs []ColumnDef // array of line definitions
	Order      []int       // []index that enables referencing the supplied csv values in an expected order
	csvfile    *os.File    // the csv file to parse
	reader     *csv.Reader // csv file reader
}

// HandlerFunc is the prototype of a csv handler function. The CSV subsystem
// will read a requested file, parse each line, and call a HandlerFunc to do
// something with the data.
// INPUTS
// []string = values found in this row
// int      = line number of csv file
//
// RETURNS
// nothing at this time
//------------------------------------------------------------------------------
type HandlerFunc func(Context, []string, int) []error

// nameMatch compares the supplied name to the names listed in the columndef.
// It returns true if there is a match, false otherwise.
//
// INPUTS
// ns - name listed in the csv column heading. May contain whitespace, which
//      will be removed
// ss - the column heading we're looking to match this to.
//
// RETURNS
// bool - true if it's a match, false otherwise
//------------------------------------------------------------------------------
func nameMatch(n string, def *ColumnDef) bool {
	count := len(def.Name)

	if len(n) == 0 {
		return false
	}
	for i := 0; i < count; i++ {
		var match bool
		if def.CaseSensitive {
			match = n == def.Name[i]
		} else {
			match = strings.EqualFold(def.Name[i], n)
		}
		if match {
			return true
		}
	}
	return false
}

// GetContext is the context structure for a specific csv file
//------------------------------------------------------------------------------
func GetContext(fname string, h HandlerFunc) (Context, error) {
	var ctx Context
	var err error
	var cols, rec []string

	//----------------------------------
	// make sure we can open the file
	//----------------------------------
	ctx.csvfile, err = os.Open(fname)
	if err != nil {
		return ctx, err
	}
	ctx.Handler = h

	//-------------------------------------------------------------------------
	// Read the first line of the file. There's no way to tell if the first
	// line of the csv file is data or column headers. So, we just require it
	// to be column headers
	//-------------------------------------------------------------------------
	ctx.ColumnDefs = deepCopy(CanonicalPropertyList)
	// util.Console("len ctx.ColumnDefs = %d\n", len(ctx.ColumnDefs))
	ctx.reader = csv.NewReader(ctx.csvfile)
	if cols, err = ctx.reader.Read(); err != nil {
		return ctx, err
	}

	//--------------------------------------------------
	// Strip the whitespace from the column titles...
	//--------------------------------------------------
	var pcount = len(cols)
	ctx.Order = make([]int, pcount)
	for i := 0; i < pcount; i++ {
		rec = append(rec, util.Stripchars(cols[i], " \r\n\t"))
		ctx.Order[i] = -1 // initialize to indicate no mapping for this column
	}
	//-------------------------------------------------------------------------
	// rec now has the column headers of the file we're looking at. Spin through
	// each column heading and match it to the canonical file column.
	//-------------------------------------------------------------------------

	//------------------------------------------
	// spin through the canonical list in order
	//------------------------------------------
	var count = len(CanonicalPropertyList)
	for i := 0; i < count; i++ {
		//------------------------------------------------------------------
		// search for a match to ctx.ColumnDefs[i].Name in rec...
		// The csv file need not be ordered the same as the CanonicalPropertyList
		// That is, the value for CanonicalPropertyList[i] can be found at
		// rec[ ctx.ColumnDefs[i].Index ]
		//------------------------------------------------------------------
		ctx.ColumnDefs[i].Index = -1 // assume we don't find it
		for j := 0; j < pcount; j++ {
			if nameMatch(rec[j], &ctx.ColumnDefs[i]) {
				rec[j] = ""                 // don't need to look for this any further
				ctx.ColumnDefs[i].Index = j // col j of the csv matches CanonicalIndex i
				ctx.Order[i] = j            // for quick indexing
				break                       // no need to loop further
			}
		}
		if ctx.ColumnDefs[i].Index < 0 {
			// util.Console("Didn't find it!\n")
		}
		//-------------------------------------------------------------------
		// If this CanonicalIndex was not matched AND it is required, then
		// return now with an error.
		//-------------------------------------------------------------------
		if ctx.ColumnDefs[i].Required && ctx.ColumnDefs[i].Index < 0 {
			return ctx, fmt.Errorf("Required column %s was not found", ctx.ColumnDefs[i].Name[0])
		}
	}

	dbg := false
	if dbg {
		for i := 0; i < len(ctx.Order); i++ {
			util.Console("%s is in col %d\n", CanonicalPropertyList[i].Name, ctx.Order[i])
		}
	}
	return ctx, nil
}

// ReadPropertyFile reads the csvfile line by line and handles each line.
//
// INPUTS
// fname = file to parse
// h     = handler function to be called with an array of strings containing
//         all the column values for a line in the csv file. It will be called
//         once for each data line of the csv file.
//
// RETURNS
// errlist = slice of all errors encountered
//------------------------------------------------------------------------------
func ReadPropertyFile(fname string, h HandlerFunc) []error {
	var record []string
	var errlist []error
	ctx, err := GetContext(fname, h)
	if err != nil {
		errlist = append(errlist, err)
		return errlist
	}
	line := 2
	for {
		// Read each record from csv
		record, err = ctx.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errlist = append(errlist, err)
			return errlist
		}
		errlist = ctx.Handler(ctx, record, line)
		if len(errlist) > 0 {
			break
		}
		line++
	}
	return errlist
}
