/*global
    w2ui, app, $, console, dateFmtStr, getDropDownSelectedIndex,
    setDropDownSelectedIndex,saveRentSteps,saveRenewOptions, varToUTCString,
    propertyStateOnLoad,setTimeout,
*/

"use strict";

var propData = {
    PRID: 0,                    // which property is currently being edited
    RSLID: 0,
    ROLID: 0,
    bPropLoaded: false,         // false -> either it's new or user clicked property in the propertyGrid, true -> just switching tabls
    bRentStepsLoaded: false,    // "  same as above for RentSteps
    bRenewOptionsLoaded: false, // "  same as above for RenewOptions
    bTrafficLoaded: false,      // " for Traffic
    bStateLoaded: false,        // " for state info
    statefilter: [1,2,3,4,5,6], // how to filter properties  (1-6) = open, (7) = closed
    formWidth: 575,             // how wide is the entry / edit form
    numStates: 7,               // states go from 1 to 7 -- this is a full complement of sates, the states array may have less
    states: [],                 // the server will be queried for these on existing properties, or filled with an inital state on new
    doneText: "#0611AA",
    doneBG: "#e0f0ff",
    notStartedText: "#888888",
    notStartedBG: "#e0e0e0",
    inProgressText: "#11AA11",
    inProgressBG: "#e0ffe0",
};

function initializeStateRecord() {
    var time0 = new Date("Jan 1, 1970");
    var now = new Date();
    var rec = {
        SIID: 0,
        PRID: 0,
        InitiatorUID: app.uid,
        InitiatorDt: now,
        InitiatorName: app.name,
        ApproverUID: 0,
        ApproverDt: time0,
        ApproverName: "",
        FlowState: 1,
        Reason: "",
        FLAGS: 0,
        LastModTime: now,
        LastModBy: app.uid,
        CreateTime: now,
        CreateBy: app.uid,
        CreateByName: "",
        LastModByName:"",
    };
    return rec;
}

function initializePropertyRecord() {
    var time0 = new Date("Jan 1, 1970");
    var now = new Date();
    var rec = {
            recid: 0,
            PRID: 0,
            Name: "",
            YearsInBusiness: 0,
            ParentCompany: "",
            URL: "",
            Symbol: "",
            Price: 0,
            DownPayment: 0,
            RentableArea: 0,
            RentableAreaUnits: 0,
            LotSize: 0,
            LotSizeUnits: 0,
            CapRate: 0,
            AvgCap: 0,
            FlowState: 1,  // initialize to State = 1
            FLAGS: 0,
            Ownership: 0,
            TenantTradeName: "",
            LeaseGuarantor: 0,
            LeaseType: 0,
            OriginalLeaseTerm: 0,
            TermRemainingOnLease: 0,
            ROLID: 0,
            RSLID: 0,
            Address: "",
            Address2: "",
            City: "",
            State: "",
            PostalCode: "",
            Country: "",
            LLResponsibilities: "",
            NOI: 0,
            HQAddress: "",
            HQAddress2: "",
            HQCity: "",
            HQState: "",
            HQPostalCode: "",
            HQCountry: "",
            Img1: "",
            Img2: "",
            Img3: "",
            Img4: "",
            Img5: "",
            Img6: "",
            Img7: "",
            Img8: "",

            CreateBy: 0,
            LastModBy: 0,

            BuildDate: time0,
            DeliveryDt: time0,
            RentCommencementDt: time0,
            LeaseExpirationDt: time0,
            CreateTime: now,
            LastModTime: now,
    };
    return rec;
}


