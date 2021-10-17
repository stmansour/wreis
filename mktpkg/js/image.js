// fitItem  center an image on the page, resize to maintain aspect ratio
//
// item - the image
// p    - the size and location of the rectangle into which the image is fitted
//------------------------------------------------------------------------------
function fitItem(item, p) {
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

// placeImage - inserts the image with filename fname int area p
//
// layer     = layer to add image to
// fname     = name of image file in the jb.cwd directory
// nameInAI  = name to give the new Illustrator item
// p         = the size and location of the rectangle into which the image is
//             fitted.  It is an object with members: left, top, width, height
//------------------------------------------------------------------------------
function placeImage(layer,fname,nameInAI,p) {
    var placedItem = layer.placedItems.add();
    var fqname = jb.cwd + "/" + fname;  // fully qualified name
    try {
        placedItem.file = new File(fqname);
    } catch (error) {
        alert(fqname + ': ' + error);
        return;
    }
    placedItem.name = nameInAI;
    fitItem(placedItem,p);
}


// fitFullPageItem  center an image on the page, resize to maintain aspect ratio
//
// item - the image
// p    - the size and location of the artboard (the page)
// hdr  - the name of the path defining the header of the page.  It is assumed
//        to be a rectangle located at the top of the artboard.
//------------------------------------------------------------------------------
function fitFullPageItem(item, p, hdr) {
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
// imgFName  - the image file name
// imgAIName - the name to give the image in the AI file
// pName     - name of the path item that sets the bounds into which the image
//             will be placed
// layers    - an array of layer names. The image will be placed in the last
//------------------------------------------------------------------------------
function placeImageInArea(imgFName, imgAIName, pName, layer) {
    var b = jb.doc.pathItems.getByName(pName);
    if (null == layer) {
        alert("placeImageInArea:  layer could not be found");
        return;
    }
    placeImage(layer,imgFName,imgAIName,b);
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
    fitFullPageItem(placedItem,b,"coverPageHeaderBar");
}

function placeAerialImage() {
    var layer = jb.doc.layers.getByName("Aerial Photo");
    var placedItem = layer.placedItems.add();
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
    fitFullPageItem(placedItem,b,"aerialPhotoHeaderBar");
}
