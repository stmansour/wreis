function arboardPositionString() {
    var s = '';
    for (var i = 0; i < app.activeDocument.artboards.length; i++) {
        s += '' + i + '. ' + app.activeDocument.artboards[i].name + '  [' + app.activeDocument.artboards[i].artboardRect + "]\n";
    }
    return s;
}

// layer = layer containing objects to move
// dx,dy = the offset to apply to each object
//------------------------------------------------------------------------------
function moveLayerObjects(layer,dx,dy) {
    for (i = 0; i < layer.pathItems.length; i++) {
        layer.pathItems[i].left += dx;
        layer.pathItems[i].top += dy;
    }
    for (i = 0; i < layer.placedItems.length; i++) {
        layer.placedItems[i].left += dx;
        layer.placedItems[i].top += dy;
    }
    for (i = 0; i < layer.textFrames.length; i++) {
        layer.textFrames[i].left += dx;
        layer.textFrames[i].top += dy;
    }
}

function findArtboardIndex(a) {
    for (var i = 0; i < app.activeDocument.artboards.length; i++) {
        if (a == app.activeDocument.artboards[i].name) {
            return i;
        }
    }
    return -1;
}

function setTextContents(pname,val) {
    var t = jb.doc.textFrames.getByName(pname);
    t.contents = val;
}

// returns the total number of photos that will be shown in the Subject Property
// area.
// -----------------------------------------------------------------------------
function subjectPropertyPhotoCount() {
    var n = 0;
    // cover photo
    if (property.Img0 != "") {
        n++;
    }
    // other photos
    var s;
    for (var j = jb.subjProp; j <= 12; j++) {
        s = "Img" + j;
        if (property[s] != "") {
            n++;
        }
    }
    return n;
}

function addSubjectImages() {
    var n = 0;
    app.selection = null;  // ensure nothing is selected
    var cpABR = app.activeDocument.artboards.getByName("Offering").artboardRect;

    for (var j = jb.subjProp; j <= 12; j++) {
        var s = "Img" + j;
        if (property[s] == "") {
            continue;
        }
        //-------------------------------------
        // Add an artboard for this image...
        //-------------------------------------
        var artb = app.activeDocument.artboards.getByName("Subject Property 1");
        var ab = artb.artboardRect;
        var idxSP1 = findArtboardIndex(artb.name);
        if (idxSP1 < 0) {
            alert("no index found for Subject Property 1 artboard!");
            return;
        }
        //                            0    1     2       3
        // artboardRect contains: [ left, top, right, bottom ]
        var width = ab[2] - ab[0];
        var dx = 36 + n*(36 + width);
        var x1 = ab[2] + dx;
        var x2 = x1 + width;
        var nabRect = [x1, ab[1], x2, ab[3]];  // this is where the new artboard goes
        var insertIndex = idxSP1 + n+1;
        app.activeDocument.artboards.insert(nabRect,insertIndex);
        var nab =  app.activeDocument.artboards[insertIndex];
        n += 1;
        nab.name = "Subject Property " + (n+1);
        var nabABR = nab.artboardRect;
        var layer = app.activeDocument.layers.add();
        layer.name = nab.name;
        var idxnab = findArtboardIndex(nab.name);

        //----------------------------------------------------------------
        // copy the header and footer from Subject Property 1.
        //   * set current artboard to Subject Property 1
        //   * Get Layer SubjectProperty1
        //   * mark it has objects selected
        //   * deselect the image named SubjectProperty1
        //----------------------------------------------------------------
        app.activeDocument.artboards.setActiveArtboardIndex(idxSP1);
        var sourceLayer = app.activeDocument.layers.getByName("Subject Property 1");  // page layer
        sourceLayer.hasSelectedArtwork = true;
        var img = sourceLayer.placedItems.getByName("SubjectProperty1");  // the image
        img.selected = false;
        app.copy();

        //----------------------------------------------------------------
        // Now paste, this will make new copies in sourceLayer.
        // Set current artboard to new artboard first...
        // Then we
        // have to move them out of the source layer into the new Subject
        // Property n layer we created above.
        //----------------------------------------------------------------
        app.activeDocument.artboards.setActiveArtboardIndex(idxnab);
        app.paste();
        var docSelected = app.activeDocument.selection;
        var anObj = null;
        for (s = 0; s < docSelected.length; s++) {
             anObj = docSelected[s];
             anObj.move(layer, ElementPlacement.PLACEATEND);
        }

        //----------------------------------------------------------------
        // All the objects were pasted in place based on the current
        // window size... it will always be different.
        // To get all objects where they need to be is a 2-step delta
        // calculation.
        // First compute the delta to line them up with the coverpage.
        // Second, move from cover page to Subject Property page location
        //----------------------------------------------------------------
        var pi = layer.pathItems.getByName("SP1-PageOutline");
        var ddx = cpABR[0] - pi.left;
        var ddy = cpABR[1] - pi.top;

        dx = ddx + nabABR[0] - cpABR[0];      // delta x from Cover page to new page
        var dy = ddy + nabABR[1] - cpABR[1];  // delta y from Cover page to new page
        moveLayerObjects(layer,dx,dy);

        //----------------------------------------------------------------
        // All the objects were pasted in place, so we need to move them
        // to the right.
        //----------------------------------------------------------------
        placeImageInArea(j,"SubjectProperty"+(j-4),"SP1-Background",app.activeDocument.layers.getByName("Subject Property "+(j-4)));
    }
}

