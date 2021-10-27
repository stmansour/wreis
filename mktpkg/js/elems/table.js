// genRect creates a rectangle and adds it to pathItems.
//            AOD = Annualized Operating Data
//
//    l   = layer
//    x,y = top left corner of rectangle
//    w,h = width and height of the rectangle
//    sw  = stroke width in points. if 0 then no stroke
//    sc  = stroke color. null means no stroke
//    fc  = fill color. If null, don set fill color
//-----------------------------------------------------------------------------
function genRect(l,x,y,w,h,sw,sc,fc) {
    var r = l.pathItems.add();
    r.filled = false;
    r.stroked = false;
    if (sw > 0 && sc != null) {
        r.stroked = true;
        r.strokeWidth = sw; // in points
        r.strokeColor = sc;
    }
    if (fc != null) {
        r.fillColor = fc;
        r.filled = true;
    }

    var t;
    var x1 = x;
    var x2 = x+w;
    var y1 = y;
    var y2 = y-h;

    if (x1 > x2) {
        t = x2;
        x2 = x1;
        x1 = t;
    }
    if (y1 > y2) {
        t = y2;
        y2 = y1;
        y1 = t;
    }

    r.setEntirePath([[x1,y1],[x2,y1],[x2,y2],[x1,y2],[x1,y1]]);
}

// genAODRect creates a rectangle and adds it to pathItems.
//            AOD = Annualized Operating Data
//            This routine sets the color for the tables
//
//    l   = layer
//    x,y = top left corner of rectangle
//    w,h = width and height of the rectangle
//    f   = boolean for filled or not
//-----------------------------------------------------------------------------
function genAODRect(l,x,y,w,h,f) {
    var strokeColor = aiGenColor(0x173e35);
    var fillColor = f ? aiGenColor(0xcdd1ce) : null;
    var strokeWidth = 0.5;
    genRect(l,x,y,w,h,strokeWidth,strokeColor,fillColor);
}

// tableext creates text
//
//    l     = layer
//    x,y   = top left corner of text
//    s     = string
//    sz    = text size
//    fname = fontname
//    sw    = stroke weight. if 0, stroke color is ignored
//    sc    = stroke color.  null if you don't want to set the color
//    fc    = fill color
//-----------------------------------------------------------------------------
function tableText(l,x,y,s,sz,fname,sw,sc,fc) {
    if (jb.chattr == null) {
        alert('jb.chattr is being used before it has been set!');
        return;
    }
    var t = l.textFrames.add();

    t.contents = s;
    t.textRange.characterAttributes = jb.chattr;
    t.textRange.characterAttributes.size = sz;
    t.textRange.characterAttributes.strokeWeight = sw;
    if (sw > 0 && sc != null) {
        t.textRange.characterAttributes.strokeColor = sc;
    }
    if (fc != null) {
        t.textRange.characterAttributes.fillColor = fc;
    }

    var font = app.textFonts.getByName(fname);
    if (font != null) {
        t.textRange.characterAttributes.textFont = font;
    }

    // alert("x,y = " + x + ', ' + y);
    t.textRange.justification = Justification.LEFT;
    t.position = [x,y];
}
// tableAODText creates text
//
//    l   = layer
//    x,y = top left corner of text
//    s   = string
//-----------------------------------------------------------------------------
function tableAODText(l,x,y,s) {
    var c = aiGenColor(0x000000);
    tableText(l,x,y,s,11,"ArialNarrow",0,null,c);
}


