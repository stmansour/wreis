//
// TO RUN THIS SCRIPT IN ADOBE ILLUSTRATOR
//
//  File -> Scripts -> Other Script...
//  then select this file.
//
//---------------------------------------------------------------
var mydoc = app.documents.add(DocumentColorSpace.RGB,1920,1080);
mydoc.layers[0].name = "Background";

// resizedoc();
centerText("Hello, world!");


function centerText(s) {
    var doc = app.activeDocument;
    var myTextFrame = doc.textFrames.add();
    myTextFrame.position = [doc.width/2,doc.height/2];
    myTextFrame.contents = "Hello World!";
    var tRange = myTextFrame.textRange;
    tRange.size = 32;
    tRange.justification = Justification.CENTER;
}

function showRectVals( r ) {
    alert("r.top = " + r.top + ", r.left = " + r.left + "r.width = " + r.width + ", r.height = " + r.height);

}

function inputtest(){
    // if there are no documents open, alert the user and exit
    if( app.documents.length == 0 ) {
        alert( "Please open a document!" );
        return;
    }

    // get a ref to the current document
    var doc = app.activeDocument;
    var artboard = doc.artboards[0];

    // flag for user input
    var validInput = false;
    var promptMessage = "Please enter the size of the rectangle in inches";
    var dims, rectWidth, rectHeight;

    // this section of code will be repeated until either of the two flags are set equal to true
    while( validInput == false ) {
        // prompt for the dimensions
        dims = prompt( promptMessage, "8.5 x 11" );
        // if the dims var is equal to null, assume the user clicked the cancel button and exit early
        if( dims === null ) {
            return;
        } else {
            // otherwise, assume the user clicked OK, try to parse the input
            // turn the string into an array, splitting on the letter 'x'
            dims = dims.split('x');

            // if there are only two elements in the array
            if( dims.length == 2 ) {
                // try to convert them to number values
                rectWidth = parseFloat( dims[0] );
                rectHeight = parseFloat( dims[1] );

                // if neither value is equal to NaN (not a number), then set the flag to let us out of the loop
                if( !isNaN( rectWidth ) && !isNaN( rectHeight ) )
                    validInput = true;
            }
        }

        // if we have to prompt again, update the message to have some more info
        promptMessage = "Please enter the size of the rectangle in inches\nValues should be in the format '<width> x <height>'";
    }

    // if we made it out of the loop, assume that the values are valid numbers!
    // convert the numbers to points and make sure they are both positive
    rectWidth = Math.abs( rectWidth * 72 );
    rectHeight = Math.abs( rectHeight * 72 );

    // get current artboard dimensions
    var bounds = artboard.artboardRect;

    // add new rectangle                                    top                   left                  width           height
    var rect = doc.pathItems.rectangle( bounds[0], bounds[1], rectWidth, rectHeight );
    artboard.artboardRect = rect.geometricBounds;

    // scale rectangle, numbers
    rect.resize(
            25,                // scale horizontal in %
            100,               // scale vertical in %
            true,              // change position
            false,             // change Fill Pattern
            false,             // change Fill Gradient
            false,             // change Stroke Pattern
            1,                 // change Line Widths
            Transformation.RIGHT );  // scale about which point
}
