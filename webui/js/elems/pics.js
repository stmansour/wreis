/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document, FormData,
    fetch,w2confirm, doDeletePhoto,
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

    //--------------------------------------
    //  Now load all the proper images...
    //--------------------------------------
    var r = w2ui.propertyForm.record;
    for (i = 1; i <= 6; i++) {
        var ImgID = 'Img'+i;
        if (r[ImgID].length > 0) {
            var id = "phototable" + i;
            var image = document.getElementById(id);
            if (image != null) {
                image.src = r[ImgID];
                document.getElementById('spnFilePath' + i).innerHTML = '';
            }
        }
    }
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
        filePath1.innerHTML = "uploading " + fileName;
        w2ui.propertyForm.record.Img1 = fileName;
        SavePhotoToServer(fileName,1,fileupload1.files[0]);
    };

    var fileupload2 = document.getElementById("FileUpload2");
    var filePath2 = document.getElementById("spnFilePath2");
    var image2 = document.getElementById("phototable2");
    image2.onclick = function() {
        fileupload2.click();
    };
    fileupload2.onchange = function() {
        var fileName2 = fileupload2.value.split('\\')[fileupload2.value.split('\\').length - 1];
        filePath2.innerHTML = "uploading " + fileName2;
        w2ui.propertyForm.record.Img2 = fileName2;
        SavePhotoToServer(fileName2,2,fileupload2.files[0]);
    };

    var fileupload3 = document.getElementById("FileUpload3");
    var filePath3 = document.getElementById("spnFilePath3");
    var image3 = document.getElementById("phototable3");
    image3.onclick = function() {
        fileupload3.click();
    };
    fileupload3.onchange = function() {
        var fileName3 = fileupload3.value.split('\\')[fileupload3.value.split('\\').length - 1];
        filePath3.innerHTML = "uploading " + fileName3;
        w2ui.propertyForm.record.Img3 = fileName3;
        SavePhotoToServer(fileName3,3,fileupload3.files[0]);
    };

    var fileupload4 = document.getElementById("FileUpload4");
    var filePath4 = document.getElementById("spnFilePath4");
    var image4 = document.getElementById("phototable4");
    image4.onclick = function() {
        fileupload4.click();
    };
    fileupload4.onchange = function() {
        var fileName4 = fileupload4.value.split('\\')[fileupload4.value.split('\\').length - 1];
        filePath4.innerHTML = "uploading " + fileName4;
        w2ui.propertyForm.record.Img4 = fileName4;
        SavePhotoToServer(fileName4,4,fileupload4.files[0]);
    };

    var fileupload5 = document.getElementById("FileUpload5");
    var filePath5 = document.getElementById("spnFilePath5");
    var image5 = document.getElementById("phototable5");
    image5.onclick = function() {
        fileupload5.click();
    };
    fileupload5.onchange = function() {
        var fileName5 = fileupload5.value.split('\\')[fileupload5.value.split('\\').length - 1];
        filePath5.innerHTML = "uploading " + fileName5;
        w2ui.propertyForm.record.Img5 = fileName5;
        SavePhotoToServer(fileName5,5,fileupload5.files[0]);
    };

    var fileupload6 = document.getElementById("FileUpload6");
    var filePath6 = document.getElementById("spnFilePath6");
    var image6 = document.getElementById("phototable6");
    image6.onclick = function() {
        fileupload6.click();
    };
    fileupload6.onchange = function() {
        var fileName6 = fileupload6.value.split('\\')[fileupload6.value.split('\\').length - 1];
        filePath6.innerHTML = "uploading " + fileName6;
        w2ui.propertyForm.record.Img6 = fileName6;
        SavePhotoToServer(fileName6,6,fileupload6.files[0]);
    };
}

// SavePhotoToServer.  Save the specified image to the property currently in
// w2ui.propertyForm and save it under the supplied index.
//
// INPUTS
//    f = fileName of image
//    x = index
//    file = the actual file from the <input ...> object
//-----------------------------------------------------------------------------
 function SavePhotoToServer(f,x,file)
{
    if (typeof file === "undefined") {
        document.getElementById('spnFilePath' + x).innerHTML = '';
        return;
    }
    let data = { cmd:'save', PRID: propData.PRID, idx: x, fileName: f };
    let formData = new FormData();
    let url = '/v1/propertyphoto/' + propData.PRID + '/' + x;
    var rr;

    formData.append("request", JSON.stringify(data));
    formData.append("photo", file);

    doSaveImage(url,formData)
    .then(resp => {
        //---------------------------------------------------------------------------
        // now we need to set the URL of the image to what has just been returned...
        //---------------------------------------------------------------------------
        var id = "phototable" + x;
        var image = document.getElementById(id);
        if (image == null) {
            console.log('ERROR: could not find image: ' + id);
            return;
        }
        image.src = resp.url;
        document.getElementById('spnFilePath' + x).innerHTML = '';
    });

}

async function doSaveImage(url,formData) {
    var rr;
    try {
       rr = await fetch(url, {method: "POST", body: formData});
    } catch(e) {
       console.log('Error: ', e);
       return null;
    }

    let resp = await rr.json();
    return resp;
}

function deleteImg(x) {
    w2confirm('Are you sure you want to delete image '+x)
    .yes(function () {
        doDeletePhoto(x)
    })
    .no(function () {
        console.log('NO');
     }
 );
}

function doDeletePhoto(x) {
    let data = { cmd:'delete', PRID: propData.PRID, idx: x };
    var dat = JSON.stringify(data);
    var url = '/v1/propertyphotodelete/' + propData.PRID + '/' + x;

    return $.post(url, dat, null, "json")
        .done(function(data) {
            // if (data.status === "success") {
            // }
            if (data.status === "error") {
                w2ui.propertyGrid.error('ERROR: '+ data.message);
                return;
            }
            var id = "phototable" + x;
            var image = document.getElementById(id);
            if (image == null) {
                console.log('ERROR: could not find image: ' + id);
                return;
            }
            image.src = '/static/html/images/building-100.png';
        })
        .fail(function(data){
                w2ui.propertyGrid.error("Save RentableLeaseStatus failed. " + data);
        });
}
