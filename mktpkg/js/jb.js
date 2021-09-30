

//----------------------------------------------------------------------------
//  MAIN ROUTINE
//----------------------------------------------------------------------------

function getArtboardBounds(artboard) {
    var bounds = artboard.artboardRect;
    var left = bounds[0];
    var top = bounds[1];
    var right = bounds[2];
    var bottom = bounds[3];
    var width = right - left;
    var height = top - bottom;
    var b = {
        left: left,
        top: top,
        width: width,
        height: height
    };
    // alert('artboard bounds:  top=' + b.top + ' left=' + b.left + ' width=' + b.width + ' height=' + b.height );
    return b;
}

// fitItem  center an image on the page, resize to maintain aspect ratio
//
// item - the image
// p    - the size and location of the artboard (the page)
// hdr  - the name of the path defining the header of the page.  It is assumed
//        to be a rectangle located at the top of the artboard.
//------------------------------------------------------------------------------
function fitItem(item, p, hdr) {
    var bar = jb.doc.pathItems.getByName(hdr);
    var bt = bar.visibleBounds;
    var b = {
        left: bt[0],
        top:  bt[1],
        right: bt[2],
        bottom: bt[3],
        width: bt[2] - bt[0],
        height: bt[1] - bt[3],
    };

    //-----------------------------------------------------
    // reduce values in p to account for the top bar...
    //-----------------------------------------------------
    p.top -= b.height;
    p.height -= b.height;

    if (item.width > item.height) {
        // landscape, scale height using ratio from width
        var newheight = (p.width * item.height) / item.width;
        item.width = p.width;
        item.height = newheight;
    } else {
        // portrait, scale width using ratio from height
        var nw = (p.height * item.width) / item.height;
        item.height = p.height;
        item.width = nw;
    }

    var cx = p.left + (p.width/2);
    var cy = p.top - (p.height/2);
    var il = cx - (item.width/2);
    var it = cy + (item.height/2);

    // alert('cx='+cx+' cy='+cy+' il='+il+' it='+it);

    item.top = it;
    item.left = il;
    item.selected = false;
}

function placeCoverImage() {
    var fname = jb.cwd + "/Img1.png";
    var placedItem = jb.doc.placedItems.add();
    try {
        placedItem.file = new File(fname);
    } catch (error) {
        alert(fname + ': ' + error);
        return;
    }

    placedItem.name = "coverPicture";
    var b = getArtboardBounds(jb.ab);
    fitItem(placedItem,b,"coverPageHeaderBar");
}

function placeAerialImage() {
    var placedItem = jb.doc.placedItems.add();
    var fname = jb.cwd + "/Img2.png";
    try {
        placedItem.file = new File(fname);
    } catch (error) {
        alert(fname + ': ' + error);
        return;
    }
    placedItem.name = "aerialPhoto";

    var aab = jb.doc.artboards.getByName("Aerial Photo");
    if (aab == null) {
        alert("artboard not found:  Aerial Photo");
        return;
    }

    var b = getArtboardBounds(aab);
    fitItem(placedItem,b,"aerialPhotoHeaderBar");
}

