/*global
    w2ui, app, $, console, dateFmtStr, getDropDownSelectedIndex,
    setDropDownSelectedIndex,saveRentSteps,saveRenewOptions, varToUTCString,
    propertyStateOnLoad,setTimeout,closeStateChangeDialog,setPropertyFormActionButtons,
    closePropertyForm, saveTraffic, setTermRemaining, monthDiff,
    setInnerHTML,getDropDownSafe
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
    statefilter: [1, 2, 3, 4, 5, 6, 7, 8], // how to filter properties  (1-8) = open, (9) = closed
    showTerminated: 0,          // 0 = don't show terminated properties, 1 = show terminated properties
    myQueue: 0,                 // 0 = don't show my queue, 1 = show my queue
    formWidth: 575,             // how wide is the entry / edit form
    numStates: 8,               // states go from 1 to 8 -- this is a full complement of sates, the states array may have less
    states: [],                 // the server will be queried for these on existing properties, or filled with an inital state on new
    rsType: 0,                  // 0 = options, 1 = date
    roType: 0,                  // 0 = options, 1 = date
    doneText: "#0611AA",
    doneBG: "#e0f0ff",
    notStartedText: "#777777",
    notStartedBG: "#f8f8f8",
    inProgressText: "#11AA11",
    inProgressBG: "#e0ffe0",
};

function initializeStateRecord() {
    var time0 = new Date("Jan 1, 1970");
    var now = new Date();
    var rec = {
        SIID: 0,
        PRID: 0,
        OwnerUID: app.uid,
        OwnerDt: now,
        OwnerName: app.name,
        ApproverUID: app.uid,
        ApproverDt: time0,
        ApproverName: app.name,
        FlowState: 1,
        Reason: "",
        FLAGS: 0,
        LastModTime: now,
        LastModBy: app.uid,
        CreateTime: now,
        CreateBy: app.uid,
        CreateByName: "",
        LastModByName: "",
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
        YearFounded: 0,
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
        OwnershipType: 0,
        TenantTradeName: "",
        LeaseGuarantor: 0,
        LeaseType: 0,
        OriginalLeaseTerm: 0,
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
        HQCity: "",
        HQState: "",
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

        BuildYear: 0,
        RenovationYear: 0,
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
        method: 'POST',
        show: {
            toolbar: true,
            footer: true,
            toolbarAdd: true,    // indicates if toolbar add new button is visible
            toolbarDelete: false,   // indicates if toolbar delete button is visible
            toolbarSave: false,   // indicates if toolbar save button is visible
            selectColumn: false,
            expandColumn: false,
            toolbarEdit: false,
            toolbarSearch: true,
            toolbarInput: true,
            searchAll: false,
            toolbarReload: true,
            toolbarColumns: true,
        },
        postData: {
            cmd: 'get',
            // statefilter: propData.statefilter,
            // showTerminated: propData.showTerminated,
            // myQueue: propData.myQueue,
        },
        //======================================================================
        // FLAGS
        //     1<<0  Drive Through?  0 = no, 1 = yes
        //	   1<<1  Roof & Structure Responsibility: 0 = Tenant, 1 = Landlord
        //	   1<<2  Right Of First Refusal: 0 = no, 1 = yes
        //======================================================================
        columns: [
            { field: 'recid', text: 'recid', size: '60px', sortable: true, hidden: true },
            { field: 'PRID', text: 'PRID', size: '60px', sortable: true, hidden: false },
            { field: 'Name', text: 'Name', size: '200px', sortable: true, hidden: false },
            { field: 'YearFounded', text: 'YearFounded', size: '60px', sortable: true, hidden: true },
            { field: 'ParentCompany', text: 'ParentCompany', size: '60px', sortable: true, hidden: true },
            { field: 'URL', text: 'URL', size: '60px', sortable: true, hidden: true },
            { field: 'Symbol', text: 'Symbol', size: '60px', sortable: true, hidden: true },
            { field: 'Price', text: 'Price', size: '60px', sortable: true, hidden: true },
            { field: 'DownPayment', text: 'DownPayment', size: '60px', sortable: true, hidden: true },
            { field: 'RentableArea', text: 'RentableArea', size: '60px', sortable: true, hidden: true },
            { field: 'RentableAreaUnits', text: 'RentableAreaUnits', size: '60px', sortable: true, hidden: true },
            { field: 'LotSize', text: 'LotSize', size: '60px', sortable: true, hidden: true },
            { field: 'LotSizeUnits', text: 'LotSizeUnits', size: '60px', sortable: true, hidden: true },
            { field: 'CapRate', text: 'CapRate', size: '60px', sortable: true, hidden: true },
            { field: 'AvgCap', text: 'AvgCap', size: '60px', sortable: true, hidden: true },
            { field: 'BuildYear', text: 'BuildYear', size: '60px', sortable: true, hidden: true },
            { field: 'RenovationYear', text: 'RenovationYear', size: '60px', sortable: true, hidden: true },
            { field: 'FlowState', text: 'FlowState', size: '60px', sortable: true, hidden: true },
            { field: 'FLAGS', text: 'FLAGS', size: '60px', sortable: true, hidden: true },
            { field: 'OwnershipType', text: 'OwnershipType', size: '60px', sortable: true, hidden: true },
            { field: 'Ownership', text: 'Ownership', size: '60px', sortable: true, hidden: true },
            { field: 'TenantTradeName', text: 'TenantTradeName', size: '60px', sortable: true, hidden: true },
            { field: 'LeaseGuarantor', text: 'LeaseGuarantor', size: '60px', sortable: true, hidden: true },
            { field: 'LeaseType', text: 'LeaseType', size: '60px', sortable: true, hidden: true },
            { field: 'OriginalLeaseTerm', text: 'OriginalLeaseTerm', size: '60px', sortable: true, hidden: true },
            { field: 'RentCommencementDt', text: 'RentCommencementDt', size: '60px', sortable: true, hidden: true },
            { field: 'LeaseExpirationDt', text: 'LeaseExpirationDt', size: '60px', sortable: true, hidden: true },
            { field: 'ROLID', text: 'ROLID', size: '60px', sortable: true, hidden: true },
            { field: 'RSLID', text: 'RSLID', size: '60px', sortable: true, hidden: true },
            { field: 'Address', text: 'Address', size: '60px', sortable: true, hidden: true },
            { field: 'Address2', text: 'Address2', size: '60px', sortable: true, hidden: true },
            { field: 'City', text: 'City', size: '100px', sortable: true, hidden: false },
            { field: 'State', text: 'State', size: '60px', sortable: true, hidden: false },
            { field: 'PostalCode', text: 'PostalCode', size: '60px', sortable: true, hidden: false },
            { field: 'Country', text: 'Country', size: '60px', sortable: true, hidden: true },
            { field: 'LLResponsibilities', text: 'LLResponsibilities', size: '60px', sortable: true, hidden: true },
            { field: 'NOI', text: 'NOI', size: '60px', sortable: true, hidden: true, render: 'money' },
            { field: 'HQCity', text: 'HQCity', size: '60px', sortable: true, hidden: true },
            { field: 'HQState', text: 'HQState', size: '60px', sortable: true, hidden: true },
            { field: 'Img1', text: 'Img1', size: '100px', sortable: true, hidden: true },
            { field: 'Img2', text: 'Img2', size: '100px', sortable: true, hidden: true },
            { field: 'Img3', text: 'Img3', size: '100px', sortable: true, hidden: true },
            { field: 'Img4', text: 'Img4', size: '100px', sortable: true, hidden: true },
            { field: 'Img5', text: 'Img5', size: '100px', sortable: true, hidden: true },
            { field: 'Img6', text: 'Img6', size: '100px', sortable: true, hidden: true },
            { field: 'Img7', text: 'Img7', size: '100px', sortable: true, hidden: true },
            { field: 'Img8', text: 'Img8', size: '100px', sortable: true, hidden: true },
            { field: 'CreateTime', text: 'CreateTime', size: '60px', sortable: true, hidden: true },
            { field: 'CreateBy', text: 'CreateBy', size: '60px', sortable: true, hidden: true },
            { field: 'LastModTime', text: 'LastModTime', size: '60px', sortable: true, hidden: true },
            { field: 'LastModBy', text: 'LastModBy', size: '60px', sortable: true, hidden: true },
        ],
        onClick: function (event) {
            event.onComplete = function (event) {
                setPropertyNotLoaded();
                loadPropertyForm(w2ui.propertyGrid.records[event.recid].PRID);
            };
        },
        onAdd: function (/*event*/) {
            var f = w2ui.propertyForm;
            f.record = initializePropertyRecord();
            f.recid = 0;
            f.url = "";
            f.refresh();
            setPropertyNotLoaded();
            if (typeof w2ui.propertyRentStepsGrid.records === "object") {
                w2ui.propertyRentStepsGrid.records = [];
            }
            if (typeof w2ui.propertyRenewOptionsGrid.records === "object") {
                w2ui.propertyRenewOptionsGrid.records = [];
            }
            if (typeof w2ui.propertyTrafficGrid.records === "object") {
                w2ui.propertyTrafficGrid.records = [];
            }
            propData.PRID = 0;  // new entry
            var s = initializeStateRecord();
            propData.states = [s];
            w2ui.propertyFormLayout.html('main', w2ui.propertyForm);
            w2ui.propertyFormLayout.html("bottom", w2ui.propertyFormBtns);
            w2ui.toplayout.html('right', w2ui.propertyFormLayout);
            w2ui.toplayout.sizeTo('right', propData.formWidth);
            w2ui.toplayout.render();
            w2ui.toplayout.show('right', true);
            var l = w2ui.propertyFormLayout.get('main');
            if (typeof l.tabs != "undefined") {
                if (typeof l.tabs.name == "string") {
                    l.tabs.click('proptabGeneral');
                }
            }
        },
        onRequest: function (/*event*/) {
            propertySetPostData();
        },
        onLoad: function (event) {
            var f = w2ui.propertyForm;
            for (var i = 0; i < w2ui.propertyGrid.records.length; i++) {
                w2ui.propertyGrid.records[i].recid = w2ui.propertyGrid.records[i].PRID;
            }
        },
    });

    w2ui.propertyGrid.toolbar.add([
        { type: 'break' },
        { type: 'check', id: 'myQueue', text: 'Queue', checked: false },
        { type: 'break' },
        { type: 'radio', id: 'openProperties', group: '1', text: 'Open', /* icon: 'fa fa-star',*/ checked: true },
        { type: 'radio', id: 'closedProperties', group: '1', text: 'Closed', /*icon: 'fa fa-heart'*/ },
        { type: 'radio', id: 'allProperties', group: '1', text: 'All', /*icon: 'fa fa-heart'*/ },
        { type: 'break' },
        { type: 'check', id: 'showTerminated', text: 'Terminated', /*icon: 'fa fa-heart'*/ },
    ]);

    w2ui.propertyGrid.toolbar.onClick = function (event) {
        event.onComplete = function (event) {
            var found = false;
            switch (event.item.id) {
                case "openProperties": found = true; propData.statefilter = [1, 2, 3, 4, 5, 6, 7, 8]; break;
                case "closedProperties": found = true; propData.statefilter = [9]; break;
                case "allProperties": found = true; propData.statefilter = [1, 2, 3, 4, 5, 6, 7, 8, 9]; break;
                case "showTerminated": found = true; propData.showTerminated = propData.showTerminated ? 0 : 1; break;
                case "myQueue": found = true; propData.myQueue = (propData.myQueue == 1) ? 0 : 1; break;
            }
            if (found) {
                propertySetPostData();
                closePropertyForm();
                w2ui.propertyGrid.reload();
            }
        };
    };

    function propertySetPostData() {
        w2ui.propertyGrid.postData = {
            cmd: 'get',
            statefilter: propData.statefilter,
            showTerminated: propData.showTerminated,
            myQueue: propData.myQueue,
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
        // header: 'Property Detail',
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
                        switch (event.target) {
                            case 'btnClose':
                                closePropertyForm();
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
                        { id: 'proptabState', text: 'State' },
                        { id: 'proptabGeneral', text: 'General' },
                        { id: 'proptabRentSteps', text: 'Rent Steps' },
                        { id: 'proptabRenewOptions', text: 'Renew Options' },
                        { id: 'proptabTraffic', text: 'Traffic' },
                        { id: 'proptabPhotos', text: 'Photos' },
                    ],
                    //---------------------------------
                    //  HANDLE THE TAB CLICKS...
                    //---------------------------------
                    onClick: function (event) {
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
            {
                type: 'right', size: 0,
            },
        ],
    });

    //------------------------------------------------------------------------
    //          Property Form
    //------------------------------------------------------------------------
    $().w2form({
        name: 'propertyForm',
        style: 'border: 0px; background-color: transparent;',
        formURL: '/static/html/formproperty.html',
        url: '/v1/property',
        fields: [
            { field: 'recid', type: 'int', required: false },
            { field: 'PRID', type: 'int', required: false },
            { field: 'Name', type: 'text', required: true },
            { field: 'YearFounded', type: 'int', required: false },
            { field: 'ParentCompany', type: 'text', required: false },
            { field: 'URL', type: 'text', required: false },
            { field: 'Symbol', type: 'text', required: false },
            { field: 'Price', type: 'money', required: false },
            { field: 'DownPayment', type: 'money', required: false },
            { field: 'RentableArea', type: 'int', required: false },
            { field: 'RentableAreaUnits', type: 'hidden', required: false },
            { field: 'LotSize', type: 'float', required: false },
            { field: 'LotSizeUnits', type: 'hidden', required: false },
            { field: 'CapRate', type: 'percent', required: false },
            { field: 'AvgCap', type: 'percent', required: false },
            { field: 'BuildYear', type: 'number', required: false },
            { field: 'RenovationYear', type: 'number', required: false },
            { field: 'FlowState', type: 'hidden†', required: false },
            { field: 'FLAGS', type: 'text', required: false },
            { field: 'OwnershipType', type: 'hidden', required: false },
            { field: 'TenantTradeName', type: 'text', required: false },
            { field: 'LeaseGuarantor', type: 'text', required: false },
            { field: 'LeaseType', type: 'hidden', required: false },
            { field: 'OriginalLeaseTerm', type: 'int', required: false },
            { field: 'RentCommencementDt', type: 'date', required: false },
            { field: 'LeaseExpirationDt', type: 'date', required: false },
            { field: 'ROLID', type: 'hidden', required: false },
            { field: 'RSLID', type: 'hidden', required: false },
            { field: 'Address', type: 'text', required: false },
            { field: 'Address2', type: 'text', required: false },
            { field: 'City', type: 'text', required: false },
            { field: 'State', type: 'text', required: false },
            { field: 'PostalCode', type: 'text', required: false },
            { field: 'Country', type: 'text', required: false },
            { field: 'LLResponsibilities', type: 'text', required: false },
            { field: 'NOI', type: 'money', required: false },
            { field: 'HQCity', type: 'text', required: false },
            { field: 'HQState', type: 'text', required: false },
            { field: 'Img1', type: 'hidden', required: false },
            { field: 'Img2', type: 'hidden', required: false },
            { field: 'Img3', type: 'hidden', required: false },
            { field: 'Img4', type: 'hidden', required: false },
            { field: 'Img5', type: 'hidden', required: false },
            { field: 'Img6', type: 'hidden', required: false },
            { field: 'Img7', type: 'hidden', required: false },
            { field: 'Img8', type: 'hidden', required: false },
            { field: 'CreateTime', type: 'text', required: false },
            { field: 'CreateBy', type: 'text', required: false },
            { field: 'LastModTime', type: 'text', required: false },
            { field: 'LastModBy', type: 'text', required: false },
        ],
        onLoad: function (event) {
            event.onComplete = function (event) {
                var r = this.record;
                r.RentCommencementDt = displayDateString(r.RentCommencementDt);
                r.LeaseExpirationDt = displayDateString(r.LeaseExpirationDt);
                r.CapRate *= 100;
                r.AvgCap *= 100;
                propData.bPropLoaded = true;
                propertyStateOnLoad(); // need to call this now that we know the state
            };
        },
        onRefresh: function (event) {
            event.onComplete = function () {
                var r = w2ui.propertyForm.record;
                // var Own = ((r.FLAGS & (1<<3)) == 0) ? 0 : 1;
                setDropDownSelectedIndex("LotSizeUnitsDD", r.LotSizeUnits);
                setDropDownSelectedIndex("OwnershipTypeDD", r.OwnershipType);
                setDropDownSelectedIndex("OwnershipDD", ((r.FLAGS & (1 << 3)) == 0) ? 0 : 1);
                setDropDownSelectedIndex("LeaseTypeDD", r.LeaseType);
                setDropDownSelectedIndex("LeaseGuarantorDD", r.LeaseGuarantor);
                setDropDownSelectedIndex("RoofResponsibilityDD", ((r.FLAGS & (1 << 1)) == 0) ? 0 : 1);
                // setTermRemaining();
            };
        },
        onChange: function (event) {
            event.onComplete = function (event) {
                switch (event.target) {
                    case "RentCommencementDt":
                    case "LeaseExpirationDt":
                        // setTermRemaining();
                        break;
                }
            };
        },
    });

    $().w2form({
        name: 'propertyFormBtns',
        url: '/v1/property',
        formURL: '/static/html/propertyFormBtns.html',

        actions: {
            save: function () {
                if (propData.tabGenDispCount > 0) {
                    //temporary fix...
                    if (w2ui.propertyForm.record.Name == "") {
                        return;
                    }
                    // var x=w2ui.propertyForm.validate(true);
                    // if (x.length > 0) {
                    //     return;
                    // }
                }
                savePropertyParts();

            },
            cancel: function () {
                closePropertyForm();
            }
        },
    });
}

