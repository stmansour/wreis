function number_format(number, decimals, dec_point, thousands_sep) {
    var n = !isFinite(+number) ? 0 : +number,
        prec = !isFinite(+decimals) ? 0 : Math.abs(decimals),
        sep = (typeof thousands_sep === 'undefined') ? ',' : thousands_sep,
        dec = (typeof dec_point === 'undefined') ? '.' : dec_point,
        toFixedFix = function (n, prec) {
            // Fix for IE parseFloat(0.55).toFixed(0) = 0;
            var k = Math.pow(10, prec);
            return Math.round(n * k) / k;
        },
        s = (prec ? toFixedFix(n, prec) : Math.round(n)).toString().split('.');
    if (s[0].length > 3) {
        s[0] = s[0].replace(/\B(?=(?:\d{3})+(?!\d))/g, sep);
    }
    if ((s[1] || '').length < prec) {
        s[1] = s[1] || '';
        s[1] += new Array(prec - s[1].length + 1).join('0');
    }
    return s.join(dec);
}

function fmtWithCommas(x) {
    return number_format(x,0,'.',',');
}

function fmtAsFloatWithCommas(x) {
    return number_format(x,2,'.',',');
}

function fmtCurrency(x) {
    return '$' + number_format(x,2,'.',',');
}

function fmtAsPercent(x) {
    return number_format(100*x,2,'.',',') + '%';
}

function fmtIndexedName(i,aiName,arr,errLabel) {
    if (arr.length < 1) {
        alert("fmtIndexedName: arr.length < 1 for ainame = " + aiName);
        return;
    }
    if (i < 0) {
        i = 0;
    }
    t = jb.doc.textFrames.getByName(aiName);
    if (i + 1 > arr.length) {
        t.contents = "(unknown "+errLabel+")";
    } else {
        if (arr[i] == null) {
            alert("fmtIndexedName: arr.length = " + arr.length + ", i = " + i + "ainame = " + aiName);
            return;
        }
        t.contents = arr[i];
    }
}

function fmtDate( s,aiName) {
    t = jb.doc.textFrames.getByName(aiName);

    // The date parser in AI's javascript seems to really be out of date.
    // We do this only because it doesn't know how to parse UTC string dates
    // formatted the way the server formats them...
    //
    // The date strings look like this:    2020-03-22 07:00:00 UTC
    //----------------------------------------------------------------------
    var a1 = s.split(" ");
    var a2 = a1[0].split("-");
    t.contents = a2[1] + "/" + a2[2] + "/" + a2[0];
}

// It is expected that the string is in this format:
//    2020-03-22 07:00:00 UTC
function AIDate(s) {
    var a1 = s.split(" ");
    var a2 = a1[0].split("-");
    var y = parseInt(a2[0]);
    var m = parseInt(a2[1]);  // month index, zero-based
    var d = parseInt(a2[2]);
    var a3 = a1[1].split(":");
    var H = parseInt(a3[0]);
    var M = parseInt(a3[1]);
    var S = parseInt(a3[2]);
    var dt = new Date(y,m-1,d,H,M,S);
    return dt;
}

var DateDiff = {
    inDays: function(d1, d2) {
        var t2 = d2.getTime();
        var t1 = d1.getTime();

        return parseInt((t2-t1)/(24*3600*1000));
    },

    inWeeks: function(d1, d2) {
        var t2 = d2.getTime();
        var t1 = d1.getTime();

        return parseInt((t2-t1)/(24*3600*1000*7));
    },

    inMonths: function(d1, d2) {
        var d1Y = d1.getFullYear();
        var d2Y = d2.getFullYear();
        var d1M = d1.getMonth();
        var d2M = d2.getMonth();

        return (d2M+12*d2Y)-(d1M+12*d1Y);
    },

    inYears: function(d1, d2) {
        return d2.getFullYear()-d1.getFullYear();
    }
};

function dateDiffInYears(dateold, datenew) {
    var ynew = datenew.getFullYear();
    var mnew = datenew.getMonth();
    var dnew = datenew.getDate();
    var yold = dateold.getFullYear();
    var mold = dateold.getMonth();
    var dold = dateold.getDate();
    var diff = ynew - yold;
    if (mold > mnew) diff--;
    else {
        if (mold == mnew) {
            if (dold > dnew) diff--;
        }
    }
    return diff;
}

function fmtDateDiffInYears(d1,d2) {
    var t;
    if (typeof d1 === "string") {
        t = AIDate(d1);
        d1 = t;
    }
    if (typeof d2 === "string") {
        t = AIDate(d2);
        d2 = t;
    }
    var diff = DateDiff.inYears(d1,d2);
    if (diff > 0) {
        return '' + diff + " Years";
    }
    return 'Lease has expired';
}

//
// Breaks up the supplied hex number into 3 parts: rgb
//
// x = colorComponents(0x2020ff);
//
// returns: RGBColor object
// ---------------------------------------------------------------------------
function colorComponents(h) {
    var x = {
        r: (h >> 16) & 0xff,
        g: (h >> 8) & 0xff,
        b: h & 0xff,
    };
    return x;
}

function aiGenColor(h) {
    var rgb = colorComponents(h);
    var c = new RGBColor();
    c.red = rgb.r;
    c.green = rgb.g;
    c.blue = rgb.b;
    return c;
}

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
