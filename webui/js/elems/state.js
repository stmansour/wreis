/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document,
*/

"use strict";

function propertyStateOnLoad() {
    if (propData.bStateLoaded) {
        updatePropertState();
        return
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
    var greenText = "#11AA11";
    var greenBG = "#e0ffe0";
    var grayText = "#888888";
    var grayBG = "#e0e0e0";
    var x;
    var color;
    var r = w2ui.propertyForm.record;
    console.log('propertyStateOnLoad has been called. FlowState = ' + r.FlowState);
    for (var i = 1; i <= propData.numStates; i++) {
        color = (r.FlowState >= i ) ? greenText : grayText;
        setStateColor('stateStepNo'+i,color);
        setStateColor('stateStepName'+i,color);
        color = (r.FlowState >= i ) ? greenBG : grayBG;
        setStateBGColor('stateLabelCell'+i,color);
        setStateBGColor('stateDataCell'+i,color);
        color = (r.FlowState >= i ) ? "black" : grayText;
        setStateLabelColor(color,i);
    }
    for (i = 0; i < propData.states.length; i++) {
        var s = "";
        var id;
        if (propData.states[i].InitiatorUID > 0) {
            s = propData.states[i].InitiatorName + '  (' + propData.states[i].InitiatorUID + ')';
            id = "stateCreateUser" + propData.states[i].FlowState;
            x = document.getElementById(id);
            if (x != null) {
                x.innerHTML = s;
            }
        }
        if (propData.states[i].ApproverUID > 0) {
            s = propData.states[i].ApproverName + '(' + propData.states[i].ApproverUID + ')';
            id = "stateApproveUser" + propData.states[i].FlowState;
            x = document.getElementById(id);
            if (x != null) {
                x.innerHTML = s;
            }
        }
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
