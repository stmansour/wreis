
function imageFilename(n) {
    var s = property["Img"+n];
    if (s == "") {
        return s;
    }
    return jb.cwd + "/Img" + n + "." + fileExtension(s);
}


//------------------------------------------------------------------------------
// fitItem  center an image on the page, resize to maintain aspect ratio
//
// item - the image
// p    - the size and location of the rectangle into which the image is fitted
//------------------------------------------------------------------------------
function fitItem(item, p) {
    var wratio = p.width / item.width;
    var hratio = p.height / item.height;
    var MINSIZE = 10;
    // alert('fitItem:  item w,h = ' + item.width + ', ' + item.height + '\np: w,h = ' + p.width + ', ' + p.height);
    if (hratio > wratio) {
        // landscape, scale height using ratio from width
        var newheight = (p.width * item.height) / item.width;
        if (newheight < MINSIZE) {
            return;
        }
        item.width = p.width;
        item.height = newheight;
    } else {
        // portrait, scale width using ratio from height
        var nw = (p.height * item.width) / item.height;
        if (nw < MINSIZE) {
            return;
        }
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

// placeImage - inserts the image with filename fname int area p
//
// layer     = layer to add image to
// n         = index number of the image (1..n)
// nameInAI  = name to give the new Illustrator item
// p         = the size and location of the rectangle into which the image is
//             fitted.  It is an object with members: left, top, width, height
//------------------------------------------------------------------------------
function placeImage(layer,n,nameInAI,p) {
    var placedItem = layer.placedItems.add();
    var fqname = imageFilename(n);  // fully qualified name
    if (fqname == "") {
        return;
    }
    try {
        placedItem.file = new File(fqname);
    } catch (error) {
        alert(fqname + ': ' + error);
        return;
    }
    placedItem.name = nameInAI;
    fitItem(placedItem,p);
}

// resizePlaceAndCrop -
//      resizes, centers, and crops the image to totally fill the area defined
//       by bb (bounding box)
//
// img       = the image to resize, center, and crop
// bb        = the object of the bounds. Members left, top, width, and height
//             define the useable area
// name      = name for the newly created image
//------------------------------------------------------------------------------
function resizePlaceAndCrop(img, bb, name) {
    var doc = app.activeDocument;
    //--------------------------------------------------------------------------
    // 1. scale it in such a way that the image fully covers the bounding area.
    // 2. center the scaled image in the bb
    // 3. crop it to fit the bounding box
    //--------------------------------------------------------------------------
    var wratio = bb.width / img.width;
    var hratio = bb.height / img.height;
    var bcx = bb.left + bb.width/2;
    var bcy = bb.top - bb.height/2;

    var MINSIZE = 10;
    if (hratio < wratio) {
        //---------------------------------------------------------------------
        // landscape, scale height using ratio from width
        //---------------------------------------------------------------------
        var newheight = (bb.width * img.height) / img.width;
        if (newheight < MINSIZE) {
            return;
        }
        img.width = bb.width;
        img.height = newheight;

        //---------------------------------------------------------------------
        // center the scaled image
        //---------------------------------------------------------------------
        var top = img.top;
        img.left = bb.left;
        icy = img.top - img.height/2;
        img.top = top + bcy - icy;

    } else {
        //---------------------------------------------------------------------
        // portrait, scale width using ratio from height
        //---------------------------------------------------------------------
        var nw = (bb.height * img.width) / img.height;
        if (nw < MINSIZE) {
            return;
        }
        img.height = bb.height;
        img.width = nw;

        //---------------------------------------------------------------------
        // center the scaled image
        //---------------------------------------------------------------------
        var left = img.left;
        icx = img.left + img.width/2;
        img.left = left + bcx - icx;
    }

    // alert('img.left = ' + img.left + ', img.top = ' + img.top + ', img.width = ' + img.width + ', img.height = ' + img.height);

    //---------------------------------------------------------------------
    // crop it
    //---------------------------------------------------------------------
    var rasterOpts = new RasterizeOptions();
    rasterOpts.antiAliasingMethod = AntiAliasingMethod.ARTOPTIMIZED; // the other option is TYPEOPTIMIZED
    rasterOpts.resolution = 72;
    var newimg = doc.rasterize(img, bb.geometricBounds, rasterOpts);
    newimg.name = name;
}

// placeResizeCenterCropImage - inserts the image with filename fname int area p. The
//              image will be expanded so that it is centered in the  space
//
// layer     = layer to add image to
// n         = index number of the image (1..n)
// nameInAI  = name to give the new Illustrator item
// bb        = the object of the bounds. members left, top, width, and height define
//             the useable area
//------------------------------------------------------------------------------
function placeResizeCenterCropImage(layer,n,nameInAI,bb) {
    var img = layer.placedItems.add();
    // img.name = nameInAI;

    var fqname = imageFilename(n);  // fully qualified name
    if (fqname == "") {
        return;
    }
    try {
        img.file = new File(fqname);
    } catch (error) {
        alert(fqname + ': ' + error);
        return;
    }
    resizePlaceAndCrop(img,bb,nameInAI);
}

// fillWithImage - fills an area with an image - will resize and crop image as
//                 needed to fill the area.
//
// INPUTS
// lname     = namd of layer to add image to
// rname     = name of rectangle (pathItem) defining area to fill
// n         = index number of the image (1..n)
// nameInAI  = name to give the new Illustrator item
//------------------------------------------------------------------------------
function fillWithImage(lname,rname,n,nameInAI) {
    var lyr = app.activeDocument.layers.getByName(lname);
    var bb = lyr.pathItems.getByName(rname);
    placeResizeCenterCropImage(lyr,n,nameInAI,bb);
}

// fitFullImageInPageItem  center an image on the page, resize to maintain aspect ratio
//
// item - the image
// p    - the size and location of the artboard (the page)
// hdr  - the name of the path defining the header of the page.  It is assumed
//        to be a rectangle located at the top of the artboard.
//------------------------------------------------------------------------------
function fitFullImageInPageItem(item, p, hdr) {
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
    fitItem(item,p);
}

// placeImageInArea
//
// n         - the image index number (1 .. n)
// imgAIName - the name to give the image in the AI file
// pName     - name of the path item that sets the bounds into which the image
//             will be placed
// layer     - layer into which image will be placed
//------------------------------------------------------------------------------
function placeImageInArea(n, imgAIName, pName, layer) {
    var b = jb.doc.pathItems.getByName(pName);
    if (null == layer) {
        alert("placeImageInArea:  layer could not be found");
        return;
    }
    placeImage(layer,n,imgAIName,b);
}

function placeCoverImage() {
    var fname = imageFilename(1);
    var placedItem = jb.doc.placedItems.add();
    try {
        placedItem.file = new File(fname);
    } catch (error) {
        alert(fname + ': ' + error);
        return;
    }
    placedItem.name = "coverPicture";
    var b = getArtboardBounds(jb.ab);
    fitFullImageInPageItem(placedItem,b,"coverPageHeaderBar");
}
