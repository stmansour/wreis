/*global
    w2ui, app, $, console, dateFmtStr, propData,
*/

"use strict";
function newRentStepRecord() {
    var rs = {
        RSLID: 0,
        RSID: 0,
        Opt: 0,
        Dt: new Date(),
        Rent: 0,
        FLAGS: 0,
    }
}

function buildRentStepsUIElements() {

    $().w2grid({
        name: 'propertyRentStepsGrid',
        url: '/v1/rentsteps',
        show: {
            toolbar         : true,
            footer          : false,
            toolbarAdd      : true,   // indicates if toolbar add new button is visible
            toolbarDelete   : true,   // indicates if toolbar delete button is visible
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
            {field: 'RSID',        caption: 'RSID',       size: '60px', sortable: true, hidden: true},
            {field: 'RSLID',       caption: 'RSLID',      size: '60px', sortable: true, hidden: true},
            {field: 'FLAGS',       caption: 'FLAGS',      size: '60px', sortable: true, hidden: true},
            {field: 'Opt',         caption: 'Period',     size: '250px', sortable: true, hidden: false},
            {field: 'Dt',          caption: 'Date',       size: '80px', sotrable: true, hidden: false},
            {field: 'Rent',        caption: 'Rent',       size: '80px', sortable: true, hidden: false, render: 'money'},
            {field: 'CreateTime',  caption: 'CreateTime', size: '60px', sortable: true, hidden: true},
            {field: 'CreatedBy',   caption: 'CreatedBy',  size: '60px', sortable: true, hidden: true},
            {field: 'LastModTime', caption: 'LastModTime',size: '60px', sortable: true, hidden: true},
            {field: 'LastModBy',   caption: 'LastModBy',  size: '60px', sortable: true, hidden: true},
        ],
         onLoad: function(event) {
            event.onComplete = function() {
                propData.bRentStepsLoaded = true;
                SetRentStepColumns(GetOptionTextMode());  // since all records are the same in BIT 0, just look at first
            };
        },
        onAdd: function(event) {

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
                w2ui.rentStepsLayout.content('right',w2ui.propertyRentStepForm);
                w2ui.rentStepsLayout.sizeTo('right',400);
                w2ui.rentStepsLayout.show('right',true);
            }
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
            {field: 'CreatedBy',   type: 'text',  required: false},
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
        onClick: function(event) {
            event.onComplete = function() {
                var r = this.record;
            };
        },
        onRefresh: function(event) {
            event.onComplete = function(event) {
                console.log('propertyRentStepForm: Refresh completed');
                EnableRentStepFormFields();
            }
        }
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

function RentStepTypeChange(event) {
    if (event.type != "click") {
        return;
    }
    //console.log('EVENT: '+ event.type + ' TARGET: '+ event.target, event);
    switch (event.target) {
    case "rsListType:rsListOpt":
        SetRentStepColumns(0);
        SetRentStepFLAGs(0);
        break;
    case "rsListType:rsListDate":
        SetRentStepColumns(1);
        SetRentStepFLAGs(1);
        break;
    }
}

function SetRentStepColumns(FLAGS) {
    var b = FLAGS & 0x1;
    if (b == 0) {
        w2ui.propertyRentStepsGrid.hideColumn("Dt");
        w2ui.propertyRentStepsGrid.showColumn("Opt");
    } else {
        w2ui.propertyRentStepsGrid.hideColumn("Opt");
        w2ui.propertyRentStepsGrid.showColumn("Dt");
    }
}

function SetRentStepFLAGs(FLAGS) {
    for (var i = 0; i < w2ui.propertyRentStepsGrid.records.length; i++) {
        w2ui.propertyRentStepsGrid.records[i].FLAGS &= (~1);
        w2ui.propertyRentStepsGrid.records[i].FLAGS |= FLAGS;
    }
}

function EnableRentStepFormFields() {
    if (GetOptionTextMode()) {
        $('#Opt').prop('disabled', true);
        $('#Dt').prop('disabled', false);
    } else {
        $('#Opt').prop('disabled', false);
        $('#Dt').prop('disabled', true);
    }
}

// GetOptionTextMode return 0 if the mode is Opt, or 1 if it is Date
function GetOptionTextMode() {
    var FLAGS = 0;
    if (w2ui.propertyRentStepsGrid.records.length > 0) {
        FLAGS = w2ui.propertyRentStepsGrid.records[0].FLAGS;  // since all records are the same in BIT 0, just look at first
    }
    return FLAGS & 0x1;
}
