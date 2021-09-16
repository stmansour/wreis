/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document,
    updatePropertyState, stateStatus, setTimeout, setInnerHTML, setStateChangeDialogValues,
    loadPropertyForm, closePropertyForm,
*/

"use strict";

var stateChangeFormRedrawInProgress = false;

function newStateChangeRecord() {
    var rec = {
        ApproverName: {},
        OwnerName: {},
    };
    return rec;
}

function propertyStateChgOnLoad() {
    // first, clear out any associated values
    w2ui.stateChangeForm.record = newStateChangeRecord();  // clean it out, every time
    w2ui.propertyStateLayout.sizeTo('right', 450);
    w2ui.propertyStateLayout.show('right');
    // w2ui.propertyStateLayout.load('right', '/static/html/statechg.html', 1, setStateChangeDialogValues);
    w2ui.propertyStateLayout.content('right', w2ui.stateChangeForm);
    w2ui.propertyStateLayout.render();
    setTimeout(setStateChangeDialogValues, 150 );
}

function buildStateChangeForm() {
    $().w2form({
        name: 'stateChangeForm',
        // style: 'border: 0px; background-color: transparent;',
        // header: 'State Change Form',
        formURL: '/static/html/statechg.html',
        fields: [
            {name: 'ApproverName',   type: 'enum',       required: true,     html: {caption: "ApproverName"},
                options: {
                    url: '/v1/usertd/',
                    max: 1,
                    renderItem: function (item) {
                        return item.Name;
                    },
                    renderDrop: function (item) {
                        return item.Name;
                    },
                    compare: function (item, search) {
                        var s = item.Name.toLowerCase();
                        var srch = search.toLowerCase();
                        var match = (s.indexOf(srch) >= 0);
                        return match;
                    },
                    // onRemove: function(event) {
                    //     event.onComplete = function() {
                    //         w2ui.RAPeopleSearchForm.actions.reset();
                    //     };
                    // }
                }
            },
            {name: 'OwnerName',   type: 'enum',       required: true,     html: {caption: "OwnerName"},
                options: {
                    url: '/v1/usertd/',
                    max: 1,
                    renderItem: function (item) {
                        return item.Name;
                    },
                    renderDrop: function (item) {
                        return item.Name;
                    },
                    compare: function (item, search) {
                        var s = item.Name.toLowerCase();
                        var srch = search.toLowerCase();
                        var match = (s.indexOf(srch) >= 0);
                        console.log('' + match + ':  '+item.Name+' to '+srch);
                        return match;
                    },
                }
            },
        ],
    });

}

// setStateChangeDialogValues is called when the state change dialog has rendered.
// We need to make a few updates based on the current stateInfo
function setStateChangeDialogValues() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = 0;
    var s = "";

    if (stateChangeFormRedrawInProgress) {
        console.log("LOOP: stateChangeForm redraw in progress");
        return;
    }

    stateChangeFormRedrawInProgress = true;
    for (var i = 0; i < propData.states.length; i++) {
        // look for current flow state, not done
        if (propData.states[i].FlowState == FlowState && (propData.states[i].FLAGS & 0x4) == 0) {
            si = propData.states[i];
            break;
        }
    }
    if (typeof si === "number") {
        console.log('Could not determine the current stateInfo object');
        stateChangeFormRedrawInProgress = false;
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
        setInnerHTML("stateReadyButtonLbl","Back To<br>In-Progress");
    }

    var x = document.getElementById("stateReadyButton");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].OwnerUID);
        setButtonPadding(x);
    }
    x = document.getElementById("approveStateButton");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].ApproverUID);
        setButtonPadding(x);
    }
    x = document.getElementById("approveStateButton");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].ApproverUID);
        setButtonPadding(x);
    }
    x = document.getElementById("rejectStateButton");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].ApproverUID);
        setButtonPadding(x);
    }
    x = document.getElementById("smRejectReason");
    if (x != null) {
        x.disabled = (app.uid != propData.states[i].ApproverUID);
        setButtonPadding(x);
    }
    x = document.getElementById("revertStateButton");
    if (x != null) {
        setButtonPadding(x);
    }
    x = document.getElementById("btnSetNewApprover");
    if (x != null) {
        setButtonPadding(x);
    }
    x = document.getElementById("btnSetNewOwner");
    if (x != null) {
        setButtonPadding(x);
    }
    x = document.getElementById("btnDone");
    if (x != null) {
        setButtonPadding(x);
    }
    x = document.getElementById("terminateStateButton");
    if (x != null) {
        setButtonPadding(x);
    }

    stateChangeFormRedrawInProgress = false;
}

function setButtonPadding(x) {
    x.style.paddingLeft="4px";
    x.style.paddingRight="4px";
    x.style.paddingTop="4px";
    x.style.paddingBottom="4px";
    x.style.margin="4px";
}

function closeStateChangeDialog() {
    w2ui.propertyStateLayout.hide('right');
    //w2ui.propertyStateLayout.render();
}

function getCurrentStateInfo() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = 0;
    for (var i = 0; i < propData.states.length; i++) {
        // look for the one that matches current flow state and is NOT done
        if (propData.states[i].FlowState == FlowState && (propData.states[i].FLAGS & 0x4) == 0) {
            return propData.states[i];
        }
    }
    return null;
}

