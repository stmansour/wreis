
// genRect creates a rectangle and adds it to pathItems
//
//    x,y = top left corner of rectangle
//    w,h = width and height of the rectangle
//-----------------------------------------------------------------------------
function genRect(x,y,w,h) {
    var rgb = colorComponents(0x173e35);
    var c = new RGBColor();
    c.red = rgb.r;
    c.green = rgb.g;
    c.blue = rgb.b;

    var r = jb.doc.pathItems.add();
    r.filled = false;
    r.stroked = true;
    r.strokeWidth = 1.0; // in points
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

function tableText(x,y,s) {
    if (jb.chattr == null) {
        alert('jb.chattr is being used before it has been set!');
        return;
    }
    layer = jb.doc.layers.getByName("Financial Overview");
    var t = layer.textFrames.add();


    t.contents = s;
    t.textRange.characterAttributes = jb.chattr;

    var font = app.textFonts.getByName("Arial Narrow");
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
    hbn = "AnnOpDt-Hdr";

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
    var offsetx = 10;  // from the left
    var offsety = 3;   // from the top

    //---------------------------------
    // Do the Rent Steps first...
    //---------------------------------
    for (var i = 0; i < property.rentSteps.length; i++) {
        genRect(x,y,width,height);
        tableText(x+offsetx, y-offsety, property.rentSteps[i].Opt);
        y -=height;
    }
    //---------------------------------
    // Now the Renewal Options
    //---------------------------------
    for (i = 0; i < property.renewOptions.length; i++) {
        genRect(x,y,width,height);
        tableText(x+offsetx, y-offsety, property.renewOptions[i].Opt);
        y -=height;
    }
}
