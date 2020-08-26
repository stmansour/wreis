"use strict";
"esversion 6";
/*global
  app, document, setDateControl, dateMonthBack, getDateFromDT, getTimeFromDT, dateFromString,
  dateFmtStr, zeroPad, yearFwd, dateYearFwd, yearBack, dateYearBack, applyLocaltimeDateOffset,
  UTCDateStringToW2UIValidDate,stringToDate,
*/

// getDropDownSelectedIndex returns the selected index of a dropdown menu with
// with the supplied id.
//
// RETURNS:
// 0 - n -> the dropdown for id was found and the index is the value returned
// -1    -> the element with id could not be found or was not a select element
//------------------------------------------------------------------------------
function getDropDownSelectedIndex(id) {
    var x = document.getElementById(id);
    if (x == null || typeof x.selectedIndex == "undefined") {
        return -1;
    }
    return x.selectedIndex;
}

// setDropDownSelectedIndex sets the selection of the supplied dropdown menu
// to the supplied value. If val < 0 it returns -1. If there is a problem with
// the id element, the return value is -1.
//
// RETURNS:
// 0  -> success
// -1 -> element with id could not be found or was not a select element, or the
//       supplied index val was < 0.
//------------------------------------------------------------------------------
function setDropDownSelectedIndex(id,val) {
    if (val < 0) {
        return -1;
    }
    var x = document.getElementById(id);
    if (x == null || typeof x.selectedIndex == "undefined") {
        return -1;
    }
    x.selectedIndex = val;
    return 0;
}