function stateReadyForApproval() {
    // Find the "in progress" record for the state selected...
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = getCurrentStateInfo();
    var cmd = "";
    if (si == null) {
        console.log('Could not determine the current stateInfo object');
        return;
    }
    if ((si.FLAGS & 0x2) == 0) {
        // currently marked as NOT READY, change to READY
        si.FLAGS |= 0x2;
        cmd = "ready";
    } else {
        // currently marked as ready, change to not ready
        si.FLAGS &= 0xeffffffffffffffd;
        cmd = "notready";
    }
    propData.bStateLoaded = false;
    finishStateChange(si,cmd);
}

// stateApproved calls the server with a request to approve the current state
//----------------------------------------------------------------------------
function stateApproved() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = getCurrentStateInfo();
    if (si == null) {
        console.log('Could not determine the current stateInfo object');
        return;
    }
    finishStateChange(si,"approve");
}

// stateRejected calls the server with a request to reject the current state
//----------------------------------------------------------------------------
function stateRejected() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = getCurrentStateInfo();
    if (si == null) {
        console.log('Could not determine the current stateInfo object');
        return;
    }

    si.Reason = "";
    var x = document.getElementById("smRejectReason");
    if (x != null) {
        si.Reason = x.value;
    }
    if (si.Reason.length < 2) {
        w2ui.stateChangeForm.error('ERROR: You must supply a reason');
        return;
    }
    finishStateChange(si,"reject");
}

// stateReverted calls the server with a request to revert the current state
// to the previous state.
//----------------------------------------------------------------------------
function stateReverted() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = getCurrentStateInfo();
    if (si == null) {
        console.log('Could not determine the current stateInfo object');
        return;
    }
    si.Reason = "";
    var x = document.getElementById("smRevertReason");
    if (x != null) {
        si.Reason = x.value;
    }
    if (si.Reason.length < 2) {
        w2ui.stateChangeForm.error('ERROR: You must supply a reason');
        return;
    }
    finishStateChange(si,"revert");
}

// stateTerminated calls the server with a request to terminate a property
//----------------------------------------------------------------------------
function stateTerminated() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = getCurrentStateInfo();
    if (si == null) {
        console.log('Could not determine the current stateInfo object');
        return;
    }

    si.Reason = "";
    var x = document.getElementById("smTerminateReason");
    if (x != null) {
        si.Reason = x.value;
    }
    if (si.Reason.length < 2) {
        w2ui.stateChangeForm.error('ERROR: You must supply a reason');
        return;
    }
    finishStateChange(si,"terminate");
    closePropertyForm();
}


// stateSetApprover calls the server with a request to change the approver to
// the newly selected user
//----------------------------------------------------------------------------
function stateSetApprover() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = getCurrentStateInfo();
    var uid = 0;

    if (si == null) {
        console.log('Could not determine the current stateInfo object');
        return;
    }
    if (typeof w2ui.stateChangeForm.record.ApproverName == "object" && w2ui.stateChangeForm.record.ApproverName != null) {
        if (w2ui.stateChangeForm.record.ApproverName.length > 0) {
            uid = w2ui.stateChangeForm.record.ApproverName[0].UID;
        }
    }
    if (uid == 0) {
        w2ui.stateChangeForm.error('ERROR: You must select a valid user');
        return;
    }
    si.ApproverUID = uid;
    si.Reason = "";
    var x = document.getElementById("smApproverReason");
    if (x != null) {
        si.Reason = x.value;
    }
    if (si.Reason.length < 2) {
        w2ui.stateChangeForm.error('ERROR: You must supply a reason');
        return;
    }
    finishStateChange(si,"setapprover");
}

// stateSetOwner calls the server with a request to change the approver to
// the newly selected user
//----------------------------------------------------------------------------
function stateSetOwner() {
    var FlowState = w2ui.propertyForm.record.FlowState;
    var si = getCurrentStateInfo();
    var uid = 0;

    if (si == null) {
        console.log('Could not determine the current stateInfo object');
        return;
    }
    if (typeof w2ui.stateChangeForm.record.OwnerName == "object" && w2ui.stateChangeForm.record.OwnerName != null) {
        if (w2ui.stateChangeForm.record.OwnerName.length > 0) {
            uid = w2ui.stateChangeForm.record.OwnerName[0].UID;
        }
    }
    if (uid == 0) {
        w2ui.stateChangeForm.error('ERROR: You must select a valid user');
        return;
    }
    si.OwnerUID = uid;
    si.Reason = "";
    var x = document.getElementById("smOwnerReason");
    if (x != null) {
        si.Reason = x.value;
    }
    if (si.Reason.length < 2) {
        w2ui.stateChangeForm.error('ERROR: You must supply a reason');
        return;
    }
    finishStateChange(si,"setowner");
}

// finishStateChange performs the repetitive tasks for a state update.
//
// INPUTS:
// si = stateinfo object to update
// c = command name
//---------------------------------------------------------------------------
function finishStateChange(si,c) {
    //-----------------------------------------------------------------------
    // if the current state is in-progress, change to READY, and vice-versa
    //-----------------------------------------------------------------------
    var params = {
        cmd: c,
        records: [si],
    };
    var dat = JSON.stringify(params);
    var url = '/v1/stateinfo/' + w2ui.propertyForm.record.PRID;
    // console.log('url = ' + url);
    // console.log('data = ' + dat);
    // console.log(" ");

    return $.post(url, dat, null, "json")
    .done(function(data) {
        if (data.status === "error") {
            w2ui.stateChangeForm.error('ERROR: '+ data.message);
            return;
        }
        var prid = propData.PRID;
        loadPropertyForm(prid);
    })
    .fail(function(data){
        var err = JSON.parse(data.responseText);
        w2ui.stateChangeForm.error("Update failed: " + err.message);
    });

}
