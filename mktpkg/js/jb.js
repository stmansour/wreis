//
//  JB - the portfolio writer :-)
//

//--------------------------------------------------------------
// Make sure to close any currently open document...
//--------------------------------------------------------------
if (!app.homeScreenVisible) {
    app.activeDocument.close(SaveOptions.PROMPTTOSAVECHANGES);
}

//---------------------------------------------------------------------------
// By convention, we will keep templates in ~/Documents/wreis.
// We get the myDocuments folder from the Adobe environment which maps it
// to the file system correctly even on Windows.  We will look for a folder
// named ~/Documents/wreis and open the file template00.ai
//---------------------------------------------------------------------------
var template = 'template00.ai';
var fname = Folder.myDocuments + '/wreis/' + template;
var f = new File(fname);
app.open(f);

//---------------------------------------------------------------------------
// immediately save this as a new document: portfolio.ai
//---------------------------------------------------------------------------
fname = Folder.myDocuments + '/wreis/portfolio.ai';
var portfolio = new File(fname);
app.activeDocument.saveAs(portfolio);
