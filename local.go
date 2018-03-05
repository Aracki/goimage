package main

import (
	"log"

	"github.com/aracki/gohexis/gohexis/resize"
	"github.com/aracki/gohexis/gohexis/s3aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	bucketSrc := "gohexis-source"
	bucketDst := "gohexis-destination"

	// returns new Session based on ~/.aws/config & ~/.aws/credentials
	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	})

	img, err := s3aws.GetImageFromS3(svc, bucketSrc, "under_the_sun.jpg")

	_, err = resize.Resize(img, "new.jpg")
	if err != nil {
		log.Fatal(err)
	}

	s3aws.PutObjectToS3(svc, bucketDst, "new.jpg")

	if err != nil {
		log.Fatal(err)
	}
}
