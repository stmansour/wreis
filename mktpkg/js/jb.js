

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

// item = the image
// p = the size and location of the artboard (the page)
function fitItem(item, p) {
    var bar = jb.doc.pathItems.getByName("headerBar");
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

function placeImage1() {
    var placedItem = jb.doc.placedItems.add();
    placedItem.file = new File(jb.cwd + "/Img1.png");
    placedItem.name = "coverPicture";
    var b = getArtboardBounds(jb.ab);
    fitItem(placedItem,b);
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
    //  FIRST PAGE
    //---------------------------------------------------------------------------
    var t = jb.doc.textFrames.getByName("propertyName");
    t.contents = property.Name;
    t = jb.doc.textFrames.getByName("streetAddress");
    t.contents = property.Address;
    t = jb.doc.textFrames.getByName("cityStateZip");
    t.contents = property.City + ", " + property.State + "  " + property.PostalCode;

    placeImage1();
}

generateMarketPackage();