//----------------------------------------------------------------------------
//  MAIN ROUTINE
//----------------------------------------------------------------------------


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
    var template = jb.cwd + '/res/template00.ai';
      var f = new File(template);
    var lyr;  // layer
    var bb;   // bounding box (any object with top,left,width,height )
    app.open(f);

    //---------------------------------------------------------------------------
    // immediately save this as a new document: portfolio.ai
    //---------------------------------------------------------------------------
    fname = jb.cwd + '/portfolio.ai';
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
    setTextContents("propertyName",property.Name);
    setTextContents("streetAddress",property.Address);
    setTextContents("cityStateZip",property.City + ", " + property.State + "  " + property.PostalCode);
    //placeCoverImage();
    if (property.Img1 != "") {
        fillWithImage("COVER PAGE - Offering","CP-CoverPhotoArea",1,"CoverShot");
    }
    if (property.Img4 != "") {
        placeImageInArea(4,"CoverLogo","CP-LogoArea",jb.doc.layers.getByName("COVER PAGE - Offering"));
    }

    //---------------------------------------------------------------------------
    //  PAGE 2 - Table of Contents
    //---------------------------------------------------------------------------
    var sp = subjectPropertyPhotoCount();
    var s = "9";
    var pn = 9;
    if (sp > 1) {
        pn += sp - 1;
        s += " - " + pn;
    }
    setTextContents("Page-Property Photos",s);
    setTextContents("TOC-MarketOverview",pn + 1);
    setTextContents("TOC-DemographicReport",pn + 2);
    if (property.Img4 != "") {
        placeImageInArea(4,"TOCLogo","TOC-LogoArea",jb.doc.layers.getByName("TOC - Table Of Contents"));
    }

    //---------------------------------------------------------------------------
    //  PAGE 3 - Financial Overview
    //---------------------------------------------------------------------------
    var t = jb.doc.textFrames.getByName("FO-Price");
    var dt = new Date();
    jb.chattr = t.textRange.characterAttributes;    // we save this for use later
    setTextContents("FO-Price",fmtCurrency(property.Price));
    setTextContents("FO-DownPayment", "100% / " + fmtCurrency(property.DownPayment));
    setTextContents("FO-RentableSF",fmtWithCommas(property.RentableArea) + " SF");
    setTextContents("FO-BuildRenovationYear",property.BuildYear + ((property.RenovationYear > 0) ? ' / ' + property.RenovationYear : "") );
    fmtIndexedName(((property.FLAGS & (1<<1)) == 0) ? 0 : 1,"FO-Roof",jb.roofStructureLabels,"roof structure flag");
    setTextContents("FO-RightOfFirstRefusal", (property.FLAGS & (1<<2) > 0) ? "Yes" : "No");
    setTextContents("FO-CapRate",fmtAsPercent(property.CapRate));
    setTextContents("FO-LeaseTermRemaining",fmtDateDiffInYears(dt, property.LeaseExpirationDt,property.OriginalLeaseTerm));
    setTextContents("FO-LotSize", (property.LotSizeUnits + 1 > jb.lotSizeLabels.length) ? "(unknown units)" : fmtAsFloatWithCommas(property.LotSize) + ' ' + jb.lotSizeLabels[property.LotSizeUnits]);
    fmtIndexedName(property.OwnershipType,"FO-TypeOwnership",jb.ownershipTypeLabels,"ownership type");
    setTextContents("FO-TenantTradeName",property.TenantTradeName);
    fmtIndexedName(property.LeaseGuarantor,"FO-LeaseGuarantor",jb.guarantorLabels,"guarantor");
    fmtIndexedName(property.LeaseType,"FO-LeaseType",jb.leaseTypeLabels,"lease type");
    setTextContents("FO-OriginalLeaseTerm",property.OriginalLeaseTerm + " years");
    fmtDate(property.LeaseExpirationDt, "FO-LeaseExpirationDate");
    fmtIndexedName(((property.FLAGS & (1<<3)) == 0) ? 0 : 1,"FO-Ownership",jb.ownershipLabels,"ownership type");
    genTable();

    //---------------------------------------------------------------------------
    //  PAGE 4 - Tenant Overview
    //---------------------------------------------------------------------------
    setTextContents("TO-TenantTradeName",property.TenantTradeName);
    setTextContents("TO-PropertyName",property.Name);
    setTextContents("TO-PropertyAddressLine1",property.Address);
    setTextContents("TO-PropertyAddressLine2",property.City + ", " + property.State + "  " + property.PostalCode);
    setTextContents("TO-LeaseTermRemaining",fmtDateDiffInYears(dt, property.LeaseExpirationDt,property.OriginalLeaseTerm));
    fmtIndexedName(((property.FLAGS & (1<<3)) == 0) ? 0 : 1,"TO-Ownership",jb.ownershipLabels,"ownership type");
    setTextContents("TO-Company",property.TenantTradeName);
    fmtIndexedName(property.LeaseGuarantor,"TO-LeaseGuarantor",jb.guarantorLabels,"guarantor");
    setTextContents("TO-StockSymbol",property.Symbol);
    setTextContents("TO-Year1NOI",fmtCurrency(property.NOI));
    s = " ";
    if (typeof property.RenewOptions != "undefined") {
        s = "(" + property["renewOptions"].length + ")";
    }
    setTextContents("TO=OptionsToRenew",s);
    fmtIndexedName(property.LeaseType,"TO-LeaseType",jb.leaseTypeLabels,"lease type");
    fmtIndexedName(((property.FLAGS & (1<<1)) == 0) ? 0 : 1,"TO-RoofStructure",jb.roofStructureLabels,"roof structure flag");
    t = jb.doc.textFrames.getByName("TO-Headquarters");
    t.contents = property.HQCity + "," + property.HQState;
    t = jb.doc.textFrames.getByName("TO-Website");
    t.contents = property.URL;
    t = jb.doc.textFrames.getByName("TO-YearsInTheBusiness");
    var now = new Date();
    // t.contents = (property.YearFounded > 0) ? "" + now.getFullYear() - property.YearFounded : " ";
    t.contents = (property.YearFounded > 0) ? "" + property.YearFounded : " ";

    if (property.Img4 != "") {
        placeImageInArea(4,"TOLogo","TO-LogoArea",jb.doc.layers.getByName("Tenant Overview"));
    }

    //---------------------------------------------------------------------------
    //  PAGE 5 - Executive Summary
    //---------------------------------------------------------------------------

    //---------------------------------------------------------------------------
    //  PAGE 6 - Aerial Photo
    //---------------------------------------------------------------------------
    if (property.Img2 != "") {
        fillWithImage("Aerial Photo","AP-background",2,"AP-AerialPhoto");
    }

    //---------------------------------------------------------------------------
    //  PAGE 7 - Subject Property
    //
    //  These start with the cover image (Img1.)+fileExtension(property.Img1)
    //  and will include images 6 - 8 if present.
    //---------------------------------------------------------------------------
    if (property.Img1 != "") {
        placeImageInArea(1,"SubjectProperty1","SP1-Background",jb.doc.layers.getByName("Subject Property 1"));
        // fillWithImage("Subject Property 1","SP1-Background",1,"SubjectProperty1");
    }
    addSubjectImages();

    //---------------------------------------------------------------------------
    //  PAGE 7+n - Area Map  (n = total count of subject photos + cover photo)
    //---------------------------------------------------------------------------
    if (property.Img3 != "") {
        placeImageInArea(3,"AM-AreaMap","AM-Background",jb.doc.layers.getByName("Area Map"));
        // fillWithImage("Area Map","AM-Background",3,"AM-AreaMap");
    }

    //---------------------------------------------------------------------------
    //  PAGE 8+n - Location Map  (n = total count of subject photos + cover photo)
    //---------------------------------------------------------------------------
    if (property.Img3 != "") {
        placeImageInArea(5,"LM-LocationMap","LM-Background",jb.doc.layers.getByName("Location Map"));
        // fillWithImage("Location Map","LM-Background",5,"LM-LocationMap");
    }

}

generateMarketPackage();
