function addSubjectImages() {
    var n = 0;
    app.selection = null;  // ensure nothing is selected

    for (var j = 5; j <= 8; j++) {
        var s = "Img" + j;
        if (property[s] == "") {
            continue;
        }
        //-------------------------------------
        // Add an artboard for this image...
        //-------------------------------------
        var ab = jb.doc.artboards.getByName("Subject Property 1").artboardRect;
        //                            0    1     2       3
        // artboardRect contains: [ left, top, right, bottom ]
        var width = ab[2] - ab[0];
        var dx = 36 + n*(36 + width);
        var x1 = ab[2] + dx;
        var x2 = x1 + width;
        var nabRect = [x1, ab[1], x2, ab[3]];  // this is where the new artboard goes
        var nab =  jb.doc.artboards.add(nabRect);
        n += 1;
        nab.name = "Subject Property " + (n+1);
        var layer = jb.doc.layers.add();
        layer.name = nab.name;

        //----------------------------------------------------------------
        // copy the header and footer from Subject Property 1.
        //   * Get Layer SubjectProperty1
        //   * mark it has objects selected
        //   * deselect the image named SubjectProperty1
        //----------------------------------------------------------------
        var sourceLayer = jb.doc.layers.getByName("Subject Property 1");  // page layer
        sourceLayer.hasSelectedArtwork = true;
        var img = sourceLayer.placedItems.getByName("SubjectProperty1");  // the image
        img.selected = false;
        app.copy();

        //----------------------------------------------------------------
        // Now paste, this will make new copies in sourceLayer.  Then we
        // have to move them out of the source layer into the new Subject
        // Property n layer we created above.
        //----------------------------------------------------------------
        app.paste();
        var docSelected = app.activeDocument.selection;
        var anObj = null;
        for (s = 0; s < docSelected.length; s++) {
             anObj = docSelected[s];
             anObj.move(layer, ElementPlacement.PLACEATEND);
        }

        //----------------------------------------------------------------
        // All the objects were pasted in place, so we need to move them
        // to the right.
        //----------------------------------------------------------------
        dx = n*(37 + width);
        for (i = 0; i < layer.pathItems.length; i++) {
            layer.pathItems[i].left += dx;
        }
        for (i = 0; i < layer.placedItems.length; i++) {
            layer.placedItems[i].left += dx;
        }
        for (i = 0; i < layer.textFrames.length; i++) {
            layer.textFrames[i].left += dx;
        }

        //----------------------------------------------------------------
        // All the objects were pasted in place, so we need to move them
        // to the right.
        //----------------------------------------------------------------
        placeImageInArea("Img"+j+".png","SubjectProperty"+(j-3),"SP1-Background",jb.doc.layers.getByName("Subject Property "+(j-3)));
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
    jb.chattr = t.textRange.characterAttributes;    // we save this for use later

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
    t = jb.doc.textFrames.getByName("FO-TenantTradeName");
    t.contents = property.TenantTradeName;
    fmtIndexedName(property.LeaseGuarantor,"FO-LeaseGuarantor",jb.guarantorLabels,"guarantor");
    fmtIndexedName(property.LeaseType,"FO-LeaseType",jb.leaseTypeLabels,"lease type");
    t = jb.doc.textFrames.getByName("FO-OriginalLeaseTerm");
    t.contents = property.OriginalLeaseTerm + " years";
    fmtDate(property.LeaseExpirationDt, "FO-LeaseExpirationDate");
    var own = ((property.FLAGS & (1<<3)) == 0) ? 0 : 1;
    fmtIndexedName(own,"FO-Ownership",jb.ownershipLabels,"ownership type");
    genTable();

    //---------------------------------------------------------------------------
    //  PAGE 4 - Tenant Overview
    //---------------------------------------------------------------------------
    t = jb.doc.textFrames.getByName("TO-TenantTradeName");
    t.contents = property.TenantTradeName;
    t = jb.doc.textFrames.getByName("TO-PropertyName");
    t.contents = property.Name;
    t = jb.doc.textFrames.getByName("TO-PropertyAddressLine1");
    t.contents = property.Address;
    t = jb.doc.textFrames.getByName("TO-PropertyAddressLine2");
    t.contents = property.City + ", " + property.State + "  " + property.PostalCode;
    own = ((property.FLAGS & (1<<3)) == 0) ? 0 : 1;
    fmtIndexedName(own,"TO-Ownership",jb.ownershipLabels,"ownership type");
    t = jb.doc.textFrames.getByName("TO-ParentCompany");
    t.contents = property.ParentCompany;
    fmtIndexedName(property.LeaseGuarantor,"TO-LeaseGuarantor",jb.guarantorLabels,"guarantor");
    t = jb.doc.textFrames.getByName("TO-StockSymbol");
    t.contents = property.Symbol;
    t = jb.doc.textFrames.getByName("TO=OptionsToRenew");
    t.contents = "(" + property["renewOptions"].length + ")";
    fmtIndexedName(property.LeaseType,"TO-LeaseType",jb.leaseTypeLabels,"lease type");
    own = ((property.FLAGS & (1<<1)) == 0) ? 0 : 1;
    fmtIndexedName(own,"TO-RoofStructure",jb.roofStructureLabels,"roof structure flag");
    t = jb.doc.textFrames.getByName("TO-Headquarters");
    t.contents = property.HQCity + "," + property.HQState;
    t = jb.doc.textFrames.getByName("TO-Website");
    t.contents = property.URL;
    t = jb.doc.textFrames.getByName("TO-YearsInTheBusiness");
    var now = new Date();
    t.contents = (property.YearFounded > 0) ? "" + now.getFullYear() - property.YearFounded : " ";

    //---------------------------------------------------------------------------
    //  PAGE 5 - Executive Summary
    //---------------------------------------------------------------------------

    //---------------------------------------------------------------------------
    //  PAGE 6 - Aerial Photo
    //---------------------------------------------------------------------------
    placeAerialImage();

    //---------------------------------------------------------------------------
    //  PAGE 7 - Area Map
    //---------------------------------------------------------------------------
    placeImageInArea("Img3.png","AM-AreaMap","AM-Background",jb.doc.layers.getByName("Area Map"));

    //---------------------------------------------------------------------------
    //  PAGE 8 - Subject Property
    //
    //  These start with the cover image (Img1.png) and will include images
    //  5 - 8 if present.
    //---------------------------------------------------------------------------
    placeImageInArea("Img1.png","SubjectProperty1","SP1-Background",jb.doc.layers.getByName("Subject Property 1"));
    addSubjectImages();

}

generateMarketPackage();
