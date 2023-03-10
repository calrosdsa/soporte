package util

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func ListObjects(bucket string, sess *session.Session) {
    svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
	fmt.Println("")
}

func UplaodHtmlTemplate(file *os.File,bucket string, sess *session.Session)(err error) {
    uploader := s3manager.NewUploader(sess)
    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key: aws.String(file.Name()),
        Body: file,
		ContentType: aws.String("text/html"),
		ContentEncoding: aws.String("utf-8"),
    })
    if err != nil {
        // Print the error and exit.
		return
    }

    fmt.Printf("Successfully uploaded %q to %q\n", file.Name(), bucket)
	return
}

func UplaodObject(file *multipart.FileHeader,bucket string, sess *session.Session)(url string,err error) {
	src, err := file.Open()
	if err != nil {
		return
	}

	defer src.Close()
	content,err := GetFileContentType(src)
    uploader := s3manager.NewUploader(sess)
	log.Println("content",content)
    output, err := uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key: aws.String("media/"+file.Filename),
        Body: src,
		ContentType: aws.String(content),
		// ContentEncoding: aws.String("utf-8"),

    })
    if err != nil {
        // Print the error and exit.
		return
    }

    fmt.Printf("Successfully uploaded %q to %q\n", file.Filename, bucket)
	return output.Location,nil
}



func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}


func GetFileContentType(ouput multipart.File) (string, error) {

	// to sniff the content type only the first
	// 512 bytes are used.
	buf := make([]byte, 512)
	_, err := ouput.Read(buf)
 
	if err != nil {
	   return "", err
	}
 
	// the function that actually does the trick
	contentType := http.DetectContentType(buf)
 
	   return contentType, nil
 }