function buildPropertyUIElements() {
    //------------------------------------------------------------------------
    //          propertyGrid
    //------------------------------------------------------------------------
    $().w2grid({
        name: 'propertyGrid',
        url: '/v1/property',
        show: {
            toolbar         : true,
            footer          : true,
            toolbarAdd      : true,    // indicates if toolbar add new button is visible
            toolbarDelete   : false,   // indicates if toolbar delete button is visible
            toolbarSave     : false,   // indicates if toolbar save button is visible
            selectColumn    : false,
            expandColumn    : false,
            toolbarEdit     : false,
            toolbarSearch   : true,
            toolbarInput    : true,
            searchAll       : false,
            toolbarReload   : true,
            toolbarColumns  : true,
        },
        //======================================================================
        // FLAGS
        //     1<<0  Drive Through?  0 = no, 1 = yes
        //	   1<<1  Roof & Structure Responsibility: 0 = Tenant, 1 = Landlord
        //	   1<<2  Right Of First Refusal: 0 = no, 1 = yes
        //======================================================================
         columns: [
            {field: 'Recid',                size: '60px', caption: 'Recid', sortable: true, hidden: true},
            {field: 'PRID',                 size: '60px', caption: 'PRID', sortable: true, hidden: true},
            {field: 'Name',                 size: '200px', caption: 'Name', sortable: true, hidden: false},
            {field: 'YearsInBusiness',      size: '60px', caption: 'YearsInBusiness', sortable: true, hidden: true},
            {field: 'ParentCompany',        size: '60px', caption: 'ParentCompany', sortable: true, hidden: true},
            {field: 'URL',                  size: '60px', caption: 'URL', sortable: true, hidden: true},
            {field: 'Symbol',               size: '60px', caption: 'Symbol', sortable: true, hidden: true},
            {field: 'Price',                size: '60px', caption: 'Price', sortable: true, hidden: true},
            {field: 'DownPayment',          size: '60px', caption: 'DownPayment', sortable: true, hidden: true},
            {field: 'RentableArea',         size: '60px', caption: 'RentableArea', sortable: true, hidden: true},
            {field: 'RentableAreaUnits',    size: '60px', caption: 'RentableAreaUnits', sortable: true, hidden: true},
            {field: 'LotSize',              size: '60px', caption: 'LotSize', sortable: true, hidden: true},
            {field: 'LotSizeUnits',         size: '60px', caption: 'LotSizeUnits', sortable: true, hidden: true},
            {field: 'CapRate',              size: '60px', caption: 'CapRate', sortable: true, hidden: true},
            {field: 'AvgCap',               size: '60px', caption: 'AvgCap', sortable: true, hidden: true},
            {field: 'BuildDate',            size: '60px', caption: 'BuildDate', sortable: true, hidden: true},
            {field: 'FlowState',            size: '60px', caption: 'FlowState', sortable: true, hidden: true},
            {field: 'FLAGS',                size: '60px', caption: 'FLAGS', sortable: true, hidden: true},
            {field: 'Ownership',            size: '60px', caption: 'Ownership', sortable: true, hidden: true},
            {field: 'TenantTradeName',      size: '60px', caption: 'TenantTradeName', sortable: true, hidden: true},
            {field: 'LeaseGuarantor',       size: '60px', caption: 'LeaseGuarantor', sortable: true, hidden: true},
            {field: 'LeaseType',            size: '60px', caption: 'LeaseType', sortable: true, hidden: true},
            {field: 'DeliveryDt',           size: '60px', caption: 'DeliveryDt', sortable: true, hidden: true},
            {field: 'OriginalLeaseTerm',    size: '60px', caption: 'OriginalLeaseTerm', sortable: true, hidden: true},
            {field: 'RentCommencementDt',   size: '60px', caption: 'RentCommencementDt', sortable: true, hidden: true},
            {field: 'LeaseExpirationDt',    size: '60px', caption: 'LeaseExpirationDt', sortable: true, hidden: true},
            {field: 'TermRemainingOnLease', size: '60px', caption: 'TermRemainingOnLease', sortable: true, hidden: true},
            {field: 'ROLID',                size: '60px', caption: 'ROLID', sortable: true, hidden: true},
            {field: 'RSLID',                size: '60px', caption: 'RSLID', sortable: true, hidden: true},
            {field: 'Address',              size: '60px', caption: 'Address', sortable: true, hidden: true},
            {field: 'Address2',             size: '60px', caption: 'Address2', sortable: true, hidden: true},
            {field: 'City',                 size: '100px', caption: 'City', sortable: true, hidden: false},
            {field: 'State',                size: '60px', caption: 'State', sortable: true, hidden: false},
            {field: 'PostalCode',           size: '60px', caption: 'PostalCode', sortable: true, hidden: false},
            {field: 'Country',              size: '60px', caption: 'Country', sortable: true, hidden: true},
            {field: 'LLResponsibilities',   size: '60px', caption: 'LLResponsibilities', sortable: true, hidden: true},
            {field: 'NOI',                  size: '60px', caption: 'NOI', sortable: true, hidden: true, render: 'money'},
            {field: 'HQAddress',            size: '60px', caption: 'HQAddress', sortable: true, hidden: true},
            {field: 'HQAddress2',           size: '60px', caption: 'HQAddress2', sortable: true, hidden: true},
            {field: 'HQCity',               size: '60px', caption: 'HQCity', sortable: true, hidden: true},
            {field: 'HQState',              size: '60px', caption: 'HQState', sortable: true, hidden: true},
            {field: 'HQPostalCode',         size: '60px', caption: 'HQPostalCode', sortable: true, hidden: true},
            {field: 'HQCountry',            size: '60px', caption: 'HQCountry', sortable: true, hidden: true},
            {field: 'Img1',                 size: '100px', caption: 'Img1', sortable: true, hidden: true},
            {field: 'Img2',                 size: '100px', caption: 'Img2', sortable: true, hidden: true},
            {field: 'Img3',                 size: '100px', caption: 'Img3', sortable: true, hidden: true},
            {field: 'Img4',                 size: '100px', caption: 'Img4', sortable: true, hidden: true},
            {field: 'Img5',                 size: '100px', caption: 'Img5', sortable: true, hidden: true},
            {field: 'Img6',                 size: '100px', caption: 'Img6', sortable: true, hidden: true},
            {field: 'Img7',                 size: '100px', caption: 'Img7', sortable: true, hidden: true},
            {field: 'Img8',                 size: '100px', caption: 'Img8', sortable: true, hidden: true},
            {field: 'CreateTime',           size: '60px', caption: 'CreateTime', sortable: true, hidden: true},
            {field: 'CreateBy',            size: '60px', caption: 'CreateBy', sortable: true, hidden: true},
            {field: 'LastModTime',          size: '60px', caption: 'LastModTime', sortable: true, hidden: true},
            {field: 'LastModBy',            size: '60px', caption: 'LastModBy', sortable: true, hidden: true},
        ],
        onClick: function(event) {
            event.onComplete = function (event) {
                var f = w2ui.propertyForm;
                var rec = w2ui.propertyGrid.get(event.recid);
                w2ui.propertyForm.recid = rec.PRID;
                propData.PRID = rec.PRID;

                f.url = "/v1/property/"+rec.PRID;
                f.refresh();
                f.reload();  // get this going as quickly as possible

                propData.RSLID = rec.RSLID;
                propData.ROLID = rec.ROLID;
                propData.bPropLoaded = false;
                propData.bRentStepsLoaded = false;
                propData.bRenewOptionsLoaded = false;
                propData.bTrafficLoaded = false;
                propData.bStateLoaded = false;
                w2ui.propertyFormLayout_main_tabs.click('proptabGeneral'); // click the general tab
                var l = w2ui.propertyFormLayout.get('main');
                if (typeof l.tabs != "undefined"){
                    if (typeof l.tabs.name == "string") {
                        l.tabs.click('proptabState');
                    }
                }
            };
        },
        onRequest: function(event) {
            // Include any postData needed
            // w2ui.propertyGrid.postData = {groupName: app.groupFilter};
        },
        onAdd: function (/*event*/) {
            var f = w2ui.propertyForm;
            f.record = initializePropertyRecord();
            f.recid = 0;
            f.url = "";
            f.refresh();
            propData.PRID = 0;  // new entry
            var s = initializeStateRecord();
            propData.states = [s];
            w2ui.propertyFormLayout.content('main', w2ui.propertyForm);
            w2ui.propertyFormLayout.content("bottom", w2ui.propertyFormBtns);
            w2ui.toplayout.content('right', w2ui.propertyFormLayout);
            w2ui.toplayout.sizeTo('right', propData.formWidth);
            w2ui.toplayout.render();
            w2ui.toplayout.show('right', true);
            var l = w2ui.propertyFormLayout.get('main');
            if (typeof l.tabs != "undefined"){
                if (typeof l.tabs.name == "string") {
                    l.tabs.click('proptabGeneral');
                }
            }
        },
        onRefresh: function(/*event*/) {
            // console.log('propertyGrid.onRefresh')
            //document.getElementById('mojoGroupFilter').value = app.groupFilter;
        },
        onLoad: function(event) {
            event.onComplete = function(event) {
                propData.statefilter = [1,2,3,4,5,6];
                propertySetPostData();
            };
            //document.getElementById('mojoGroupFilter').value = app.groupFilter;
        },
        onSearch: function(event) {
            console.log('onSearch event fired. event = ' + event);
        }
    });

    w2ui.propertyGrid.toolbar.add([
        { type: 'break' },
        { type: 'radio', id: 'openProperties', group: '1', text: 'Open', /* icon: 'fa fa-star',*/ checked: true },
        { type: 'radio', id: 'closedProperties', group: '1', text: 'Closed', /*icon: 'fa fa-heart'*/ },
        { type: 'radio', id: 'allProperties', group: '1', text: 'All', /*icon: 'fa fa-heart'*/ },
    ]);

    w2ui.propertyGrid.toolbar.onClick = function(event) {
        event.onComplete = function (event) {
            var found = false;
            switch (event.item.id) {
            case "openProperties": found=true; propData.statefilter = [1,2,3,4,5,6]; break;
            case "closedProperties": found=true; propData.statefilter = [7]; break;
            case "allProperties": found=true; propData.statefilter = [1,2,3,4,5,6,7]; break;
            }
            if (found) {
                propertySetPostData();
                w2ui.propertyGrid.reload();
            }
        };
    };

    function propertySetPostData() {
        w2ui.propertyGrid.postData = {
            statefilter: propData.statefilter,
        };
    }

    // create a layout.
    //  top    - toolbar
    //  main   - tabs plus whatever is needed to go under the tab to make the form
    //           we need.  It could be a standard form, or it could be a grid, or
    //           any combination.
    //  bottom - buttons
    //
    //  right  - an area to open if we need to pop up another form dialog... for
    //           example.
    //--------------------------------------------------------------------------
    $().w2layout({
        name: 'propertyFormLayout',
        header: 'Property Detail',
        panels: [
            {
                type: 'top',
                size: 35,
                style: 'border: 1px solid silver;',
                content: "",
                toolbar: {
                    style: "height: 35px; background-color: #eee; border: 0px;",
                    items: [
                        { id: 'btnNotes', type: 'button', icon: 'fa fa-sticky-note-o' },
                        { id: 'bt3', type: 'spacer' },
                        { id: 'btnClose', type: 'button', icon: 'fa fa-times' },
                    ],
                    onClick: function (event) {
                        switch(event.target) {
                        case 'btnClose':
                            w2ui.toplayout.hide('right', true);
                            w2ui.propertyGrid.render();
                            break;
                        }
                    },
                },
            },
            {
                type: 'main',
                overflow: "hidden",
                style: 'background-color: white; border: 1px solid silver; padding: 0px;',
                tabs: {
                    style: "padding-top: 10px;",
                    active: 'proptabGeneral',
                    tabs: [
                        { id: 'proptabState', caption: 'State' },
                        { id: 'proptabGeneral', caption: 'General' },
                        { id: 'proptabRentSteps', caption: 'Rent Steps' },
                        { id: 'proptabRenewOptions', caption: 'Renew Options' },
                        { id: 'proptabTraffic', caption: 'Traffic' },
                        { id: 'proptabPhotos', caption: 'Photos' },
                    ],
                    //---------------------------------
                    //  HANDLE THE TAB CLICKS...
                    //---------------------------------
                    onClick: function (event) {
                        // console.log('event.target = ' + event.target);
                        switch (event.target) {
                            case "proptabState":
                            setPropertyLayout(event.target);
                            break;

                            case "proptabGeneral":
                            setPropertyLayout(event.target);
                            break;

                            case 'proptabRentSteps':
                            setPropertyLayout(event.target);
                            break;

                            case 'proptabRenewOptions':
                            setPropertyLayout(event.target);
                            break;

                            case 'proptabTraffic':
                            setPropertyLayout(event.target);
                            break;

                            case 'proptabPhotos':
                            setPropertyLayout(event.target);
                            break;
                        }
                    }
                }
            },
            {
                type: 'bottom', size: 60, // style: 'background-color: white;  border-top: 1px solid silver; text-align: center; padding: 15px;',
            },
        ],
    });

    //------------------------------------------------------------------------
    //          Property Form
    //------------------------------------------------------------------------
    $().w2form({
        name: 'propertyForm',
        style: 'border: 0px; background-color: transparent;',
        // header: 'Property Detail',
        formURL: '/static/html/formproperty.html',
        url: '/v1/property',
        fields: [
            {field: 'recid',                type: 'int',  required: false },
            {field: 'PRID',                 type: 'int',  required: false},
            {field: 'Name',                 type: 'text', required: false},
            {field: 'YearsInBusiness',      type: 'int',  required: false},
            {field: 'ParentCompany',        type: 'text', required: false},
            {field: 'URL',                  type: 'text', required: false},
            {field: 'Symbol',               type: 'text',  required: false},
            {field: 'Price',                type: 'money', required: false},
            {field: 'DownPayment',          type: 'money', required: false},
            {field: 'RentableArea',         type: 'int',   required: false},
            {field: 'RentableAreaUnits',    type: 'hidden', required: false},
            {field: 'LotSize',              type: 'int',    required: false},
            {field: 'LotSizeUnits',         type: 'hidden', required: false},
            {field: 'CapRate',              type: 'percent',  required: false},
            {field: 'AvgCap',               type: 'percent',  required: false},
            {field: 'BuildDate',            type: 'date', required: false},
            {field: 'FlowState',            type: 'hidden†', required: false},
            {field: 'FLAGS',                type: 'text', required: false},
            {field: 'Ownership',            type: 'hidden', required: false},
            {field: 'TenantTradeName',      type: 'text', required: false},
            {field: 'LeaseGuarantor',       type: 'text', required: false},
            {field: 'LeaseType',            type: 'hidden', required: false},
            {field: 'DeliveryDt',           type: 'date', required: false},
            {field: 'OriginalLeaseTerm',    type: 'int', required: false},
            {field: 'RentCommencementDt',   type: 'date', required: false},
            {field: 'LeaseExpirationDt',    type: 'date', required: false},
            {field: 'TermRemainingOnLease', type: 'int', required: false},
            {field: 'ROLID',                type: 'hidden', required: false},
            {field: 'RSLID',                type: 'hidden', required: false},
            {field: 'Address',              type: 'text', required: false},
            {field: 'Address2',             type: 'text', required: false},
            {field: 'City',                 type: 'text', required: false},
            {field: 'State',                type: 'text', required: false},
            {field: 'PostalCode',           type: 'text', required: false},
            {field: 'Country',              type: 'text', required: false},
            {field: 'LLResponsibilities',   type: 'text', required: false},
            {field: 'NOI',                  type: 'money', required: false},
            {field: 'HQAddress',            type: 'text', required: false},
            {field: 'HQAddress2',           type: 'text', required: false},
            {field: 'HQCity',               type: 'text', required: false},
            {field: 'HQState',              type: 'text', required: false},
            {field: 'HQPostalCode',         type: 'text', required: false},
            {field: 'HQCountry',            type: 'text', required: false},
            {field: 'Img1',                 type: 'hidden', required: false},
            {field: 'Img2',                 type: 'hidden', required: false},
            {field: 'Img3',                 type: 'hidden', required: false},
            {field: 'Img4',                 type: 'hidden', required: false},
            {field: 'Img5',                 type: 'hidden', required: false},
            {field: 'Img6',                 type: 'hidden', required: false},
            {field: 'Img7',                 type: 'hidden', required: false},
            {field: 'Img8',                 type: 'hidden', required: false},
            {field: 'CreateTime',           type: 'text', required: false},
            {field: 'CreateBy',            type: 'text', required: false},
            {field: 'LastModTime',          type: 'text', required: false},
            {field: 'LastModBy',            type: 'text', required: false},
        ],
        // toolbar: {
        //     items: [
        //         { id: 'btnNotes', type: 'button', icon: 'fa fa-sticky-note-o' },
        //         { id: 'bt3', type: 'spacer' },
        //         { id: 'btnClose', type: 'button', icon: 'fa fa-times' },
        //     ],
        //     onClick: function (event) {
        //         if (event.target == 'btnClose') {
        //                     w2ui.toplayout.hide('right',true);
        //                     w2ui.propertyGrid.render();
        //         }
        //     },
        // },
        onLoad: function(event) {
            event.onComplete = function() {
                var r = this.record;
                var y = new Date(r.BuildDate);
                r.BuildDate = dateFmtStr(y);
                y = new Date(r.DeliveryDt);
                r.DeliveryDt = dateFmtStr(y);
                y = new Date(r.RentCommencementDt);
                r.RentCommencementDt = dateFmtStr(y);
                y = new Date(r.LeaseExpirationDt);
                r.LeaseExpirationDt = dateFmtStr(y);
                r.CapRate *= 100;
                r.AvgCap *= 100;
                setDropDownSelectedIndex("LotSizeUnitsDD",r.LotSizeUnits);
                setDropDownSelectedIndex("OwnershipDD",r.Ownership);
                setDropDownSelectedIndex("TermRemainingOnLeaseUnitsDD",r.TermRemainingOnLeaseUnits);
                setDropDownSelectedIndex("LeaseTypeDD",r.LeaseType);
                setDropDownSelectedIndex("LeaseGuarantorDD",r.LeaseGuarantor);

                propData.bPropLoaded = true;
                propertyStateOnLoad(); // need to call this now that we know the state
            };
        },
    });

    $().w2form({
        name: 'propertyFormBtns',
        url: '/v1/property',
        formURL: '/static/html/propertyFormBtns.html',

        actions: {
            save: function () {
                    $.when(
                        savePropertyForm(),
                        saveRentSteps(),
                        saveRenewOptions(),
                        saveTraffic()
                    )
                    .done( function() {
                        propertySaveDoneCB();
                    })
                    .fail( function() {
                        var s = 'Save Property encountered an error';
                        w2ui.propertyGrid.error(s);
                        propertySaveDoneCB();
                    });
            },
            delete: function() {
                var request={cmd:"delete",selected: [w2ui.propertyForm.record.PRID]};
                $.post('/v1/person/'+w2ui.propertyForm.record.PRID, JSON.stringify(request))
                .done(function(data) {
                    if (typeof data == 'string') {  // it's weird, a successful data add gets parsed as an object, an error message does not
                        var msg = JSON.parse(data);
                        w2ui.propertyForm.error(msg.message);
                        return;
                    }
                    if (data.status != 'success') {
                        w2ui.propertyForm.error(data.message);
                    }
                });
                w2ui.toplayout.hide('right',true);
                w2ui.propertyGrid.reload();
            },
            reset: function() {
                var f = w2ui.asmInstForm;
                console.log('reset: ASMID = ' + f.record.ASMID );
            }
        },
   });

}