// savePropertyParts will save the property in an efficient way to handle new
// properties or updates.  If the property is new, it must first save the Property
// structure in order to get back the PRID.  The PRID is needed for saving all
// the other structures.  If propertyForm.record.PRID < 1 it means we're saving
// a new property.  Save it first, then set the PRID so that the remaining calls
// can reference it.
//-------------------------------------------------------------------------------
function savePropertyParts() {
    var PRID = w2ui.propertyForm.record.PRID;
    if (PRID > 0) {
        savePropertyPartsPhase2(true);
        return;
    }
    //-------------------------------------------
    // need to save property first to get PRID
    //-------------------------------------------
    $.when(
        savePropertyForm()
    )
        .done(function () {
            var PRID = w2ui.propertyForm.record.PRID;
            if (0 == PRID) {
                // an error must have occurred.  don't proceed.
                return;
            }
            savePropertyPartsPhase2(false);  // proceed with saving everything else
        })
        .fail(function () {
            var s = 'Save Property encountered an error';
            w2ui.propertyGrid.error(s);
            propertySaveDoneCB();
        });
}
// savePropertyPartsPhase2
// If x is false, then we don't need to call savePropertyForm because it was
// a new property and had to be saved before everything else so that we could
// establish the PRID.
//
// INPUTS
//    x  -  true means save everything including the propertyForm.
//          false means save everything but exclude the propertyForm.
//-------------------------------------------------------------------------------
function savePropertyPartsPhase2(x) {
    if (x) {
        $.when(
            savePropertyForm(),
            saveRentSteps(),
            saveRenewOptions(),
            saveTraffic()
        )
            .done(function () {
                propertySaveDoneCB();
            })
            .fail(function () {
                var s = 'Save Property encountered an error';
                w2ui.propertyGrid.error(s);
                propertySaveDoneCB();
            });
        return;
    }

    $.when(
        saveRentSteps(),
        saveRenewOptions(),
        saveTraffic()
    )
        .done(function () {
            propertySaveDoneCB();
        })
        .fail(function () {
            var s = 'Save Property encountered an error';
            w2ui.propertyGrid.error(s);
            propertySaveDoneCB();
        });
}