// genTable
//
//-----------------------------------------------------------------------------
function genTable() {
    var hbn = "AnnOpDt-Hdr";
    var MINY = 750;
    var l1 = jb.doc.layers.getByName("Financial Overview");
    if (l1==null) {
        alert("could not get layer Financial Overview" );
        return;
    }
    var layer = l1.layers.getByName("Annualized Operating Data");
    if (layer == null) {
        alert("could not get sublayer Annualized Operating Data" );
        return;
    }

    //------------------------------------------------------
    // First, select the header bar and get its location
    //------------------------------------------------------
    var p = jb.doc.pathItems.getByName(hbn);
    if (p == null) {
        alert(hbn + " was not found");
        return;
    }
    var width = p.width;
    var height = p.height;
    var left = p.left;
    var top = p.top;

    // alert('(' + left.toFixed(2) + ',' + top.toFixed(2) + '),' + '[' + width.toFixed(2) + ',' + height.toFixed(2) + ']');

    //------------------------------------------------------
    // Generate the rows...
    //------------------------------------------------------
    if (typeof property.renewOptions == "undefined" && typeof property.rentSteps == "undefined") {
        return;
    }

    var nrows = 0;
    if (typeof property.renewOptions != "undefined") {
        nrows += property.renewOptions.length;
    }
    if (typeof property.rentSteps != "undefined") {
        nrows += property.rentSteps.length;
    }
    if (nrows == 0) {
        return;
    }

    var y = top - height;
    var x = left;
    var offsetx = 10;   // from the left
    var offsety = 3;    // from the top
    var col2 = 175;     // 150 from left edge
    var col3 = 310;     // from left edge
    var monthly = 0.0;  // monthly amount
    var fill = false;   // every other rect should be filled
    var i;

    //---------------------------------
    // Do the Rent Steps first...
    //---------------------------------
    if (typeof property.rentSteps != "undefined") {
        for (i = 0; i < property.rentSteps.length; i++) {
            genAODRect(layer,x,y,width,height,fill);
            tableAODText(layer,x+offsetx, y-offsety, property.rentSteps[i].Opt);
            tableAODText(layer,x+offsetx+col2, y-offsety, fmtCurrency(property.rentSteps[i].Rent));
            monthly = property.rentSteps[i].Rent / 12.0;
            tableAODText(layer,x+offsetx+col3, y-offsety, fmtCurrency(monthly));
            y -=height;
            fill = !fill;
        }
    }
    //---------------------------------
    // Now the Renewal Options
    //---------------------------------
    if (typeof property.renewOptions != "undefined") {
        for (i = 0; i < property.renewOptions.length; i++) {
            genAODRect(layer,x,y,width,height,fill);
            tableAODText(layer,x+offsetx, y-offsety, property.renewOptions[i].Opt);
            tableAODText(layer,x+offsetx+col2, y-offsety, fmtCurrency(property.renewOptions[i].Rent));
            monthly = property.renewOptions[i].Rent / 12.0;
            tableAODText(layer,x+offsetx+col3, y-offsety, fmtCurrency(monthly));
            y -=height;
            fill = !fill;
        }
    }

    //---------------------------------
    // Now the skeleton table
    //---------------------------------
    var skx = x;
    var sky = y - height;
    var upperLX = skx;
    var upperLY = sky;
    var w = width/3;

    if (sky < MINY) {
        // alert("*** JB WARNING ***  not enough room for skeleton table or image on Financial Overview");
        return;
    }

    var strokeWidth = 0.35; // points
    var strokeColor = aiGenColor(0xffffff);
    var fillColor = aiGenColor(0xcdd1ce);
    var bottomColor = aiGenColor(0x173e35);
    var c = null;
    var labels = ["BASE RENT", "NET OPERATING INCOME", "TOTAL RETURN YR-1"];
    var white = aiGenColor(0xffffff);

    var rent = "";  // if there are no rentSteps, just show blank
    if (typeof property.rentSteps != "undefined" && property.rentSteps.length > 0) {
        rent = fmtCurrency(property.rentSteps[0].Rent);
    }

    for (i = 0; i < 3; i++) {
        c = (i == 2) ? bottomColor : fillColor;
        genRect(layer,skx,sky,w,height,strokeWidth,strokeColor,c);
        genRect(layer,skx+w,sky,w,height,strokeWidth,strokeColor,c);
        genRect(layer,skx+w+w,sky,w,height,strokeWidth,strokeColor,c);
        c = (i == 2) ? white : bottomColor;
        tableText(layer,skx+offsetx, sky-offsety, labels[i], 13, "BebasNeue-Regular", 0, null,c);
        if (rent != "") {
            tableText(layer,skx+offsetx+col3, sky-offsety, rent, 11, "ArialNarrow", 0, null,c);
        }

        sky -= height;
    }
    strokeWidth = 0.5; // points
    genRect(layer,upperLX,upperLY,width,3*height,strokeWidth,bottomColor,null); // encompassing rect

    if (sky < MINY) {
        // alert("*** JB WARNING ***  not enough room for image on Financial Overview");
        return;
    }

    //-------------------------------------------------
    // Now add a smaller version of the cover image...
    //-------------------------------------------------
    var ab = jb.doc.artboards.getByName("Financial Overview");
    if (null == ab) {
        alert('could not find Financial Overview artboard');
        return;
    }
    var bounds = ab.artboardRect;
    var b = {
        left: bounds[0],
        top: bounds[1],
        right: bounds[2],
        bottom: bounds[3],
        width: bounds[2] - bounds[0],
        height: bounds[1] - bounds[3],
    };

    y = sky;  // top - moves top y pos to one cell height below skeleton table
    var lowerRY = b.bottom + 70;

    var pb = {
        left: x,
        top: y,
        width: width,
        height: y - b.bottom,
    };
    if (property.Img1 != "") {
        placeImage(layer, 1, "coverShot", pb);
    }
}
