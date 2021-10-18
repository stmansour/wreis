//
//  jbx.js - the portfolio writer :-)
//
//  File -> Scripts -> Other Script...
//     or press Cmd + F12
//  then select this file.
//
//  ARTBOARD RECTANGLES:
//  1.      Contents            [ 812,  1952, 1604,  1340]
//  2.      Financial Overview  [ 812,  1304, 1604,   692]
//  3.      Tenant Overview     [ 812,   656, 1604,    44]
//  4.      Exectutive Summary  [ 812,     8, 1604,  -604]
//  5.      Aerial Photo        [ 812,  -640, 1604, -1252]
//  6.      Area Map            [ 812, -1288, 1604, -1900]
//  7.      Subject Property 1  [1640, -1288, 2432, -1900]
//  8.      Market Overview     [ 812, -1936, 1604, -2548]
//  9.      Demographic Report  [ 812, -2584, 1604, -3196]
//  10.     Closing Cover       [ 812, -3232, 1604, -3844]
//  11.     Subject Property 2  [2468, -1288, 3260, -1900]
//  12.     Subject Property 3  [3296, -1288, 4088, -1900]
//  13.     Subject Property 4  [4124, -1288, 4916, -1900]

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
