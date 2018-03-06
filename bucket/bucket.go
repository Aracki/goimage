package bucket

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/aracki/gohexis/gohexis/resize"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetImageFromS3 gets object from s3 source-bucket by key.
// Encode that object to image interface
// and returns it and it's name.
func GetImageFromS3(svc *s3.S3, bucketName string, key string) (resize.ImageFile, error) {

	ctx := context.Background()
	res, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			log.Println("ErrCodeNoSuchKey occured")
			return resize.ImageFile{}, err
		}
		return resize.ImageFile{}, err
	} else {
		fmt.Printf("%s downloaded from %s\n", key, bucketName)
	}
	defer res.Body.Close()

	img1, _, _ := image.Decode(res.Body)

	return resize.ImageFile{Image: img1, FileName: key}, nil
}

// UploadAllToS3 uploads all the files from pathList to s3 destination-bucket.
func UploadAllToS3(svc *s3.S3, bucketName string, pathList []string) error {

	for _, p := range pathList {

		ctx := context.Background()
		f, err := os.Open(p)
		if err != nil {
			return err
		}

		fileInfo, err := f.Stat()
		if err != nil {
			return err
		}
		size := fileInfo.Size()

		buffer := make([]byte, size)
		if _, err := f.Read(buffer); err != nil {
			return err
		}

		fileBytes := bytes.NewReader(buffer)

		// Uploads the object to S3. The Context will interrupt the request if the
		// timeout expires.
		_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(p),
			Body:   fileBytes,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
				// If the SDK can determine the request or retry delay was canceled
				// by a context the CanceledErrorCode error code will be returned.
				fmt.Fprintf(os.Stderr, "Upload canceled due to timeout, %v\n", err)
			}
			return err
		}

		fmt.Printf("Successfully uploaded file to %s/%s\n", bucketName, p)
		f.Close()
	}

	return nil
}
