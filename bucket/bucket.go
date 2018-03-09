package bucket

import (
	"context"
	"fmt"
	"image"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// GetImageFromS3 gets object from s3 source-bucket by key and saves it in-memory.
// Returns Image interface encoded from object.
func GetImageFromS3(svc *s3.S3, bucketName, key string) (image.Image, error) {

	ctx := context.Background()
	res, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			return nil, errors.Wrap(err, "ErrCodeNoSuchKey occurred")
		}
		return nil, err
	}
	defer res.Body.Close()
	fmt.Printf("%s downloaded from %s\n", key, bucketName)

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s decoded to image\n", key)
	return img, nil
}

// UploadAllToS3 uploads all the files from pathList to s3 destination-bucket.
// Key for uploading files are path of file without tmp/. That is it cuts first 5 chars from path.
// Returns keys of successfully uploaded images to source-bucket.
func UploadAllToS3(svc *s3.S3, bucketName string, pathList []string) (keys []string, err error) {

	for _, p := range pathList {

		ctx := context.Background()
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}

		dstKey := "Thumbnails/" + p[5:]

		// Uploads the object to S3. The Context will interrupt the request if the
		// timeout expires.
		_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(dstKey),
			Body:   f,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
				// If the SDK can determine the request or retry delay was canceled
				// by a context the CanceledErrorCode error code will be returned.
				return nil, errors.Wrap(err, "upload canceled due to timeout")
			}
			return nil, err
		}

		fmt.Printf("%s uploaded to %s\n", p, bucketName)
		keys = append(keys, dstKey)
		f.Close()
	}

	return keys, nil
}