function resetTabDispCounts() {
    propData.tabGenDispCount = 0;
    propData.tabRStDispCount = 0;
    propData.tabROpDispCount = 0;
    propData.tabTraDispCount = 0;
    propData.tabPhoDispCount = 0;
    propData.tabStaDispCount = 0;

}

function setPropertyNotLoaded() {
    propData.bPropLoaded = false;
    propData.bRentStepsLoaded = false;
    propData.bRenewOptionsLoaded = false;
    propData.bTrafficLoaded = false;
    propData.bStateLoaded = false;
    propData.states = [];
    resetTabDispCounts();
}

// displayDateString returns a string that can be used to populate a form field.
// If the date is prior to 1971 then it will use a blank string.
//------------------------------------------------------------------------------
function displayDateString(d) {
    var y = new Date(d);
    var s = "";
    if (y.getFullYear() > 1970) {
        s = dateFmtStr(y);
    }
    return s;
}

function loadPropertyForm(PRID) {
    var f = w2ui.propertyForm;
    var rec = null;


    for (var i = 0; i < w2ui.propertyGrid.records.length; i++) {
        if (w2ui.propertyGrid.records[i].PRID == PRID) {
            rec = w2ui.propertyGrid.records[i];
            break;
        }
    }
    if (rec === null) {
        console.log('could not find PRID = ' + PRID + ' in property grid records');
        return;
    }
    w2ui.propertyForm.recid = rec.PRID;
    propData.PRID = rec.PRID;
    closeStateChangeDialog();

    // setPropertyHeader();
    f.url = "/v1/property/" + rec.PRID;
    f.reload();  // get this going as quickly as possible

    setPropertyNotLoaded();
    var l = w2ui.propertyFormLayout.get('main');
    if (typeof l.tabs != "undefined") {
        if (typeof l.tabs.name == "string") {
            l.tabs.click('proptabState');
        }
    }
}

