package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	db "wreis/db/lib"
	util "wreis/util/lib"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

/*
**
**  Required constants for S3 configuration.
**
**/
const (
	S3Region        = "us-east-1" // This parameter define the region of bucket
	ImageUploadPath = ""          // This parameter define in which folder have to upload image
	MaxPhotoIndex   = 8
)

var reImagePart = regexp.MustCompile(` filename="([^"]+)"`)
var rePhotoURL = regexp.MustCompile(`img-(\d+)-(\d+)-(.*)`)

type imageFileReq struct {
	Cmd      string
	PRID     int64
	Idx      int64
	Filename string
}

type imgSuccess struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// The contractors set up extres to not include S3Region.  This turned out to be
// needed as WREIS was set to a different region.
//------------------------------------------------------------------------------
func getS3Region() string {
	rg := S3Region // default
	if len(db.Wdb.Config.S3Region) > 0 {
		rg = db.Wdb.Config.S3Region
	}
	return rg
}

// SvcHandlerPropertyPhoto uploads a single file belonging to a property. The
// URL is of the form:
//-----------------------------------------------------------------------------
func SvcHandlerPropertyPhoto(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	if d.MimeMultipartOnly {
		handlePhotoSave(w, r, d)
		return
	}
	switch d.wsSearchReq.Cmd {
	case "delete":
		handlePhotoDelete(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// parsePhotoURL
//
//   https://s3.amazonaws.com/directory-pics/img-1-1-roller-32.png
//------------------------------------------------------------------------------
func parsePhotoURL(url string) (int64, int64, string) {
	var fname string
	var PRID, Idx int64
	s := filepath.Base(url)                  // gets us to: img-1-1-roller-32.png
	ba := rePhotoURL.FindSubmatch([]byte(s)) // [1] = "1", [2] = "1", [3] = "roller-32.png"
	if len(ba) >= 4 {
		PRID, _ = util.IntFromString(string(ba[1]), "bla")
		Idx, _ = util.IntFromString(string(ba[2]), "bla")
		fname = string(ba[3])
	}
	return PRID, Idx, fname
}

// handlePhotoDelete the photo with index idx from PRID
// URL is of the form:
//                      d.ID  d.SUBID
//    /v1/propertyphoto/PRID/idx
//
//-----------------------------------------------------------------------------
func handlePhotoDelete(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "handlePhotoDelete"
	util.Console("Entered %s\n", funcname)

	//---------------------------
	// Quick sanity check...
	//---------------------------
	if d.ID < 1 || d.SUBID < 1 || d.SUBID > MaxPhotoIndex {
		e := fmt.Errorf("%s: Error in property value PRID (%d) or Image Index (%d)", funcname, d.ID, d.SUBID)
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	var PRID, Idx int64
	var Fname, tmpURL string

	//---------------------------
	// get property
	//---------------------------
	pr, err := db.GetProperty(r.Context(), d.ID)
	if err != nil {
		e := fmt.Errorf("%s: Error from db.GetProperty:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	switch d.SUBID {
	case 1:
		PRID, Idx, Fname = parsePhotoURL(pr.Img1)
		tmpURL = pr.Img1
		pr.Img1 = ""
	case 2:
		PRID, Idx, Fname = parsePhotoURL(pr.Img2)
		tmpURL = pr.Img2
		pr.Img2 = ""
	case 3:
		PRID, Idx, Fname = parsePhotoURL(pr.Img3)
		tmpURL = pr.Img3
		pr.Img3 = ""
	case 4:
		PRID, Idx, Fname = parsePhotoURL(pr.Img4)
		tmpURL = pr.Img4
		pr.Img4 = ""
	case 5:
		PRID, Idx, Fname = parsePhotoURL(pr.Img5)
		tmpURL = pr.Img5
		pr.Img5 = ""
	case 6:
		PRID, Idx, Fname = parsePhotoURL(pr.Img6)
		tmpURL = pr.Img6
		pr.Img6 = ""
	case 7:
		PRID, Idx, Fname = parsePhotoURL(pr.Img7)
		tmpURL = pr.Img7
		pr.Img7 = ""
	case 8:
		PRID, Idx, Fname = parsePhotoURL(pr.Img8)
		tmpURL = pr.Img8
		pr.Img8 = ""
	}

	if d.ID != PRID || d.SUBID != Idx {
		e := fmt.Errorf("%s: Image %s is not associated with PRID %d, Idx %d", funcname, tmpURL, PRID, Idx)
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	if err = DeleteS3ImageFile(Fname, PRID, Idx); err != nil {
		e := fmt.Errorf("%s: Error deleting image from S3: %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	if err = db.UpdateProperty(r.Context(), &pr); err != nil {
		e := fmt.Errorf("%s: Error from db.UpdateProperty:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	var resp imgSuccess
	resp.Status = "success"
	resp.URL = ""

	SvcWriteResponse(&resp, w)
}

// handlePhotoSave uploads a single file belonging to a property. The
// URL is of the form:
//
//    /v1/propertyphoto/PRID/idx
//
//-----------------------------------------------------------------------------
func handlePhotoSave(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var req imageFileReq
	var fname string
	var imgFileBytes []byte
	foundReq := false
	foundImg := false
	funcname := "SvcHandlerPropertyPhoto"
	util.Console("Entered %s\n", funcname)
	contentType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !strings.HasPrefix(contentType, "multipart/") {
		e := fmt.Errorf("expecting a multipart message.  %d", http.StatusBadRequest)
		SvcErrorReturn(w, e)
		return
	}

	rbody := bytes.NewReader(d.b)
	reader := multipart.NewReader(rbody, params["boundary"])
	defer r.Body.Close()

	//---------------------------------------------------------
	// Scan the parts, make sure we get the info we need...
	//---------------------------------------------------------
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			e := fmt.Errorf("%s: %s, %d", "unexpected error when retrieving a part of the message", err.Error(), http.StatusInternalServerError)
			SvcErrorReturn(w, e)
			return
		}
		defer part.Close()

		fileBytes, err := ioutil.ReadAll(part)
		if err != nil {
			e := fmt.Errorf("%s:  %d", "failed to read content of the part", http.StatusInternalServerError)
			SvcErrorReturn(w, e)
			return
		}

		av := part.Header["Content-Disposition"]
		val := av[0]
		util.Console("Content-Disposition:  value = %s\n", val)

		//---------------------------------------------------------------------
		// look for Content-Disposition: name="request"  this will give us the
		// structure to unmarshal
		//---------------------------------------------------------------------
		if strings.Contains(val, ` name="request"`) {
			if err := json.Unmarshal(fileBytes, &req); err != nil {
				e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
				SvcFuncErrorReturn(w, e, funcname)
				return
			}
			foundReq = true
		}

		//---------------------------------------------------------------------
		// look for Content-Disposition: filename="<<whatever>>" this will
		// give us the image file we need to save to S3
		//---------------------------------------------------------------------
		if reImagePart.MatchString(val) {
			af := reImagePart.FindSubmatch([]byte(val))
			if len(af) > 1 {
				util.Console("af[0] = %s, af[1] = %s\n", string(af[0]), string(af[1]))
				fname = string(af[1])
			}
			foundImg = true
			imgFileBytes = fileBytes
		}
	}

	//---------------------------------------------------------
	// Reject if we didn't get everything we need...
	//---------------------------------------------------------
	if !foundReq || !foundImg || req.PRID != d.ID || fname != req.Filename || req.Idx < 1 || req.Idx > MaxPhotoIndex {
		var reasons string
		if !foundReq {
			reasons += fmt.Sprintf("No request was found. ")
		}
		if !foundImg {
			reasons += fmt.Sprintf("No image was found. ")
		}
		if req.PRID != d.ID {
			reasons += fmt.Sprintf("PRID in url is %d, but in request it is %d. ", d.ID, req.PRID)
		}
		if fname != req.Filename {
			reasons += fmt.Sprintf("filename in MIME part is %s, but in request it is %s. ", fname, req.Filename)
		}
		if req.Idx < 1 || req.Idx > 6 {
			reasons += fmt.Sprintf("idx = %d is out of range. ", req.Idx)
		}

		e := fmt.Errorf("Command request format is not correct: %s", reasons)
		SvcErrorReturn(w, e)
		return
	}

	util.Console("len(imgFileBytes)=%d\n", len(imgFileBytes))

	_, imageURL, err := UploadImageToS3(fname, imgFileBytes, req.PRID, req.Idx)
	if err != nil {
		e := fmt.Errorf("%s: Error from UploadImageToS3:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	pr, err := db.GetProperty(r.Context(), req.PRID)
	if err != nil {
		e := fmt.Errorf("%s: Error from db.GetProperty:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}
	switch req.Idx {
	case 1:
		pr.Img1 = imageURL
	case 2:
		pr.Img2 = imageURL
	case 3:
		pr.Img3 = imageURL
	case 4:
		pr.Img4 = imageURL
	case 5:
		pr.Img5 = imageURL
	case 6:
		pr.Img6 = imageURL
	case 7:
		pr.Img7 = imageURL
	case 8:
		pr.Img8 = imageURL
	}
	if err = db.UpdateProperty(r.Context(), &pr); err != nil {
		e := fmt.Errorf("%s: Error from db.UpdateProperty:  %s", funcname, err.Error())
		SvcFuncErrorReturn(w, e, funcname)
		return
	}

	util.Console("URL to image: %s\n", imageURL)
	var resp imgSuccess
	resp.Status = "success"
	resp.URL = imageURL

	SvcWriteResponse(&resp, w)
}

// GeneratePRImageFileName generate file name with PRID
//------------------------------------------------------------------------------
func GeneratePRImageFileName(filename string, PRID, idx int64) string {
	return fmt.Sprintf("img-%d-%d-%s", PRID, idx, filename)
}

// GetImageFilenameFromURL returns the image filename from a URL of the form:
//
//        https://s3.amazonaws.com/directory-pics/img-4-3-roller-32.png
//
// It will return "roller-32.png", basically it strips of the image identifier
// (img), the PRID (4), and the index (3).  Note that any dash embedded in
// the filename is also preserved.
//
// If the supplied filename has less than 3 embedded dashes it will simply
// return the entire filename portion of the url.
//------------------------------------------------------------------------------
func GetImageFilenameFromURL(s string) string {
	_, fname := filepath.Split(s)
	sa := strings.Split(fname, "-")
	if len(sa) < 4 {
		return s
	}
	return strings.Join(sa[3:], "-")
}

// GetFileContentType uses the http library capability to determine a file's type
//
// INPUTS
//    out  -  os.File to read for the type
//
// RETURNS
//    filetype string
//    any error encountered
//------------------------------------------------------------------------------
func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	//--------------------------------------------------------------------------------------
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	//--------------------------------------------------------------------------------------
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// UploadImageFileToS3 is a convenience function for UploadImageToS3.
//
// INPUTS
//     fileHeader - It contains header information of file. e.g., filename, file type
//     usrfile    - It contains file data/information
//
//     PRID       - PRID of the property
//     idx        - index number associated with this pic (1 - 6)
//
// RETURNS
//     filename   - to associate with the contents of usrfile, we don't
//                  actually use this filename to open a file system file.
//     imagePath  - of the form pic{{PRID}{{PRIDX}}
//     imageURL   - URL to access the picture
//     err        - any error encountered
//------------------------------------------------------------------------------
func UploadImageFileToS3(filename string, usrfile *os.File, PRID, idx int64) (string, string, error) {
	//-----------------------------------------
	// get the file size and read
	// the file content into a buffer
	//-----------------------------------------
	fileInfo, _ := usrfile.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	usrfile.Read(buffer)

	return UploadImageToS3(filename, buffer, PRID, idx)
}

// UploadImageToS3 upload image to AWS S3 bucket
//
// INPUTS
//     fileHeader - It contains header information of file. e.g., filename, file type
//     buffer     - Contents of the file as bytes
//
//     PRID       - PRID of the property
//     idx        - index number associated with this pic (1 - 6)
//
// RETURNS
//     filename   - to associate with the contents of buffer, we don't
//                  actually use this filename to open a file system file.
//     imagePath  - of the form pic{{PRID}{{PRIDX}}
//     imageURL   - URL to access the picture
//     err        - any error encountered
//------------------------------------------------------------------------------
func UploadImageToS3(filename string, buffer []byte, PRID, idx int64) (string, string, error) {

	//-----------------------------------------
	// setup credential
	// reading credential from the aws config
	//-----------------------------------------
	creds := credentials.NewStaticCredentials(db.Wdb.Config.S3BucketKeyID, db.Wdb.Config.S3BucketKey, "")
	_, err := creds.Get()
	if err != nil {
		util.Console("Bad credentials: %s", err.Error())
		return "", "", err
	}

	//-----------------------------------------
	// Set up configuration
	//-----------------------------------------
	cfg := aws.NewConfig().WithRegion(getS3Region()).WithCredentials(creds)

	//-----------------------------------------
	// Set up session
	//-----------------------------------------
	sess, err := session.NewSession(cfg)
	if err != nil {
		util.Console("Error creating session: %s\n", err.Error())
		return "", "", err
	}

	//-----------------------------------------
	// Create new s3 instance
	//-----------------------------------------
	svc := s3.New(sess)

	//-----------------------------------------
	// Path after the bucket name
	//-----------------------------------------
	fname := GeneratePRImageFileName(filename, PRID, idx)
	imagePath := filepath.Join(ImageUploadPath, fname)
	contentType := http.DetectContentType(buffer)

	util.Console(`
Bucket:               %s
Key:                  %s
ServerSideEncryption: %s
ContentType:          %s
CacheControl:         %s
ACL:                  %s
`, db.Wdb.Config.S3BucketName, imagePath, "AES256", contentType, "max-age=86400", "public-read")

	//-----------------------------------------
	// define parameters to upload image to S3
	//-----------------------------------------
	params := &s3.PutObjectInput{
		Bucket: aws.String(db.Wdb.Config.S3BucketName),
		Key:    aws.String(imagePath),
		Body:   bytes.NewReader(buffer),
		//ServerSideEncryption: aws.String("AES256"),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("max-age=86400"),
		ACL:          aws.String("public-read"),
	}

	//-----------------------------------------
	// Upload image to s3 bucket
	//-----------------------------------------
	_, err = svc.PutObject(params)
	if err != nil {
		util.Console("PutObject: error: %s", err.Error())
		return "", "", err
	}

	//-----------------------------------------
	// get image location
	//-----------------------------------------
	imageLocation := fmt.Sprintf("%s/%s/%s", db.Wdb.Config.S3BucketHost, db.Wdb.Config.S3BucketName, fname)
	return imagePath, imageLocation, nil
}

// DeleteS3ImageFile  image from AWS S3 bucket
//
// INPUTS
//     filename   - to associate with the contents of usrfile, we don't
//                  actually use this filename to open a file system file.
//     PRID       - PRID of the property
//     idx        - index number associated with this pic (1 - 6)
//
// RETURNS
//     err           - any error encountered
//------------------------------------------------------------------------------
func DeleteS3ImageFile(filename string, PRID, idx int64) error {
	//----------------------------------------------
	// reading credential info from the aws config
	//----------------------------------------------
	creds := credentials.NewStaticCredentials(db.Wdb.Config.S3BucketKeyID, db.Wdb.Config.S3BucketKey, "")
	_, err := creds.Get()
	if err != nil {
		util.Console("Bad credentials: %s", err.Error())
		return err
	}

	//-----------------------------------------
	// Set up configuration
	//-----------------------------------------
	cfg := aws.NewConfig().WithRegion(getS3Region()).WithCredentials(creds)

	//-----------------------------------------
	// Set up session
	//-----------------------------------------
	sess, err := session.NewSession(cfg)
	if err != nil {
		util.Console("Error creating session: %s\n", err.Error())
		return err
	}

	//-----------------------------------------
	// Create new s3 instance
	//-----------------------------------------
	svc := s3.New(sess)
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(db.Wdb.Config.S3BucketName),
		Key:    aws.String(GeneratePRImageFileName(filename, PRID, idx)),
	}
	if _, err := svc.DeleteObject(params); err != nil {
		return err
	}

	return nil
}