// savePropertyForm grabs all the data that is associated with the propertForm,
//      converts anything that needs attention and calls the server's save
//      function.
//------------------------------------------------------------------------------
function savePropertyForm() {
    var rec = w2ui.propertyForm.record;

    //-----------------------------------------
    // Handle any conversions necessary...
    //-----------------------------------------
    rec.AvgCap /= 100;  // convert back to decimal number
    rec.CapRate /= 100; // convert back to decimal number
    rec.BuildDate = varToUTCString(rec.BuildDate);
    rec.DeliveryDt = varToUTCString(rec.DeliveryDt);
    rec.RentCommencementDt = varToUTCString(rec.RentCommencementDt);
    rec.LeaseExpirationDt = varToUTCString(rec.LeaseExpirationDt);
    rec.CreateTime = varToUTCString(rec.CreateTime);
    rec.LastModTime = varToUTCString(rec.LastModTime);

    rec.LotSizeUnits = getDropDownSelectedIndex("LotSizeUnitsDD");
    rec.Ownership = getDropDownSelectedIndex("OwnershipDD");
    rec.TermRemainingOnLeaseUnits = getDropDownSelectedIndex("TermRemainingOnLeaseUnitsDD");
    rec.LeaseType = getDropDownSelectedIndex("LeaseTypeDD");
    rec.LeaseGuarantor = getDropDownSelectedIndex("LeaseGuarantorDD");

    //-----------------------------------------
    // Now send it to the server
    //-----------------------------------------
    var params = {
        cmd: "save",
        record: rec,
        states: [0,0,0,0,0,0,0]
    };

    var dat = JSON.stringify(params);
    var url = '/v1/property/' + rec.PRID;

    return $.post(url, dat, null, "json")
    .done(function(data) {
        // if (data.status === "success") {
        // }
        if (data.status === "error") {
            w2ui.propertyGrid.error('ERROR: '+ data.message);
        }
    })
    .fail(function(data){
            w2ui.propertyGrid.error("Save RentableLeaseStatus failed. " + data);
    });

}

