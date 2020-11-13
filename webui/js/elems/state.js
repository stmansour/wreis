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
        var f = w2ui.propertyFormBtns;
        var curState = 1;
        var s = "<p><table>";
        var j = 1;
        var MINYEAR = 2000;
        var notAnApproval = false;
        var statusSet = false;

        //--------------------------------------------------------
        // turn off the Save buttons for this record
        //--------------------------------------------------------
        $(f.box).find("button[name=save]").prop( "disabled", true );
        $(f.box).find("button[name=saveadd]").prop( "disabled", true );
        $(f.box).find("button[name=delete]").prop( "disabled", true );

        for (var i = 0; i < propData.states.length; i++) {
            var id;
            var dt;
            var flags = propData.states[i].FLAGS;
            j = propData.states[i].FlowState;


            //---------------------------------------------------------------
            // Before processing any further, dump the info we have if the
            // FlowState has changed
            //---------------------------------------------------------------
            if (j != curState) {
                s += "</table>";
                setHTMLByID("stateDataCell" + curState,s);
                s = "<p><table>";
                curState = j;
            }

            color = getStateTextColor(j,fs,0);
            setStateColor('stateStepNo'+j,color);
            setStateColor('stateStepName'+j,color);

            color = getStateTextColor(j,fs,1);
            setStateBGColor('stateLabelCell'+j,color);
            setStateBGColor('stateDataCell'+j,color);

            color = (r.FlowState >= j ) ? "black" : propData.notStartedText;
            setStateLabelColor(color,j);
            /*
            **  FLAGS
            **          bit  Description
            **          ---  ----------------------------------------------------------------------
                0x1      0  valid only when ApproverUID > 0, 0 = State Approved, 1 = not approved
                0x2      1  0 = work is in progress, 1 = READY: request approval for this state
                0x4      2  0 = this state is work in progress, 1 = work is concluded on this StateInfo
                0x8      3  0 = this state has not been reverted.  1 = this state was reverted
                0x10     4  0 = no owner change, 1 = owner change -- changer will be the UID of LastModBy on this StateInfo, and creator of the StateInfo with new owner
                0x10     4  0 = no owner change, 1 = owner change -- changer will be the UID of LastModBy on this StateInfo, and creator of the StateInfo with new owner
            */
            notAnApproval = false;  // assume it's an Approval
            statusSet = false;
            s += '<tr><td align="right">Status:</td><td>';
            if ((flags & 0x2) != 0) {
                s += 'READY <span style="color:#3333DD;font-weight:bold;">(approval pending)</span>';
                statusSet = true;
            }
            if ((flags & 0x8) != 0) {
                s += "REVERTED ";
                notAnApproval = true;  // this would make it not an approval
                statusSet = true;
            }
            if ((flags & 0x10) != 0) {
                s += "OWNER CHANGED ";
                notAnApproval = true;  // this would make it not an approval
                statusSet = true;
            }
            if ((flags & 0x20) != 0) {
                s += "APPROVER CHANGED ";
                notAnApproval = true;  // this would make it not an approval
                statusSet = true;
            }
            dt = new Date(propData.states[i].ApproverDt);
            if (!notAnApproval && propData.states[i].ApproverUID > 0 && dt.getFullYear() > MINYEAR) {
                s += ((flags & 0x1) != 0) ? "NOT " : "";
                s += "APPROVED ";
                statusSet = true;
            }
            if ((flags & 0x4) != 0) {
                s += '<span style="color:#117711;font-weight:bold;">&#10004;</span>';
                statusSet = true;
            }

            if (!statusSet) {
                s += '<span style="color:#117711;font-weight:bold;">IN PROGRESS</span>';
            }
            s += '</td></tr>';

            if (propData.states[i].Reason.length > 0) {
                s += '<tr><td align="right">Reason:</td><td>' + propData.states[i].Reason + '</td></tr>';
            }

            if (propData.states[i].ApproverUID > 0) {
                dt = new Date(propData.states[i].ApproverDt);
                var y = "";
                if (dt.getFullYear() > MINYEAR) {
                    y = dt.toDateString();
                }
                s += '<tr><td align="right">Approver</td><td>' + propData.states[i].ApproverName;
                if (dt.getFullYear() > MINYEAR) {
                    s += ', ' + dt.toDateString() + "</td></tr>";
                }
            }
            if (propData.states[i].OwnerUID > 0) {
                dt = new Date(propData.states[i].OwnerDt);
                s += '<tr><td align="right">Owner:</td><td>' + propData.states[i].OwnerName;
                if (dt.getFullYear() > MINYEAR ) {
                    s += ', ' + dt.toDateString() + "</td></tr>";
                }
            }
            stateStatus(propData.states[i],r.FlowState);
            s += "<tr><td>Last Update:</td><td>";
            dt = new Date(propData.states[i].LastModTime);
            s += propData.states[i].LastModByName + ", " + dt.toDateString();

            //-----------------------------------------------------------
            // now add a spacer line
            //-----------------------------------------------------------
            s += '<tr><td colspan="2" height="10"></td>';

            //-----------------------------------------------------------
            // ADD CHANGE BUTTON TO CURRENT STATE...
            //-----------------------------------------------------------
            if (propData.states[i].FlowState == w2ui.propertyForm.record.FlowState) {
                setStateChange(w2ui.propertyForm.record.FlowState);
            }
        }
        s += "</table>";
        setHTMLByID("stateDataCell" + curState,s);
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
