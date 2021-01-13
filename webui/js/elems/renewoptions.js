/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise,
*/

"use strict";

var RenewOptionModule = {
    id: 0,
};

function getNextRenewOptionID() {
    RenewOptionModule.id -= 1;
    return RenewOptionModule.id;
}

function newRenewOptionRecord() {
    var id = getNextRenewOptionID();
    var rs = {
        recid: id,
        ROLID: 0,
        ROID: id,
        Opt: 0,
        Dt: new Date(),
        Rent: 0,
        FLAGS: propData.roType,
    };
    return rs;
}

function buildRenewOptionsUIElements() {
    $().w2grid({
        name: 'propertyRenewOptionsGrid',
        url: '/v1/renewoptions',
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
            {field: 'recid',       caption: 'recid',      size: '60px', sortable: true, hidden: true},
            {field: 'ROID',        caption: 'ROID',       size: '60px', sortable: true, hidden: true},
            {field: 'ROLID',       caption: 'ROLID',      size: '60px', sortable: true, hidden: true},
            {field: 'FLAGS',       caption: 'FLAGS',      size: '60px', sortable: true, hidden: true},
            {field: 'Opt',         caption: 'Period',     size: '250px', sortable: true, hidden: false},
            {field: 'Dt',          caption: 'Date',       size: '80px', sotrable: true, hidden: false},
            {field: 'Rent',        caption: 'Rent',       size: '80px', sortable: true, hidden: false, render: 'money'},
            {field: 'CreateTime',  caption: 'CreateTime', size: '60px', sortable: true, hidden: true},
            {field: 'CreateBy',   caption: 'CreateBy',  size: '60px', sortable: true, hidden: true},
            {field: 'LastModTime', caption: 'LastModTime',size: '60px', sortable: true, hidden: true},
            {field: 'LastModBy',   caption: 'LastModBy',  size: '60px', sortable: true, hidden: true},
        ],
        onLoad: function(event) {
            event.onComplete = function() {
                propData.bRenewOptionsLoaded = true;
                w2ui.propertyRenewOptionsGrid.url = ''; // don't go back to the server until we're ready to save
                for (var i = 0; i < w2ui.propertyRenewOptionsGrid.records.length; i++) {
                    w2ui.propertyRenewOptionsGrid.records[i].recid = w2ui.propertyRenewOptionsGrid.records[i].ROID;
                }
                SetRenewOptionColumns(propData.roType);  // since all records are the same in BIT 0, just look at first
            };
        },
        onAdd: function(event) {
            w2ui.propertyRenewOptionForm.record = newRenewOptionRecord();
            var ev = {
                type: "click",
                target: ((w2ui.propertyRenewOptionForm.record.FLAGS & 1) == 0) ? "roListType:roListOpt" : "roListType:roListDate",
            };
            RenewOptionTypeChange(ev); // make sure the opt vs date mode is set correctly
            showRenewOptionForm();
        },
        onClick: function(event) {
            event.onComplete = function(event) {
                var r = w2ui.propertyRenewOptionForm.record;
                var x = this.getSelection();
                if (x.length < 1) {return;}
                var idx = this.get(x[0],true); // get the index of the selection
                var fr = w2ui.propertyRenewOptionsGrid.records[idx];
                r.ROLID = fr.ROLID;
                r.ROID = fr.ROID;
                r.Opt = fr.Opt;
                r.Dt = fr.Dt;
                r.Rent = fr.Rent;
                r.FLAGS = fr.FLAGS;
                r.recid = fr.ROID;
                showRenewOptionForm();
            };
        }
    });

    $().w2form({
        name: 'propertyRenewOptionForm',
        style: 'border: 0px; background-color: transparent;',
        // header: 'Property Detail',
        formURL: '/static/html/formRenewOption.html',
        // url: '/v1/property',
        fields: [
            {field: 'ROLID',       type: 'int',   required: false},
            {field: 'ROID',        type: 'int',   required: false},
            {field: 'Opt',         type: 'text',  required: false},
            {field: 'Dt',          type: 'date',  required: false},
            {field: 'Rent',        type: 'money', required: true, render: 'money' },
            {field: 'CreateTime',  type: 'text',  required: false},
            {field: 'CreateBy',   type: 'text',  required: false},
            {field: 'LastModTime', type: 'text',  required: false},
            {field: 'LastModBy',   type: 'text',  required: false},
        ],
        toolbar: {
            items: [
                { id: 'btnNotes', type: 'button', icon: 'fa fa-sticky-note-o' },
                { id: 'bt3', type: 'spacer' },
                { id: 'btnClose', type: 'button', icon: 'fa fa-times' },
            ],
            onClick: function (event) {
                if (event.target == 'btnClose') {
                    w2ui.renewOptionsLayout.hide('right',true);
                    w2ui.propertyRenewOptionsGrid.render();
                }
            },
        },
        onRefresh: function(event) {
            event.onComplete = function(event) {
                EnableRenewOptionFormFields();
            };
        },
    });


    $().w2layout({
        name: 'renewOptionsLayout',
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

function finishRenewOptionsGridToolbar() {
    var t = w2ui.propertyRenewOptionsGrid.toolbar;
    t.add([
        { id: 'bt3', type: 'spacer' },
        { id: 'roListType', type: 'menu-radio', icon: 'fa fa-star',
            text: function (item) {
                var text = item.selected;
                var el   = this.get('roListType:' + item.selected);
                return 'Renew Type: ' + el.text;
            },
            selected: 'roListOpt',
            items: [
                { id: 'roListOpt', text: 'Period', icon: 'fa fa-tachometer' },
                { id: 'roListDate', text: 'Date', icon: 'fa fa-tachometer' },
            ]
        },
    ]);
    t.on('*', RenewOptionTypeChange);
}

function RenewOptionDelete() {
    var r = w2ui.propertyRenewOptionForm.record;
    var g = w2ui.propertyRenewOptionsGrid;
    var i = g.get(r.recid,true);
    if (typeof i == "number" && i >= 0) {
        var removed = g.records.splice(i,1);
        // console.log('removed = ' + removed);
    }
    w2ui.renewOptionsLayout.hide('right',true);
    g.render();
}

function RenewOptionSave() {
    var r = w2ui.propertyRenewOptionForm.record;
    var g = w2ui.propertyRenewOptionsGrid;

    var x=w2ui.propertyRenewOptionForm.validate(true);
    if (x.length > 0) {
        return;
    }


    if (r.ROID < 0) {
        w2ui.propertyRenewOptionsGrid.add(r);
        w2ui.propertyRenewOptionForm.record.ROID = 0;  // to make sure that this one won't be added again
    }
    g.set(r.recid,r);

    w2ui.renewOptionsLayout.hide('right',true);
    g.render();
}

function RenewOptionTypeChange(event) {
    if (event.type != "click") {
        return;
    }
    //console.log('EVENT: '+ event.type + ' TARGET: '+ event.target, event);
    switch (event.target) {
    case "roListType:roListOpt":
        SetRenewOptionColumns(0);
        SetRenewOptionFLAGs(0);
        propData.roType = 0;
        break;
    case "roListType:roListDate":
        SetRenewOptionColumns(1);
        SetRenewOptionFLAGs(1);
        propData.roType = 1;
        break;
    }
}

function SetRenewOptionColumns(FLAGS) {
    var b = FLAGS & 0x1;
    var t = w2ui.propertyRenewOptionsGrid.toolbar.get("roListType");
    if (b == 0) {
        w2ui.propertyRenewOptionsGrid.hideColumn("Dt");
        w2ui.propertyRenewOptionsGrid.showColumn("Opt");
        t.selected = "roListOpt";
    } else {
        w2ui.propertyRenewOptionsGrid.hideColumn("Opt");
        w2ui.propertyRenewOptionsGrid.showColumn("Dt");
        t.selected = "roListDate";
    }
    w2ui.propertyRenewOptionsGrid.toolbar.refresh();
}

// INPUTS
//   FLAGS =  0 means use Options
//            1 means use Dates
//-----------------------------------------------------------------------
function SetRenewOptionFLAGs(FLAGS) {
    for (var i = 0; i < w2ui.propertyRenewOptionsGrid.records.length; i++) {
        w2ui.propertyRenewOptionsGrid.records[i].FLAGS &= 0xeffffffffffffffe;
        w2ui.propertyRenewOptionsGrid.records[i].FLAGS |= FLAGS;
    }
}

function EnableRenewOptionFormFields() {
    if (1 == propData.roType) {
        $('#Opt').prop('disabled', true);
        $('#Dt').prop('disabled', false);
    } else {
        $('#Opt').prop('disabled', false);
        $('#Dt').prop('disabled', true);
    }
}

function showRenewOptionForm() {
    w2ui.renewOptionsLayout.content('right',w2ui.propertyRenewOptionForm);
    w2ui.renewOptionsLayout.sizeTo('right',400);
    w2ui.renewOptionsLayout.show('right',true);
}


function saveRenewOptions() {
    //-----------------------------------------------------------------------
    // If we never loaded the renewoptions, then they weren't changed, so just
    // return success.
    //-----------------------------------------------------------------------
    if (!propData.bRenewOptionsLoaded) {
        return new Promise( function(resolve,reject) {
            if (true) {
                resolve("success");
            } else {
                reject("error");
            }
        });
    }

    //-----------------------------------------------------------------------
    // We have loaded the renewoptions, so we need to go through the save...
    //-----------------------------------------------------------------------
    var params = {
        cmd: "save",
        PRID: w2ui.propertyForm.record.PRID,
        records: w2ui.propertyRenewOptionsGrid.records
    };
    for (var i = 0; i < params.records.length; i++) {
        var d = new Date(params.records[i].Dt);
        params.records[i].Dt = d.toUTCString();
        if (params.records[i].Opt == "0") {
            RenewOptionModule.maxOpt++;
            params.records[i].Opt = '' + RenewOptionModule.maxOpt;
        }
    }
    var dat = JSON.stringify(params);
    var url = '/v1/renewoptions/' + w2ui.propertyForm.record.ROLID;

    return $.post(url, dat, null, "json")
    .done(function(data) {
        if (data.status === "error") {
            w2ui.propertyGrid.error('ERROR: '+ data.message);
        }
        propData.bRenewOptionsLoaded = false;
    })
    .fail(function(data){
            w2ui.propertyGrid.error("Save RenewOptions failed. " + data);
    });
}
