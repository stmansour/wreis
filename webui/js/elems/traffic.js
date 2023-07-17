/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise,
    setPropertyFormActionButtons,
*/

"use strict";

var TrafficModule = {
    id: 0,
};

function getNextTrafficID() {
    TrafficModule.id -= 1;
    return TrafficModule.id;
}

function newTrafficRecord() {
    var id = getNextTrafficID();
    var t = {
        recid: id,
        PRID: w2ui.propertyForm.record.PRID,
        TID: id,
        Description: "",
        Count: 0,
        FLAGS: 0,
    };
    return t;
}

function buildTrafficUIElements() {
    $().w2grid({
        name: 'propertyTrafficGrid',
        url: '/v1/trafficitems',
        method: 'POST',
        postData: {
            cmd: 'get',
        },
        show: {
            toolbar         : true,
            footer          : false,
            toolbarAdd      : true,   // indicates if toolbar add new button is visible
            toolbarDelete   : false,   // indicates if toolbar delete button is visible
            toolbarSave     : false,   // indicates if toolbar save button is visible
            selectColumn    : false,
            expandColumn    : false,
            toolbarEdit     : false,
            toolbarSearch   : false,
            toolbarInput    : false,
            searchAll       : false,
            toolbarReload   : false,
            toolbarColumns  : false,
        },
        //======================================================================
        // FLAGS
        //     1<<0  Drive Through?  0 = no, 1 = yes
        //	   1<<1  Roof & Structure Responsibility: 0 = Tenant, 1 = Landlord
        //	   1<<2  Right Of First Refusal: 0 = no, 1 = yes
        //======================================================================
        columns: [
            {field: 'recid',       text: 'recid',       size: '60px',  sortable: true, hidden: true},
            {field: 'TID',         text: 'TID',         size: '60px',  sortable: true, hidden: true},
            {field: 'PRID',        text: 'PRID',        size: '60px',  sortable: true, hidden: true},
            {field: 'FLAGS',       text: 'FLAGS',       size: '60px',  sortable: true, hidden: true},
            {field: 'Count',       text: 'Count',       size: '80px',  sortable: true, hidden: false},
            {field: 'Description', text: 'Description', size: '350px', sortable: true, hidden: false},
            {field: 'CreateTime',  text: 'CreateTime',  size: '60px',  sortable: true, hidden: true},
            {field: 'CreateBy',    text: 'CreateBy',    size: '60px',  sortable: true, hidden: true},
            {field: 'LastModTime', text: 'LastModTime', size: '60px',  sortable: true, hidden: true},
            {field: 'LastModBy',   text: 'LastModBy',   size: '60px',  sortable: true, hidden: true},
        ],
         onLoad: function(event) {
            event.onComplete = function() {
                propData.bTrafficLoaded = true;
                w2ui.propertyTrafficGrid.url = ''; // don't go back to the server until we're ready to save
                for (var i = 0; i < w2ui.propertyTrafficGrid.records.length; i++) {
                    w2ui.propertyTrafficGrid.records[i].recid = w2ui.propertyTrafficGrid.records[i].TID;
                }
            };
        },
        onAdd: function(event) {
            w2ui.propertyTrafficForm.record = newTrafficRecord();
            w2ui.propertyTrafficGrid.add(w2ui.propertyTrafficForm.record);
            showTrafficForm();
        },
        onClick: function(event) {
            event.onComplete = function(event) {
                var r = w2ui.propertyTrafficForm.record;
                var x = this.getSelection();
                if (x.length < 1) {
                    return;
                }
                var idx       = this.get(x[0],true); // get the index of the selection
                var fr        = w2ui.propertyTrafficGrid.records[idx];
                r.PRID        = fr.PRID;
                r.TID         = fr.TID;
                r.Description = fr.Description;
                r.Count       = fr.Count;
                r.FLAGS       = fr.FLAGS;
                showTrafficForm();
            };
        }
    });

    $().w2form({
        name: 'propertyTrafficForm',
        style: 'border: 0px; background-color: transparent;',
        // header: 'Property Detail',
        formURL: '/static/html/formTraffic.html',
        // url: '/v1/property',
        fields: [
            {field: 'PRID',        text: 'PRID',        type: 'int',   required: false},
            {field: 'TID',         text: 'TID',         type: 'int',   required: false},
            {field: 'Description', text: 'Description', type: 'text',  required: true},
            {field: 'Count',       text: 'Count',       type: 'int',   required: true },
            {field: 'CreateTime',  text: 'CreateTime',  type: 'text',  required: false},
            {field: 'CreateBy',    text: 'CreateBy',    type: 'text',  required: false},
            {field: 'LastModTime', text: 'LastModTime', type: 'text',  required: false},
            {field: 'LastModBy',   text: 'LastModBy',   type: 'text',  required: false},
        ],
        toolbar: {
            items: [
                { id: 'btnNotes', type: 'button', icon: 'fa fa-sticky-note-o' },
                { id: 'bt3', type: 'spacer' },
                { id: 'btnClose', type: 'button', icon: 'fa fa-times' },
            ],
            onClick: function (event) {
                if (event.target == 'btnClose') {
                    closeTrafficForm();
                    w2ui.propertyTrafficGrid.render();
                }
            },
        },
        onRefresh: function(event) {
            event.onComplete = function(event) {
            };
        },
    });


    $().w2layout({
        name: 'propertyTrafficLayout',
        padding: 0,
        panels: [
            { type: 'left',    size: 0,     hidden: true },
            { type: 'top',     size: 0,     hidden: true,  content: 'top',  resizable: true, style: app.pstyle },
            { type: 'main',    size: '60%', hidden: false, content: 'main', resizable: true, style: app.pstyle },
            { type: 'preview', size: 0,     hidden: true,  content: 'PREVIEW'  },
            { type: 'bottom',  size: 0,     hidden: true,  content: 'bottom', resizable: false, style: app.pstyle },
            { type: 'right',   size: 0,     hidden: true,  content: 'right' }
        ],
    });
}

