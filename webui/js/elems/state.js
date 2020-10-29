/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document,
    updatePropertyState, stateStatus,closeStateChangeDialog,
*/

"use strict";

function buildStateUIElements() {

    $().w2layout({
        name: 'propertyStateLayout',
        padding: 0,
        panels: [
            { type: 'left',    size: 0,     hidden: true },
            { type: 'top',     size: 0,     hidden: true,  content: 'top',  resizable: true, style: app.pstyle },
            { type: 'main',    size: '60%', hidden: false, content: 'main', resizable: true, style: app.pstyle },
            { type: 'preview', size: 0,     hidden: true,  content: 'PREVIEW'  },
            { type: 'bottom',  size: 0,     hidden: true,  content: 'bottom', resizable: false, style: app.pstyle },
            { type: 'right',   size: 0,     hidden: true,  content: 'right',
                toolbar: {
                    items: [
                        { id: 'btnNotes', type: 'button', icon: 'fa fa-sticky-note-o' },
                        { id: 'bt3', type: 'spacer' },
                        { id: 'btnClose', type: 'button', icon: 'fa fa-times' },
                    ],
                    onClick: function (event) {
                        if (event.target == 'btnClose') {
                            closeStateChangeDialog();
                        }
                    },
                },
            },
        ],
    });
}



function propertyStateOnLoad() {
    if (propData.bStateLoaded) {
        updatePropertyState();
        return;
    }

    if (propData.PRID == 0) {
        propData.bStateLoaded = true;
        updatePropertyState();
        return;
    }

    var params = {
        cmd: "get",
    };
    var dat = JSON.stringify(params);
    var url = '/v1/stateinfo/' + propData.PRID;

    return $.post(url, dat, null, "json")
    .done(function(data) {
        if (data.status === "error") {
            w2ui.propertyGrid.error('ERROR: '+ data.message);
            return;
        }
        propData.states = data.records;
        propData.bStateLoaded = true;
        updatePropertyState();
    })
    .fail(function(data){
            w2ui.propertyGrid.error("Get states failed. " + data);
    });

}

function updatePropertyState() {
    var x;
    var color;
    var r = w2ui.propertyForm.record;
    var fs;
    if (r == null) {
        console.log('r is null.  w2ui.propertyForm.record =  ' + w2ui.propertyForm.record);
        return;
    }
    fs = r.FlowState;
    if (propData.states != null) {
        for (var i = 0; i < propData.states.length; i++) {
            var s = "";
            var id;
            var dt;
            var j = propData.states[i].FlowState;

            color = getStateTextColor(j,fs,0);
            setStateColor('stateStepNo'+j,color);
            setStateColor('stateStepName'+j,color);

            color = getStateTextColor(j,fs,1);
            setStateBGColor('stateLabelCell'+j,color);
            setStateBGColor('stateDataCell'+j,color);

            color = (r.FlowState >= j ) ? "black" : propData.notStartedText;
            setStateLabelColor(color,j);

            if (propData.states[i].InitiatorUID > 0) {
                dt = new Date(propData.states[i].InitiatorDt);
                s = propData.states[i].InitiatorName + ', ' + dt.toDateString();
                id = "stateCreateUser" + j;
                setHTMLByID(id,s);
            }
            if (propData.states[i].ApproverUID > 0) {
                s = propData.states[i].ApproverName;
                id = "stateApproveUser" + j;
                setHTMLByID(id,s);
            }
            stateStatus(propData.states[i],r.FlowState);
            id = "stateLastMod" + j;
            dt = new Date(propData.states[i].LastModTime);
            s = propData.states[i].LastModByName + ", " + dt.toDateString();
            setHTMLByID(id,s);
            if (propData.states[i].FlowState == w2ui.propertyForm.record.FlowState) {
                setStateChange(w2ui.propertyForm.record.FlowState);
            }
        }
    }
}

// getStateTextColor describes the status of the state.
//
// INPUTS
//   ts  = FlowState of stateinfo structure
//   fs = FlowState of the current property
//   g  = 0 -> foreground color, 1 -> background color
//
// RETURNS
//   the requested color string
//------------------------------------------------------------------------------
function getStateTextColor(ts,fs,g) {
    if ( ts < fs ) {
        // state is completed
        return g != 0 ? propData.doneBG : propData.doneText;
    } else if ( ts == fs ) {
        // state is in progress
        return g != 0 ? propData.inProgressBG : propData.inProgressText;
    } else {
        return g != 0 ? propData.notStartedBG : propData.notStartedText;
    }
    return "black";
}

// stateStatus describes the status of the state.
//
//   FLAGS & 0x1 is the approval status.
//                0 =>  approved
//                1 =>  rejected  and t.Reason explains why
//
// INPUTS
//   t  = stateinfo structure
//   fs = FlowState of the current property
//
// RETURNS
//
//------------------------------------------------------------------------------
function stateStatus(t,fs) {
    // If ApproverDt field has year >1970 then the approver has made a determination
    var dt = new Date(t.ApproverDt);
    var label;
    var id = "stateStatus" + t.FlowState;  // the label for this particular state

    if (dt.getFullYear() > 1970) {
        if (dt.FLAGS & 0x1 > 0) {
            // 1 means not approved
            label = "Rejected: " + dt.toDateString() + ", " + t.Reason;
        } else {
            label = "Approved: " + dt.toDateString();
        }
    } else {
        label = "";
        if (fs == t.FlowState) {
            label = "In Progress";
        }
    }

    setHTMLByID(id,label);
}

function setHTMLByID(id,s) {
    var x = document.getElementById(id);
    if (x != null) {
        x.innerHTML = s;
    }
}

function setStateColor(id,color) {
    var x = document.getElementById(id);
    if (x != null) {
        x.style.color = color;
    }
}

function setStateBGColor(id,color) {
    var x = document.getElementById(id);
    if (x != null) {
        x.style.backgroundColor = color;
    }
}

function setStateLabelColor(color,j) {
    var x = document.getElementsByName("stateLabelCell"+j);
    if (x == null) {return;}
    for (var i=0; i<x.length; i++) {
        x[i].style.color = color;
    }
}

// INPUTS
//   x = state number
function setStateChange(y) {
    //var s = `<button class="w2ui-btn" onclick="w2popup.load({url:'/static/html/statechg.html',showMax:true})">Change...</button>`;
    var s = `<br><button class="w2ui-btn" onclick="propertyStateChgOnLoad();">Change...</button>`;
    var x = document.getElementById("stateChange"+y);
    if (x != null) {
        x.innerHTML = s;
    }
}
