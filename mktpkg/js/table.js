
// genRect creates a rectangle and adds it to pathItems
//
//    x,y = top left corner of rectangle
//    w,h = width and height of the rectangle
//    f   = boolean for filled or not
//    l   = layer
//-----------------------------------------------------------------------------
function genRect(x,y,w,h,f,l) {
    var c = aiFillColor(0x173e35);
    var fc = aiFillColor(0xcdd1ce);
    var r = l.pathItems.add();
    r.filled = f;
    if (f) {
        r.fillColor = fc;
    }
    r.stroked = true;
    r.strokeWidth = 0.75; // in points
    r.strokeColor = c;

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
    var col2 = 155;     // 150 from left edge
    var col3 = 310;     // from left edge
    var monthly = 0.0;  // monthly amount
    var fill = false;   // every other rect should be filled

    //---------------------------------
    // Do the Rent Steps first...
    //---------------------------------
    for (var i = 0; i < property.rentSteps.length; i++) {
        genRect(x,y,width,height,fill,layer);
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
        genRect(x,y,width,height,fill,layer);
        tableText(x+offsetx, y-offsety, property.renewOptions[i].Opt,layer);
        tableText(x+offsetx+col2, y-offsety, fmtCurrency(property.renewOptions[i].Rent),layer);
        monthly = property.renewOptions[i].Rent / 12.0;
        tableText(x+offsetx+col3, y-offsety, fmtCurrency(monthly),layer);
        y -=height;
        fill = !fill;
    }
}
