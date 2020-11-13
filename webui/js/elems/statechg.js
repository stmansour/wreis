/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document,
    updatePropertyState, stateStatus, setTimeout, setInnerHTML, setStateChangeDialogValues,
    loadPropertyForm,
*/

"use strict";

function propertyStateChgOnLoad() {
    w2ui.propertyStateLayout.sizeTo('right', 450);
    w2ui.propertyStateLayout.show('right');
    w2ui.propertyStateLayout.load('right', '/static/html/statechg.html', 1, setStateChangeDialogValues);
    setTimeout(setStateChangeDialogValues, 1000 );
}


// setStateChangeDialogValues is called when the state change dialog has rendered.
// We need to make a few updates based on the current stateInfo
function setStateChangeDialogValues() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = 0;
    var s = "";

    for (var i = 0; i < propData.states.length; i++) {
        // look for current flow state, not done
        if (propData.states[i].FlowState == FlowState && (propData.states[i].FLAGS & 0x4) == 0) {
            si = propData.states[i];
            break;
        }
    }
    if (typeof si === "number") {
        console.log('Could not determine the current stateInfo object');
        return;
    }

    //-------------------------------------------------------------------------
    // only work on this
    // if the flags currently show this stateInfo to be READY, then the button
    // should indicate that we set it back to IN PROGRESS.  If the flags show
    // that it to be IN PROGRESS, then the button should allow us to change
    // it to READY...
    //-------------------------------------------------------------------------
    if ((si.FLAGS & 0x2) == 0) {
        setInnerHTML("stateReadyLabel","IN PROGRESS");
        setInnerHTML("stateReadyButtonLbl","Ready For Approval");
    } else {
        setInnerHTML("stateReadyLabel","READY");
        setInnerHTML("stateReadyButtonLbl","Back To In-Progress");
    }

    var x = document.getElementById("stateReadyButton");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].OwnerUID);
    }
    x = document.getElementById("approveStateButton");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].ApproverUID);
    }
    x = document.getElementById("approveStateButton");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].ApproverUID);
    }
    x = document.getElementById("rejectStateButton");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].ApproverUID);
    }


}

function closeStateChangeDialog() {
    w2ui.propertyStateLayout.hide('right');
}

function stateReadyForApproval() {
    // Find the "in progress" record for the state selected...
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = 0;
    for (var i = 0; i < propData.states.length; i++) {
        // look for the one that matches current flow state and is NOT done
        if (propData.states[i].FlowState == FlowState && (propData.states[i].FLAGS & 0x4) == 0) {
            si = propData.states[i];
            break;
        }
    }
    if (typeof si === "number") {
        console.log('Could not determine the current stateInfo object');
        return;
    }

    //-----------------------------------------------------------------------
    // if the current state is in-progress, change to READY, and vice-versa
    //-----------------------------------------------------------------------
    var params = {
        cmd: "ready",
        records: [si],
    };
    if ((si.FLAGS & 0x2) != 0) {
        params.cmd = "notready";
    }
    var dat = JSON.stringify(params);
    var url = '/v1/stateinfo/' + w2ui.propertyForm.record.PRID;

    return $.post(url, dat, null, "json")
    .done(function(data) {
        if (data.status === "error") {
            w2ui.propertyGrid.error('ERROR: '+ data.message);
            return;
        }
        var prid = propData.PRID;
        loadPropertyForm(prid);
    })
    .fail(function(data){
        var err = JSON.parse(data.responseText);
        w2ui.propertyGrid.error("Update failed: " + err.message);
    });
}
