/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, w2utils, setInnerHTML,
    number_format,
*/

"use strict";

var RentStepModule = {
    id: 0,
};

function getNextRentStepID() {
    RentStepModule.id -= 1;
    return RentStepModule.id;
}

function newRentStepRecord() {
    var id = getNextRentStepID();
    var rs = {
        recid: id,
        RSLID: 0,
        RSID: id,
        Opt: 0,
        Dt: new Date(),
        Rent: 0,
        FLAGS: propData.rsType,
    };
    return rs;
}

function buildRentStepsUIElements() {
    $().w2grid({
        name: 'propertyRentStepsGrid',
        url: '/v1/rentsteps',
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
            {field: 'recid',  caption: 'recid',        size: '60px',  sortable: true, hidden: true},
            {field: 'RSID',   caption: 'RSID',         size: '60px',  sortable: true, hidden: true},
            {field: 'RSLID',  caption: 'RSLID',        size: '60px',  sortable: true, hidden: true},
            {field: 'FLAGS',  caption: 'FLAGS',        size: '60px',  sortable: true, hidden: true},
            {field: 'Opt',    caption: 'Period',       size: '250px', sortable: true, hidden: false},
            {field: 'Dt',     caption: 'Date',         size: '80px',  sotrable: true, hidden: false},
            {field: 'Rent',   caption: 'Annual Rent',  size: '110px', sortable: true, hidden: false, render: 'money'},
            {field: 'mrent',  caption: 'Monthly Rent', size: '110px', sortable: false, style: 'text-align: right', hidden: false,
                render: function (record, index, col_index) {
                    var y = record.Rent/12;
                    var s = "$" + number_format(y,2,'.',',');
                    return s;
                }
            },
            {field: 'CreateTime',  caption: 'CreateTime',   size: '60px',  sortable: true, hidden: true},
            {field: 'CreateBy',    caption: 'CreateBy',     size: '60px',  sortable: true, hidden: true},
            {field: 'LastModTime', caption: 'LastModTime',  size: '60px',  sortable: true, hidden: true},
            {field: 'LastModBy',   caption: 'LastModBy',    size: '60px',  sortable: true, hidden: true},
        ],
        onLoad: function(event) {
            event.onComplete = function() {
                propData.bRentStepsLoaded = true;
                w2ui.propertyRentStepsGrid.url = ''; // don't go back to the server until we're ready to save
                for (var i = 0; i < w2ui.propertyRentStepsGrid.records.length; i++) {
                    w2ui.propertyRentStepsGrid.records[i].recid = w2ui.propertyRentStepsGrid.records[i].RSID;
                }
                SetRentStepColumns(propData.rsType);  // since all records are the same in BIT 0, just look at first
            };
        },
        onAdd: function(event) {
            w2ui.propertyRentStepForm.record = newRentStepRecord();
            var ev = {
                type: "click",
                target: ((w2ui.propertyRentStepForm.record.FLAGS & 1) == 0) ? "rsListType:rsListOpt" : "rsListType:rsListDate",
            };
            RentStepTypeChange(ev); // make sure the opt vs date mode is set correctly
            showRentStepForm();
        },
        onClick: function(event) {
            event.onComplete = function(event) {
                var r = w2ui.propertyRentStepForm.record;
                var x = this.getSelection();
                if (x.length < 1) {return;}
                var idx = this.get(x[0],true); // get the index of the selection
                var fr = w2ui.propertyRentStepsGrid.records[idx];
                r.RSLID = fr.RSLID;
                r.RSID = fr.RSID;
                r.Opt = fr.Opt;
                r.Dt = fr.Dt;
                r.Rent = fr.Rent;
                r.FLAGS = fr.FLAGS;
                r.recid = fr.RSID;
                showRentStepForm();
            };
        }
    });

    $().w2form({
        name: 'propertyRentStepForm',
        style: 'border: 0px; background-color: transparent;',
        // header: 'Property Detail',
        formURL: '/static/html/formRentStep.html',
        // url: '/v1/property',
        fields: [
            {field: 'RSLID',       type: 'int',   required: false},
            {field: 'RSID',        type: 'int',   required: false},
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
                    w2ui.rentStepsLayout.hide('right',true);
                    w2ui.propertyRentStepsGrid.render();
                }
            },
        },
        onRefresh: function(event) {
            event.onComplete = function(event) {
                EnableRentStepFormFields();
                SetMonthlyRentString();
            };
        },
        onChange: function(event) {
            event.onComplete = function(event) {
                SetMonthlyRentString();
            };
        },
    });


    $().w2layout({
        name: 'rentStepsLayout',
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

function finishRentStepsGridToolbar() {
    var t = w2ui.propertyRentStepsGrid.toolbar;
    t.add([
        { id: 'bt3', type: 'spacer' },
        { id: 'rsListType', type: 'menu-radio', icon: 'fa fa-star',
            text: function (item) {
                var text = item.selected;
                var el   = this.get('rsListType:' + item.selected);
                return 'Step Type: ' + el.text;
            },
            selected: 'rsListOpt',
            items: [
                { id: 'rsListOpt',  text: 'Period', icon: 'fa fa-tachometer' },
                { id: 'rsListDate', text: 'Date',   icon: 'fa fa-tachometer' },
            ]
        },
    ]);
    t.on('*', RentStepTypeChange);
}

function RentStepDelete() {
    var r = w2ui.propertyRentStepForm.record;
    var g = w2ui.propertyRentStepsGrid;
    var i = g.get(r.recid,true);
    if (typeof i == "number" && i >= 0) {
        var removed = g.records.splice(i,1);
        // console.log('removed = ' + removed);
    }
    w2ui.rentStepsLayout.hide('right',true);
    g.render();
}

function RentStepSave() {
    var r = w2ui.propertyRentStepForm.record;
    var g = w2ui.propertyRentStepsGrid;

    var x=w2ui.propertyRentStepForm.validate(true);
    if (x.length > 0) {
        return;
    }

    if (r.RSID < 0) {
        w2ui.propertyRentStepsGrid.add(r);
        w2ui.propertyRentStepForm.record.RSID = 0;  // to make sure that this one won't be added again
    }
    g.set(r.recid,r);

    w2ui.rentStepsLayout.hide('right',true);
    g.render();
}

function SetMonthlyRentString() {
    var y = w2ui.propertyRentStepForm.record.Rent/12;
    var s = "$" + number_format(y,2,'.',',');
    setInnerHTML("PRmonthly",s);
}

function RentStepTypeChange(event) {
    if (event.type != "click") {
        return;
    }
    //console.log('EVENT: '+ event.type + ' TARGET: '+ event.target, event);
    switch (event.target) {
    case "rsListType:rsListOpt":
        SetRentStepColumns(0);
        SetRentStepFLAGs(0);
        propData.rsType = 0;
        break;
    case "rsListType:rsListDate":
        SetRentStepColumns(1);
        SetRentStepFLAGs(1);
        propData.rsType = 1;
        break;
    }
}

function SetRentStepColumns(FLAGS) {
    var b = FLAGS & 0x1;
    var t = w2ui.propertyRentStepsGrid.toolbar.get("rsListType");
    if (b == 0) {
        w2ui.propertyRentStepsGrid.hideColumn("Dt");
        w2ui.propertyRentStepsGrid.showColumn("Opt");
        t.selected = "rsListOpt";
    } else {
        w2ui.propertyRentStepsGrid.hideColumn("Opt");
        w2ui.propertyRentStepsGrid.showColumn("Dt");
        t.selected = "rsListDate";
    }
    w2ui.propertyRentStepsGrid.toolbar.refresh();
}

// INPUTS
//   FLAGS =  0 means use Options
//            1 means use Dates
//-----------------------------------------------------------------------
function SetRentStepFLAGs(FLAGS) {
    for (var i = 0; i < w2ui.propertyRentStepsGrid.records.length; i++) {
        w2ui.propertyRentStepsGrid.records[i].FLAGS &= 0xeffffffffffffffe;
        w2ui.propertyRentStepsGrid.records[i].FLAGS |= FLAGS;
    }
}

function EnableRentStepFormFields() {
    if (1 == propData.rsType) {
        $('#Opt').prop('disabled', true);
        $('#Dt').prop('disabled', false);
    } else {
        $('#Opt').prop('disabled', false);
        $('#Dt').prop('disabled', true);
    }
}

function showRentStepForm() {
    w2ui.rentStepsLayout.content('right',w2ui.propertyRentStepForm);
    w2ui.rentStepsLayout.sizeTo('right',400);
    w2ui.rentStepsLayout.show('right',true);
}


function saveRentSteps() {
    //-----------------------------------------------------------------------
    // If we never loaded the rentsteps, then they weren't changed, so just
    // return success.
    //-----------------------------------------------------------------------
    if (!propData.bRentStepsLoaded) {
        return new Promise( function(resolve,reject) {
            if (true) {
                resolve("success");
            } else {
                reject("error");
            }
        });
    }

    //-----------------------------------------------------------------------
    // We have loaded the rentsteps, so we need to go through the save...
    //-----------------------------------------------------------------------
    var params = {
        cmd: "save",
        PRID: w2ui.propertyForm.record.PRID,
        records: w2ui.propertyRentStepsGrid.records
    };
    for (var i = 0; i < params.records.length; i++) {
        var d = new Date(params.records[i].Dt);
        params.records[i].Dt = d.toUTCString();
        var y = params.records[i].Opt;
        params.records[i].Opt = y.toString();  // need to make sure this is a string
    }
    var dat = JSON.stringify(params);
    var url = '/v1/rentsteps/' + w2ui.propertyForm.record.RSLID;

    return $.post(url, dat, null, "json")
    .done(function(data) {
        if (data.status === "error") {
            w2ui.propertyGrid.error('ERROR: '+ data.message);
        }
        propData.bRentStepsLoaded = false;
    })
    .fail(function(data){
            w2ui.propertyGrid.error("Save RentSteps failed. " + data);
    });
}
