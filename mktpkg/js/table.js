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

function tableText(x,y,s,layer) {
    if (jb.chattr == null) {
        alert('jb.chattr is being used before it has been set!');
        return;
    }
    var t = layer.textFrames.add();


    t.contents = s;
    t.textRange.characterAttributes = jb.chattr;
    t.textRange.characterAttributes.size = 11;

    var font = app.textFonts.getByName("ArialNarrow");
    if (font != null) {
        t.textRange.characterAttributes.textFont = font;
    }

    // alert("x,y = " + x + ', ' + y);
    t.textRange.justification = Justification.LEFT;
    t.position = [x,y];
}


// genTable
//
//-----------------------------------------------------------------------------
function genTable() {
    var hbn = "AnnOpDt-Hdr";
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
    var nrows = property.renewOptions.length + property.rentSteps.length;
    // alert("nrows = " + nrows);
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

    //---------------------------------
    // Do the Rent Steps first...
    //---------------------------------
    for (var i = 0; i < property.rentSteps.length; i++) {
        genAODRect(layer,x,y,width,height,fill);
        tableText(x+offsetx, y-offsety, property.rentSteps[i].Opt,layer);
        tableText(x+offsetx+col2, y-offsety, fmtCurrency(property.rentSteps[i].Rent),layer);
        monthly = property.rentSteps[i].Rent / 12.0;
        tableText(x+offsetx+col3, y-offsety, fmtCurrency(monthly),layer);
        y -=height;
        fill = !fill;
    }
    //---------------------------------
    // Now the Renewal Options
    //---------------------------------
    for (i = 0; i < property.renewOptions.length; i++) {
        genAODRect(layer,x,y,width,height,fill);
        tableText(x+offsetx, y-offsety, property.renewOptions[i].Opt,layer);
        tableText(x+offsetx+col2, y-offsety, fmtCurrency(property.renewOptions[i].Rent),layer);
        monthly = property.renewOptions[i].Rent / 12.0;
        tableText(x+offsetx+col3, y-offsety, fmtCurrency(monthly),layer);
        y -=height;
        fill = !fill;
    }

    //---------------------------------
    // Now the skeleton table
    //---------------------------------
    var skx = x;
    var sky = y - height;
    var upperLX = skx;
    var upperLY = sky;
    var w = width/3;

    var strokeWidth = 0.35; // points
    var strokeColor = aiGenColor(0xffffff);
    var fillColor = aiGenColor(0xcdd1ce);
    var bottomColor = aiGenColor(0x173e35);
    var c = null;

    for (i = 0; i < 3; i++) {
        c = (i == 2) ? bottomColor : fillColor;
        genRect(layer,skx,sky,w,height,strokeWidth,strokeColor,c);
        genRect(layer,skx+w,sky,w,height,strokeWidth,strokeColor,c);
        genRect(layer,skx+w+w,sky,w,height,strokeWidth,strokeColor,c);
        sky -= height;
    }
    strokeWidth = 0.5; // points
    genRect(layer,upperLX,upperLY,width,3*height,strokeWidth,bottomColor,null); // encompassing rect
}
