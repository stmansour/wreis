/*global
    w2ui, app, $, console, dateFmtStr, propData, Promise, document, FormData,
    fetch,w2confirm, doDeletePhoto,window,
*/

"use strict";

var AppPics = {
    WaitDepth: 0,
    NumImgs: 8,
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
    for (var i = 1; i <= AppPics.NumImgs; i++) {
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
    for (i = 1; i <= AppPics.NumImgs; i++) {
        var ImgID = 'Img'+i;
        if (r[ImgID].length > 0) {
            var id = "phototable" + i;
            var image = document.getElementById(id);
            if (image != null) {
                image.src = r[ImgID];
                document.getElementById('spnFilePath' + i).innerHTML = '';
            }
            image.width = 190;
        }
    }
}

// // propChangeImage allows the user to change the selected image.
// //
// // INPUTS
// // x = the index number 1 - 8 for which image they're switching
// //
// //-----------------------------------------------------------------------
// function propChangeImage(x) {
//     var f = document.getElementById("FileUpload"+x);
//     if (f == null) {
//         return;
//     }
//     var fp = document.getElementById("spnFilePath" + x);
//     if (fp == null ) {
//         return;
//     }
//     var image = "Img" + x;
//     var idx = x;
//     f.onchange = function() {
//         var fileName = f.value.split('\\')[f.value.split('\\').length - 1];
//         fp.innerHTML = "uploading " + fileName;
//         w2ui.propertyForm.record[image] = fileName;
//         SavePhotoToServer(fileName,idx,f.files[0]);
//     }
// }

function SetUpImageCatchers() {
    var fileupload1 = document.getElementById("FileUpload1");
    var filePath1 = document.getElementById("spnFilePath1");
    var pt1 = document.getElementById("phototable1");
    var image1 = document.getElementById("editPic1");
    image1.onclick = function() { fileupload1.click(); };
    pt1.onclick = function() { fileupload1.click(); };
    fileupload1.onchange = function() {
        var fileName = fileupload1.value.split('\\')[fileupload1.value.split('\\').length - 1];
        filePath1.innerHTML = "uploading " + fileName;
        w2ui.propertyForm.record.Img1 = fileName;
        SavePhotoToServer(fileName,1,fileupload1.files[0]);
    };

    var fileupload2 = document.getElementById("FileUpload2");
    var filePath2 = document.getElementById("spnFilePath2");
    var pt2 = document.getElementById("phototable2");
    var image2 = document.getElementById("editPic2");
    image2.onclick = function() { fileupload2.click(); };
    pt2.onclick = function() { fileupload2.click(); };
    fileupload2.onchange = function() {
        var fileName2 = fileupload2.value.split('\\')[fileupload2.value.split('\\').length - 1];
        filePath2.innerHTML = "uploading " + fileName2;
        w2ui.propertyForm.record.Img2 = fileName2;
        SavePhotoToServer(fileName2,2,fileupload2.files[0]);
    };

    var fileupload3 = document.getElementById("FileUpload3");
    var filePath3 = document.getElementById("spnFilePath3");
    var pt3 = document.getElementById("phototable3");
    var image3 = document.getElementById("editPic3");
    image3.onclick = function() { fileupload3.click(); };
    pt3.onclick = function() { fileupload3.click(); };
    fileupload3.onchange = function() {
        var fileName3 = fileupload3.value.split('\\')[fileupload3.value.split('\\').length - 1];
        filePath3.innerHTML = "uploading " + fileName3;
        w2ui.propertyForm.record.Img3 = fileName3;
        SavePhotoToServer(fileName3,3,fileupload3.files[0]);
    };

    var fileupload4 = document.getElementById("FileUpload4");
    var filePath4 = document.getElementById("spnFilePath4");
    var pt4 = document.getElementById("phototable4");
    var image4 = document.getElementById("editPic4");
    image4.onclick = function() { fileupload4.click(); };
    pt4.onclick = function() { fileupload4.click(); };
    fileupload4.onchange = function() {
        var fileName4 = fileupload4.value.split('\\')[fileupload4.value.split('\\').length - 1];
        filePath4.innerHTML = "uploading " + fileName4;
        w2ui.propertyForm.record.Img4 = fileName4;
        SavePhotoToServer(fileName4,4,fileupload4.files[0]);
    };

    var fileupload5 = document.getElementById("FileUpload5");
    var filePath5 = document.getElementById("spnFilePath5");
    var pt5 = document.getElementById("phototable5");
    var image5 = document.getElementById("editPic5");
    image5.onclick = function() { fileupload5.click(); };
    pt5.onclick = function() { fileupload5.click(); };
    fileupload5.onchange = function() {
        var fileName5 = fileupload5.value.split('\\')[fileupload5.value.split('\\').length - 1];
        filePath5.innerHTML = "uploading " + fileName5;
        w2ui.propertyForm.record.Img5 = fileName5;
        SavePhotoToServer(fileName5,5,fileupload5.files[0]);
    };

    var fileupload6 = document.getElementById("FileUpload6");
    var filePath6 = document.getElementById("spnFilePath6");
    var pt6 = document.getElementById("phototable6");
    var image6 = document.getElementById("editPic6");
    image6.onclick = function() { fileupload6.click(); };
    pt6.onclick = function() { fileupload6.click(); };
    fileupload6.onchange = function() {
        var fileName6 = fileupload6.value.split('\\')[fileupload6.value.split('\\').length - 1];
        filePath6.innerHTML = "uploading " + fileName6;
        w2ui.propertyForm.record.Img6 = fileName6;
        SavePhotoToServer(fileName6,6,fileupload6.files[0]);
    };

    var fileupload7 = document.getElementById("FileUpload7");
    var filePath7 = document.getElementById("spnFilePath7");
    var pt7 = document.getElementById("phototable7");
    var image7 = document.getElementById("editPic7");
    image7.onclick = function() { fileupload7.click(); };
    pt7.onclick = function() { fileupload7.click(); };
    fileupload7.onchange = function() {
        var fileName7 = fileupload7.value.split('\\')[fileupload7.value.split('\\').length - 1];
        filePath7.innerHTML = "uploading " + fileName7;
        w2ui.propertyForm.record.Img7 = fileName7;
        SavePhotoToServer(fileName7,7,fileupload7.files[0]);
    };

    var fileupload8 = document.getElementById("FileUpload8");
    var filePath8 = document.getElementById("spnFilePath8");
    var pt8 = document.getElementById("phototable8");
    var image8 = document.getElementById("editPic8");
    image8.onclick = function() { fileupload8.click(); };
    pt8.onclick = function() { fileupload8.click(); };
    fileupload8.onchange = function() {
        var fileName8 = fileupload8.value.split('\\')[fileupload8.value.split('\\').length - 1];
        filePath8.innerHTML = "uploading " + fileName8;
        w2ui.propertyForm.record.Img8 = fileName8;
        SavePhotoToServer(fileName8,8,fileupload8.files[0]);
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
    var id = "phototable" + x;
    var image = document.getElementById(id);
    if (image == null) {
        console.log('ERROR: could not find image: ' + id);
        return;
    }
    image.src = "/static/html/images/spinner.gif";
    image.width = 190;
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
        image.width = 190;
        document.getElementById('spnFilePath' + x).innerHTML = '';
        w2ui.propertyForm.record["Img"+x] = resp.url;
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

// deleteImg is called when the user presses one of the trashcan buttons in the
// photos form.  If no photo is assigned to the slot they chose then it just
// returns.  Otherwise it will ask if the user is sure they want to delete the
// photo. If they do, it will call doDeletePhoto()
//
// INPUTS
//     x = index number of the image.
//------------------------------------------------------------------------------
function deleteImg(x) {
    var id='Img' + x;
    var url=w2ui.propertyForm.record[id];
    if (url.length < 1) {
        return;
    }

    w2confirm('Are you sure you want to delete image '+x)
    .yes(function () {
            doDeletePhoto(x);
        }
    )
    .no(function () {
            console.log('NO');
         }
    );
}

// propImageView loads the full image into a new page.
//
// INPUTS
//     x = index number of the image.
//------------------------------------------------------------------------------
function propImageView(x) {
    var id='Img' + x;
    var url=w2ui.propertyForm.record[id];
    if (url.length < 1) {
        return;
    }
    window.open(url);
}

// doDeletePhoto is called when the user acknowledges that they want to delete
// one of the photos in the photos form. If the slot (1 - 8) has an associated
// image, it will be deleted.
//------------------------------------------------------------------------------
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
            image.width = 100;
            w2ui.propertyForm.record["Img"+x] = "";
        })
        .fail(function(data){
                w2ui.propertyGrid.error("Save RentableLeaseStatus failed. " + data);
        });
}
