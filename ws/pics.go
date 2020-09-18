package ws

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
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
)

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
	_, fname := path.Split(s)
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

// UploadImageFileToS3 upload image to AWS S3 bucket
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
	cfg := aws.NewConfig().WithRegion(S3Region).WithCredentials(creds)

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
	imagePath := path.Join(ImageUploadPath, fname)

	//-----------------------------------------
	// get the file size and read
	// the file content into a buffer
	//-----------------------------------------
	fileInfo, _ := usrfile.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	usrfile.Read(buffer)
	contentType := http.DetectContentType(buffer)

	//-----------------------------------------
	// define parameters to upload image to S3
	//-----------------------------------------
	params := &s3.PutObjectInput{
		Bucket:               aws.String(db.Wdb.Config.S3BucketName),
		Key:                  aws.String(imagePath),   // it include filename
		Body:                 bytes.NewReader(buffer), // data of file
		ServerSideEncryption: aws.String("AES256"),
		ContentType:          aws.String(contentType),
		CacheControl:         aws.String("max-age=86400"),
		ACL:                  aws.String("public-read"),
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
	cfg := aws.NewConfig().WithRegion(S3Region).WithCredentials(creds)

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
		Bucket: aws.String("directory-pics"),
		Key:    aws.String(GeneratePRImageFileName(filename, PRID, idx)),
	}
	if _, err := svc.DeleteObject(params); err != nil {
		return err
	}

	return nil
}
