// package util

// import (
// 	"fmt"
// 	"mime/multipart"
// 	"os"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"
// 	"github.com/aws/aws-sdk-go/service/s3/s3manager"
// )

// func ListObjects(bucket string, sess *session.Session) {
//     svc := s3.New(sess)
// 	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
// 	if err != nil {
// 		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
// 	}

// 	for _, item := range resp.Contents {
// 		fmt.Println("Name:         ", *item.Key)
// 		fmt.Println("Last modified:", *item.LastModified)
// 		fmt.Println("Size:         ", *item.Size)
// 		fmt.Println("Storage class:", *item.StorageClass)
// 		fmt.Println("")
// 	}

// 	fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
// 	fmt.Println("")
// }

// func UplaodObject(file *multipart.FileHeader,bucket string, sess *session.Session)(err error) {
// 	src, err := file.Open()
// 	if err != nil {
// 		return
// 	}
// 	defer src.Close()

// 	// Destination
// 	dst, err := os.Create(file.Filename)
// 	if err != nil {
// 		return
// 	}
// 	defer dst.Close()
//     uploader := s3manager.NewUploader(sess)
//     // same as the filename.
//     _, err = uploader.Upload(&s3manager.UploadInput{
//         Bucket: aws.String(bucket),

//         // Can also use the `filepath` standard library package to modify the
//         // filename as need for an S3 object key. Such as turning absolute path
//         // to a relative path.
//         Key: aws.String(""),
//         // The file to be uploaded. io.ReadSeeker is preferred as the Uploader
//         // will be able to optimize memory when uploading large content. io.Reader
//         // is supported, but will require buffering of the reader's bytes for
//         // each part.
//         Body: dst,
//         ContentType: aws.String("text/html"),
//         ContentEncoding: aws.String("utf-8"),
//     })
//     if err != nil {
//         // Print the error and exit.
// 		return
//     }

//     fmt.Printf("Successfully uploaded %q to %q\n", dst.Name(), bucket)
// 	return
// }

// func exitErrorf(msg string, args ...interface{}) {
// 	fmt.Fprintf(os.Stderr, msg+"\n", args...)
// 	os.Exit(1)
// }

package util

import (
	"context"
	"fmt"
	"mime/multipart"
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

// func UplaodHtmlTemplate(file *os.File,bucket string, sess *session.Session)(err error) {
//     uploader := s3manager.NewUploader(sess)
//     _, err = uploader.Upload(&s3manager.UploadInput{
//         Bucket: aws.String(bucket),
//         Key: aws.String(file.Name()),
//         Body: file,
// 		ContentType: aws.String("text/html"),
// 		ContentEncoding: aws.String("utf-8"),
//     })
//     if err != nil {
//         // Print the error and exit.
// 		return
//     }

//     fmt.Printf("Successfully uploaded %q to %q\n", file.Name(), bucket)
// 	return
// }

func UplaodObject(file *multipart.FileHeader,bucket string, sess *session.Session)(url string,err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()
    uploader := s3manager.NewUploader(sess)
    output, err := uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key: aws.String("media/"+file.Filename),
        Body: src,
    })
    if err != nil {
        // Print the error and exit.
		return
    }

    fmt.Printf("Successfully uploaded %q to %q\n", file.Filename, bucket)
	return output.Location,nil
}


func UplaodObjectWebp(ctx context.Context,file *os.File,bucket string, sess *session.Session)(url string,err error) {
	
    uploader := s3manager.NewUploader(sess)
    output, err := uploader.UploadWithContext(ctx,&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key: aws.String("media-t/"+file.Name()),
        Body: file,
    })
    if err != nil {
        // Print the error and exit.
		return
    }

    fmt.Printf("Successfully uploaded %q to %q\n", file.Name(), bucket)
	return output.Location,nil
}


func UplaodObjectWebpWithoutCxt(file *os.File,bucket string, sess *session.Session)(url string,err error) {
	
    uploader := s3manager.NewUploader(sess)
    output, err := uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key: aws.String("media-chat/"+file.Name()),
        Body: file,
    })
    if err != nil {
        // Print the error and exit.
		return
    }

    fmt.Printf("Successfully uploaded %q to %q\n", file.Name(), bucket)
	return output.Location,nil
}


func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