// TrafficDelete is designed to be called by the UI when a Traffic item
// is deleted. It does not save this change to the database, only to the grid.
//------------------------------------------------------------------------------
function TrafficDelete() {
    var r = w2ui.propertyTrafficForm.record;
    var g = w2ui.propertyTrafficGrid;
    var i = g.get(r.recid,true);
    if (i >= 0) {
        var removed = g.records.splice(i,1);
        // console.log('removed = ' + removed);
    }
    closeTrafficForm();
    g.render();
}

// TrafficSave is designed to be called by the UI when changes to a Traffic item
// are completed. It does not save to the database, only to the grid.
//------------------------------------------------------------------------------
function TrafficSave() {
    var r = w2ui.propertyTrafficForm.record;
    var g = w2ui.propertyTrafficGrid;

    var x=w2ui.propertyTrafficForm.validate(true);
    if (x.length > 0) {
        return;
    }
    g.set(r.recid,r);
    closeTrafficForm();
    g.render();
}

function showTrafficForm() {
    w2ui.propertyTrafficLayout.html('right',w2ui.propertyTrafficForm);
    w2ui.propertyTrafficLayout.sizeTo('right',400);
    w2ui.propertyTrafficLayout.show('right',true);
    setPropertyFormActionButtons(false);
}

function closeTrafficForm() {
    w2ui.propertyTrafficLayout.hide('right',true);
    setPropertyFormActionButtons(true);
}

function saveTraffic() {
    //-----------------------------------------------------------------------
    // If we never loaded the Traffics, then they weren't changed, so just
    // return success.
    //-----------------------------------------------------------------------
    if (!propData.bTrafficLoaded) {
        return new Promise( function(resolve,reject) {
            if (true) {
                resolve("success");
            } else {
                reject("error");
            }
        });
    }

    // Just a precaution...
    for (var i = 0; i < w2ui.propertyTrafficGrid.records.length; i++) {
        w2ui.propertyTrafficGrid.records[i].PRID = w2ui.propertyForm.record.PRID;
    }

    //-----------------------------------------------------------------------
    // We have loaded the Traffics, so we need to go through the save...
    //-----------------------------------------------------------------------
    var params = {
        cmd: "save",
        records: w2ui.propertyTrafficGrid.records
    };
    var dat = JSON.stringify(params);
    var url = '/v1/trafficitems/' + w2ui.propertyForm.record.PRID;

    return $.post(url, dat, null, "json")
    .done(function(data) {
        // if (data.status === "success") {
        // }
        if (data.status === "error") {
            w2ui.propertyGrid.error('ERROR: '+ data.message);
        }
    })
    .fail(function(data){
            w2ui.propertyGrid.error("Save Traffic failed. " + data);
    });
}
