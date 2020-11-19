
/*global
    w2ui, app, $, console, document,
*/

"use strict";

function getDashboard() {

    var params = {
        cmd: "get",
    };
    var dat = JSON.stringify(params);
    var url = '/v1/dashboard/';

    $.post(url, dat, null, "json")
    .done(function(data) {
        if (typeof data == 'string') {  // it's weird, a successful data add gets parsed as an object, an error message does not
            var msg = JSON.parse(data);
            console.log('Response to dashboard: ' + msg.status);
            return;
        }
        if (data.status == 'success') {
            document.getElementById('PropertyCount').innerHTML = '' + data.record.PropertyCount;
            document.getElementById('CompletedProperties').innerHTML = '' + data.record.CompletedProperties;
            document.getElementById('ActiveProperties').innerHTML = '' + data.record.ActiveProperties;
            document.getElementById('YourQueue').innerHTML = '' + data.record.YourQueue;
        } else {
            console.log('data.status = ' + data.status);
        }
    })
    .fail(function(data) {
        console.log('data = ' + data);
    });
}
