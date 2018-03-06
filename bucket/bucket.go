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

func GetImageFromS3(svc *s3.S3, bucketName string, fileName string) (resize.ImageFile, error) {

	ctx := context.Background()
	res, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		// Cast err to awserr.Error to handle specific error codes.
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			//todo remove log.fatal
			log.Fatal(aerr.Message())
		}
		return resize.ImageFile{}, err
	} else {
		fmt.Printf("%s downloaded from %s\n", fileName, bucketName)
	}
	defer res.Body.Close()

	img1, _, _ := image.Decode(res.Body)

	return resize.ImageFile{Image: img1, FileName: fileName}, nil
}

func PutObjectToS3(svc *s3.S3, bucketName string, pathList []string) error {

	for _, p := range pathList {

		ctx := context.Background()
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		// copy f to ReadSeeker
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
				fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
			}
			return err
		}

		fmt.Printf("successfully uploaded file to %s/%s\n", bucketName, p)
	}

	return nil
}