// getDropDownIndexSafe - the UI may not have been loaded if the user did not
//    goto the "General" tab.  When this is the case, just return the prev
//    value -- it's like the default value.
//
// s = name of UI element to check
// y = previous value
//------------------------------------------------------------------------------
function getDropDownSafe(s, y) {
    var x = getDropDownSelectedIndex(s);
    if (x > -1) {
        return x;
    }
    return y;
}

// prepareSaveData returns the data to push to the server to save the property
//------------------------------------------------------------------------------
function prepareSaveData() {
    var rec = w2ui.propertyForm.record;

    if (rec.Name == "") {
        rec.Name = "New Property";
    }

    //-----------------------------------------
    // Handle any conversions necessary...
    //-----------------------------------------
    rec.AvgCap /= 100;  // convert back to decimal number
    rec.CapRate /= 100; // convert back to decimal number
    rec.RentCommencementDt = varToUTCString(rec.RentCommencementDt);
    rec.LeaseExpirationDt = varToUTCString(rec.LeaseExpirationDt);
    rec.CreateTime = varToUTCString(rec.CreateTime);
    rec.LastModTime = varToUTCString(rec.LastModTime);

    rec.LotSizeUnits = getDropDownSafe("LotSizeUnitsDD", rec.LotSizeUnits);
    rec.OwnershipType = getDropDownSafe("OwnershipTypeDD", rec.OwnershipType);
    rec.LeaseType = getDropDownSafe("LeaseTypeDD", rec.LeaseType);
    rec.LeaseGuarantor = getDropDownSafe("LeaseGuarantorDD", rec.LeaseGuarantor);

    var mask = 1 << 3;
    var b = getDropDownSafe("OwnershipDD", ((rec.FLAGS & mask) == 0) ? 0 : 1);
    if (b === 0) {
        rec.FLAGS &= ~mask;
    } else {
        rec.FLAGS |= mask;
    }
    mask = 1 << 1;
    b = getDropDownSafe("RoofResponsibilityDD", ((rec.FLAGS & mask) == 0) ? 0 : 1);  // initialize to value at last download
    if (b === 0) {
        rec.FLAGS &= ~mask;
    } else {
        rec.FLAGS |= mask;
    }

    var x = rec.BuildYear;
    if (typeof x == "string") {
        rec.BuildYear = parseInt(x);
    }
    x = rec.RenovationYear;
    if (typeof x == "string") {
        rec.RenovationYear = parseInt(x);
    }
    x = rec.YearFounded;
    if (typeof x == "string") {
        rec.YearFounded = parseInt(x);
    }

    //-----------------------------------------
    // Now send it to the server
    //-----------------------------------------
    var params = {
        cmd: "save",
        record: rec,
        states: [0, 0, 0, 0, 0, 0, 0]
    };
    return params;
}

