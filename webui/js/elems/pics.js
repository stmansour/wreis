/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document,
*/

"use strict";

var AppPics = {
    WaitDepth: 0,
};

function buildPropertyPhotosUIElements() {
    $().w2layout({
        name: 'propertyPhotosLayout',
        padding: 0,
        panels: [
            { type: 'left',    size: 0,     hidden: true },
            { type: 'top',     size: 0,     hidden: true,  content: 'top',  resizable: true, style: app.pstyle },
            { type: 'main',    size: '100%', hidden: false, content: 'main', resizable: true, style: app.pstyle },
            { type: 'preview', size: 0,     hidden: true,  content: 'PREVIEW'  },
            { type: 'bottom',  size: 0,     hidden: true,  content: 'bottom', resizable: false, style: app.pstyle },
            { type: 'right',   size: 0,     hidden: true,  content: 'right' }
        ],
        onRender: function(event) {
            event.onComplete = function(event) {
                AwaitImagePanelRenderComplete();
            };
        },
    });
}

function AwaitImagePanelRenderComplete() {
    var done = true;
    for (var i = 1; i <= 6; i++) {
        var s = 'FileUpload' + i;
        var x = document.getElementById(s);
        if (x == null) {
            done = false;
            break;
        }
    }
    if (! done ) {
        AppPics.WaitDepth++;
        if ( AppPics.WaitDepth > 10) {
            AppPics.WaitDepth = 0;
            console.log('Wait depth > 10 in image panel render.  Something is wrong!');
            return;
        }
        setTimeout(() => { AwaitImagePanelRenderComplete(); /* yes, recurse*/ }, 500);
        return; // don't call SetUp here...
    }
    //-------------------------------------------------------------
    // if we hit this point, everything we need is now available
    //-------------------------------------------------------------
    AppPics.WaitDepth = 0; // now is the time to reset this
    SetUpImageCatchers();
}

function SetUpImageCatchers() {
    var fileupload1 = document.getElementById("FileUpload1");
    var filePath1 = document.getElementById("spnFilePath1");
    var image = document.getElementById("phototable1");
    image.onclick = function() {
        fileupload1.click();
    };
    fileupload1.onchange = function() {
        var fileName = fileupload1.value.split('\\')[fileupload1.value.split('\\').length - 1];
        filePath1.innerHTML = "<b>Selected File: </b>" + fileName;
        w2ui.propertyForm.record.Img1 = "fileName";
    };

    var fileupload2 = document.getElementById("FileUpload2");
    var filePath2 = document.getElementById("spnFilePath2");
    image = document.getElementById("phototable2");
    image.onclick = function() {
        fileupload2.click();
    };
    fileupload2.onchange = function() {
        var fileName = fileupload2.value.split('\\')[fileupload2.value.split('\\').length - 1];
        filePath2.innerHTML = "<b>Selected File: </b>" + fileName;
        w2ui.propertyForm.record.Img2 = "fileName";
    };

    var fileupload3 = document.getElementById("FileUpload3");
    var filePath3 = document.getElementById("spnFilePath3");
    image = document.getElementById("phototable3");
    image.onclick = function() {
        fileupload3.click();
    };
    fileupload3.onchange = function() {
        var fileName = fileupload3.value.split('\\')[fileupload3.value.split('\\').length - 1];
        filePath3.innerHTML = "<b>Selected File: </b>" + fileName;
        w2ui.propertyForm.record.Img3 = "fileName";
    };

    var fileupload4 = document.getElementById("FileUpload4");
    var filePath4 = document.getElementById("spnFilePath4");
    image = document.getElementById("phototable4");
    image.onclick = function() {
        fileupload4.click();
    };
    fileupload4.onchange = function() {
        var fileName = fileupload4.value.split('\\')[fileupload4.value.split('\\').length - 1];
        filePath4.innerHTML = "<b>Selected File: </b>" + fileName;
        w2ui.propertyForm.record.Img4 = "fileName";
    };

    var fileupload5 = document.getElementById("FileUpload5");
    var filePath5 = document.getElementById("spnFilePath5");
    image = document.getElementById("phototable5");
    image.onclick = function() {
        fileupload5.click();
    };
    fileupload5.onchange = function() {
        var fileName = fileupload5.value.split('\\')[fileupload5.value.split('\\').length - 1];
        filePath5.innerHTML = "<b>Selected File: </b>" + fileName;
        w2ui.propertyForm.record.Img5 = "fileName";
    };

    var fileupload6 = document.getElementById("FileUpload6");
    var filePath6 = document.getElementById("spnFilePath6");
    image = document.getElementById("phototable6");
    image.onclick = function() {
        fileupload6.click();
    };
    fileupload6.onchange = function() {
        var fileName = fileupload6.value.split('\\')[fileupload6.value.split('\\').length - 1];
        filePath6.innerHTML = "<b>Selected File: </b>" + fileName;
        w2ui.propertyForm.record.Img6 = "fileName";
    };
}
