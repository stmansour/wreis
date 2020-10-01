/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document,
*/

"use strict";

// function buildPropertyStateUIElements() {
//     $().w2layout({
//         name: 'propertyStateFormLayout',
//         padding: 0,
//         panels: [
//             { type: 'left',    size: 0,     hidden: true },
//             { type: 'top',     size: 0,     hidden: true, content: 'top',  resizable: true, style: app.pstyle },
//             { type: 'main',    size: '60%', hidden: false, content: 'main', resizable: true, style: app.pstyle },
//             { type: 'preview', size: 0,     hidden: true,  content: 'PREVIEW'  },
//             { type: 'bottom',  size: 50,    hidden: false, content: 'bottom', resizable: false, style: app.pstyle },
//             { type: 'right',   size: 0,     hidden: true }
//         ]
//     });
// }



function propertyStateOnLoad() {
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

function setStateLabelColor(color,i) {
    var x = document.getElementsByName("stateLabelCell"+i);
    if (x == null) {return;}
    for (var i=0; i<x.length; i++) {
        x[i].style.color = color;
    }
}