// savePropertyForm grabs all the data that is associated with the propertForm,
//      converts anything that needs attention and calls the server's save
//      function.
//------------------------------------------------------------------------------
function savePropertyForm() {
    var params = prepareSaveData();
    var dat = JSON.stringify(params);
    var url = '/v1/property/' + params.record.PRID;

    return $.post(url, dat, null, "json")
        .done(function (data) {
            // if (data.status === "success") {
            // }
            if (data.status === "error") {
                w2ui.propertyGrid.error('ERROR: ' + data.message);
            }
            w2ui.propertyForm.record.PRID = data.recid;  // this ensures that if a new property was saved we have a valid PRID in the struct
        })
        .fail(function (data) {
            w2ui.propertyGrid.error("Save RentableLeaseStatus failed. " + data);
        });
}

// savePropertyFormWithCB saves the current property and calls the supplied
// callback function. This function has two parameters:
//     data - data returned from post
//     done - boolean: true = success, false = post failed
//-----------------------------------------------------------------------------
function savePropertyFormWithCB(cbf) {
    var params = prepareSaveData();
    var dat = JSON.stringify(params);
    var url = '/v1/property/' + params.record.PRID;

    return $.post(url, dat, null, "json")
        .done(function (data) {
            cbf(data, true);
        })
        .fail(function (data) {
            cbf(data, false);
        });
}

