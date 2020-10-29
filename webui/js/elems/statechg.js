/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document,
    updatePropertyState, stateStatus, setTimeout
*/

"use strict";

function propertyStateChgOnLoad() {
    w2ui.propertyStateLayout.sizeTo('right', 450);
    w2ui.propertyStateLayout.load('right', '/static/html/statechg.html', null, setStateChangeDialogValues);
    w2ui.propertyStateLayout.show('right');
    setTimeout(setStateChangeDialogValues, 100 );
}

function setStateChangeDialogValues() {
    console.log("update the fields in the statechg dialog");
}

function closeStateChangeDialog() {
    w2ui.propertyStateLayout.hide('right');
}