function propertySaveDoneCB() {
    w2ui.toplayout.hide('right',true);
    w2ui.propertyGrid.reload();
}

// setPropertyLayout is used to display the property form in the UI and handle
// the tab clicking.
//
// INPUTS
//      PRID - int64, property id
//      tab  - string, name of the tab that was pressed
function setPropertyLayout(tab) {
    w2ui.propertyFormLayout.content("bottom", w2ui.propertyFormBtns);

    switch (tab) {

    case "proptabState":
        w2ui.propertyFormLayout.load('main', '/static/html/formState.html', null,propertyStateOnLoad);
        setTimeout(propertyStateOnLoad, 100 );
        // setTimeout(propertyStateOnLoad, 1000 );
        break;

    case "proptabGeneral":
        if (propData.bPropLoaded) {
            w2ui.propertyForm.url = '';
        } else {
            w2ui.propertyForm.url = '/v1/property/' + propData.PRID;
        }
        w2ui.propertyFormLayout.content('main', w2ui.propertyForm);
        break;

    case "proptabRentSteps":
        if (propData.bRentStepsLoaded) {
            w2ui.propertyRentStepsGrid.url = '';
        } else {
            w2ui.propertyRentStepsGrid.clear();
            w2ui.propertyRentStepsGrid.url = '/v1/rentsteps/' + propData.RSLID;
        }
        w2ui.rentStepsLayout.content('main',w2ui.propertyRentStepsGrid);
        w2ui.propertyFormLayout.content('main',w2ui.rentStepsLayout);
        break;

    case "proptabRenewOptions":
        if (propData.bRenewOptionsLoaded) {
            w2ui.propertyRenewOptionsGrid.url = '';
        } else {
            w2ui.propertyRenewOptionsGrid.clear();
            w2ui.propertyRenewOptionsGrid.url = '/v1/renewoptions/' + propData.ROLID;
        }
        w2ui.renewOptionsLayout.content('main',w2ui.propertyRenewOptionsGrid);
        w2ui.propertyFormLayout.content('main',w2ui.renewOptionsLayout);
        break;

    case "proptabTraffic":
        if (propData.bTrafficLoaded) {
            w2ui.propertyTrafficGrid.url = '';
        } else {
            w2ui.propertyTrafficGrid.clear();
            w2ui.propertyTrafficGrid.url = '/v1/trafficitems/' + propData.PRID;
        }
        w2ui.propertyTrafficLayout.content('main',w2ui.propertyTrafficGrid);
        w2ui.propertyFormLayout.content('main',w2ui.propertyTrafficLayout);
        break;
    case "proptabPhotos":
        w2ui.propertyPhotosLayout.load('main','/static/html/formPhotos.html');
        w2ui.propertyFormLayout.content('main',w2ui.propertyPhotosLayout);
        break;
    }

    showForm();
}

function showForm() {
    // SHOW the right panel now
    w2ui.toplayout.content('right', w2ui.propertyFormLayout);
    w2ui.toplayout.sizeTo('right', propData.formWidth);
    w2ui.toplayout.show('right', true);
}
