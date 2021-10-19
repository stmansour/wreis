//
//  jbx.js - the portfolio writer :-)
//
//  File -> Scripts -> Other Script...
//     or press Cmd + F12
//  then select this file.
//=========================================================================
var jb = {
    portfolio: null,        // the portfolio.ai we are auto-generating
    ab: null,               // active artboard
    doc: null,              // the working document
    chattr: null,           // the default font and attributes
    cwd: "",                // the current working directory
    lotSizeLabels: [        // what units for LotSize
        "sqft", "acres"
        ],
    ownershipTypeLabels: [      // OwnershipTypetype
        "Fee Simple",
        "Leasehold"
        ],
    ownershipLabels: [
        "Private",
        "Public"
        ],
    roofStructureLabels: [      // roof responsibility
        "Tenant Responsible",
        "Landlord Responsible"
        ],
    guarantorLabels: [      // who is guarantor
        "Corporate",
        "Franchise",
        "Individual"
        ],
    leaseTypeLabels: [
        "Absolute NNN",
        "Double Net",
        "Triple Net",
        "Gross"
    ],
};