function generateMarketPackage() {
    //--------------------------------------------------------------
    // Make sure to close any currently open document...
    //--------------------------------------------------------------
    if (!app.homeScreenVisible) {
        app.activeDocument.close(SaveOptions.PROMPTTOSAVECHANGES);
    }

    //---------------------------------------------------------------------------
    // By convention, we will keep templates in ~/Documents/wreis.
    // We get the myDocuments folder from the Adobe environment which maps it
    // to the file system correctly even on Windows.  We will look for a folder
    // named ~/Documents/wreis and open the file template00.ai
    //---------------------------------------------------------------------------
    var template = 'template00.ai';
    var fname = Folder.myDocuments + '/wreis/' + template;
    var f = new File(fname);
    app.open(f);

    //---------------------------------------------------------------------------
    // immediately save this as a new document: portfolio.ai
    //---------------------------------------------------------------------------
    fname = Folder.myDocuments + '/wreis/portfolio.ai';
    var portfolio = new File(fname);
    app.activeDocument.saveAs(portfolio);

    //---------------------------------------------------------------------------
    // Set to the first artboard and start the update process...
    //---------------------------------------------------------------------------
    jb.doc = app.activeDocument;
    jb.doc.artboards.setActiveArtboardIndex(0);  // we'll start on the offering
    jb.ab = jb.doc.artboards[0];

    //---------------------------------------------------------------------------
    //  PAGE 1 - Cover Page
    //---------------------------------------------------------------------------
    var t = jb.doc.textFrames.getByName("propertyName");
    t.contents = property.Name;
    t = jb.doc.textFrames.getByName("streetAddress");
    t.contents = property.Address;
    t = jb.doc.textFrames.getByName("cityStateZip");
    t.contents = property.City + ", " + property.State + "  " + property.PostalCode;
    placeCoverImage();

    //---------------------------------------------------------------------------
    //  PAGE 3 - Financial Overview
    //---------------------------------------------------------------------------
    t = jb.doc.textFrames.getByName("FO-Price");
    t.contents = fmtCurrency(property.Price);
    t = jb.doc.textFrames.getByName("FO-DownPayment");
    t.contents = fmtCurrency(property.DownPayment);
    t = jb.doc.textFrames.getByName("FO-RentableSF");
    t.contents = fmtWithCommas(property.RentableArea);
    t = jb.doc.textFrames.getByName("FO-Roof");
    if (property.FLAGS & 0x1 > 0) {
        t.contents = "Landlord Responsible";
    } else {
        t.contents = "Tenant Responsible";
    }
    t = jb.doc.textFrames.getByName("FO-RightOfFirstRefusal");
    if (property.FLAGS & 0x4 > 0) {
        t.contents = "Yes";
    } else {
        t.contents = "No";
    }

    t = jb.doc.textFrames.getByName("FO-CapRate");
    t.contents = fmtAsPercent(property.CapRate);

    t = jb.doc.textFrames.getByName("FO-LeaseTermRemaining");
    var dt = new Date();
    t.contents = fmtDateDiffInYears(dt, property.LeaseExpirationDt);


    t = jb.doc.textFrames.getByName("FO-BuildRenovationYear");
    if (property.RenovationYear > 0) {
        t.contents = '' + property.RenovationYear;
    } else {
        t.contents = 'n/a';
    }
    t = jb.doc.textFrames.getByName("FO-LotSize");
    if (property.LotSizeUnits + 1 > jb.lotSizeLabels.length) {
        t.contents = "(unknown units)";
    } else {
        t.contents = fmtWithCommas(property.LotSize) + ' ' + jb.lotSizeLabels[property.LotSizeUnits];
    }
    fmtIndexedName(property.OwnershipType,"FO-TypeOwnership",jb.ownershipTypeLabels,"ownership type");

    //---------------------------------------------------------------------------
    //  PAGE 4 - Financial Overview
    //---------------------------------------------------------------------------
    t = jb.doc.textFrames.getByName("FO-TenantTradeName");
    t.contents = property.TenantTradeName;
    fmtIndexedName(property.LeaseGuarantor,"FO-LeaseGuarantor",jb.guarantorLabels,"guarantor");
    fmtIndexedName(property.LeaseType,"FO-LeaseType",jb.leaseTypeLabels,"lease type");
    t = jb.doc.textFrames.getByName("FO-OriginalLeaseTerm");
    t.contents = property.OriginalLeaseTerm + " years";
    fmtDate(property.LeaseExpirationDt, "FO-LeaseExpirationDate");
    var own = ((property.FLAGS & (1<<3)) == 0) ? 0 : 1;
    fmtIndexedName(own,"FO-Ownership",jb.ownershipLabels,"ownership type");

    //---------------------------------------------------------------------------
    //  PAGE 4 - Financial Overview
    //---------------------------------------------------------------------------

    //---------------------------------------------------------------------------
    //  PAGE 5 - Executive Summary
    //---------------------------------------------------------------------------

    //---------------------------------------------------------------------------
    //  PAGE 6 - Aerial Photo
    //---------------------------------------------------------------------------
    placeAerialImage();
}

generateMarketPackage();