function propertySaveDoneCB() {
    w2ui.toplayout.hide('right', true);
    w2ui.propertyGrid.reload();
}

// setPropertyLayout is used to display the property form in the UI and handle
// the tab clicking.
//
// INPUTS
//      PRID - int64, property id
//      tab  - string, name of the tab that was pressed
//------------------------------------------------------------------------------
function setPropertyLayout(tab) {
    w2ui.propertyFormLayout.html("bottom", w2ui.propertyFormBtns);


    switch (tab) {

        case "proptabState":
            propData.tabStaDispCount++;
            // w2ui.propertyFormLayout.load('main', '/static/html/formState.html', null,propertyStateOnLoad);
            w2ui.propertyStateLayout.load('main', '/static/html/formState.html', null, propertyStateOnLoad);
            w2ui.propertyFormLayout.html('main', w2ui.propertyStateLayout);
            setTimeout(propertyStateOnLoad, 100);
            break;

        case "proptabGeneral":
            propData.tabGenDispCount++;
            if (propData.bPropLoaded) {
                w2ui.propertyForm.url = '';
            } else {
                w2ui.propertyForm.url = '/v1/property/' + propData.PRID;
            }
            w2ui.propertyFormLayout.html('main', w2ui.propertyForm);
            setPropertyFormActionButtons(true);
            break;

        case "proptabRentSteps":
            propData.tabRStDispCount++;
            if (propData.bRentStepsLoaded) {
                w2ui.propertyRentStepsGrid.url = '';
            } else {
                if (w2ui.propertyForm.record.RSLID == 0) {
                    w2ui.propertyRentStepsGrid.url = '';
                    propData.bRentStepsLoaded = true;
                } else {
                    w2ui.propertyRentStepsGrid.clear();
                    w2ui.propertyRentStepsGrid.url = '/v1/rentsteps/' + w2ui.propertyForm.record.RSLID;
                }
            }
            w2ui.rentStepsLayout.html('main', w2ui.propertyRentStepsGrid);
            w2ui.propertyFormLayout.html('main', w2ui.rentStepsLayout);
            setPropertyFormActionButtons(true);
            break;

        case "proptabRenewOptions":
            propData.tabROpDispCount++;
            if (propData.bRenewOptionsLoaded) {
                w2ui.propertyRenewOptionsGrid.url = '';
            } else {
                if (w2ui.propertyForm.record.ROLID == 0) {
                    w2ui.propertyRenewOptionsGrid.url = '';
                    propData.bRenewOptionsLoaded = true;
                } else {
                    w2ui.propertyRenewOptionsGrid.clear();
                    w2ui.propertyRenewOptionsGrid.url = '/v1/renewoptions/' + w2ui.propertyForm.record.ROLID;
                }
            }
            w2ui.renewOptionsLayout.html('main', w2ui.propertyRenewOptionsGrid);
            w2ui.propertyFormLayout.html('main', w2ui.renewOptionsLayout);
            setPropertyFormActionButtons(true);
            break;

        case "proptabTraffic":
            propData.tabTraDispCount++;
            if (propData.bTrafficLoaded) {
                w2ui.propertyTrafficGrid.url = '';
            } else {
                w2ui.propertyTrafficGrid.clear();
                w2ui.propertyTrafficGrid.url = '/v1/trafficitems/' + propData.PRID;
            }
            w2ui.propertyTrafficLayout.html('main', w2ui.propertyTrafficGrid);
            w2ui.propertyFormLayout.html('main', w2ui.propertyTrafficLayout);
            setPropertyFormActionButtons(true);
            break;
        case "proptabPhotos":
            propData.tabPhoDispCount++;
            w2ui.propertyPhotosLayout.load('main', '/static/html/formPhotos.html');
            w2ui.propertyFormLayout.html('main', w2ui.propertyPhotosLayout);
            setTimeout(function () {
                setPropertyFormActionButtons(false);
            }, 100);
            break;
    }

    showForm();
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function showForm() {
    // SHOW the right panel now
    w2ui.toplayout.sizeTo('right', propData.formWidth);
    await sleep(1000);
    w2ui.toplayout.show('right', true);
    await sleep(1000);
    w2ui.toplayout.html('right', w2ui.propertyFormLayout);
}

// setPropertyFormActionButtons is used to turn off and on the Save, Save And Add, and
// Delete buttons at the bottom of the form.
//
// INPUTS:
//  t = true -> turn buttons on,  false -> turn buttons off
//--------------------------------------------------------------------------------
function setPropertyFormActionButtons(t) {
    var f = w2ui.propertyFormBtns;
    var x = !t;
    $(f.box).find("button[name=save]").prop("disabled", x);
    $(f.box).find("button[name=cancel]").prop("disabled", x);
}


function closePropertyForm() {
    w2ui.toplayout.hide('right', true);
    w2ui.propertyGrid.render();
}

// function setTermRemaining() {
//     var s = "n/a";
//     var now = new Date();
//     var s1=w2ui.propertyForm.record.RentCommencementDt;
//     if (s1 == null || typeof s1 != "string" ) {
//         setInnerHTML("PRTermRemaining",s);
//         return;
//     }
//     if (s1.length == 0) {
//         setInnerHTML("PRTermRemaining",s);
//         return;
//     }
//     var d1=new Date(s1);
//     var s2=w2ui.propertyForm.record.LeaseExpirationDt;
//     if (s2 == null || typeof s2 != "string") {
//         setInnerHTML("PRTermRemaining",s);
//         return;
//     }
//     if (s2.length == 0) {
//         setInnerHTML("PRTermRemaining",s);
//         return;
//     }
//     var d2=new Date(s2);
//     var m=monthDiff(d1,d2);
//     // Handle the case where the expiration date has passed...
//     if (now.getTime() > d2.getTime()) {
//         s = "Lease has expired";
//     } else {
//         s = Math.floor(m/12) + ' years ' + Math.floor(m%12) + ' months';
//     }
//     setInnerHTML("PRTermRemaining",s);
// }
