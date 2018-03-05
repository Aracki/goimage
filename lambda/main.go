package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aracki/gohexis/gohexis/bucket"
	"github.com/aracki/gohexis/gohexis/resize"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Proccessing Lambda request %s\n", request.RequestContext.RequestID)

	bucketSrc := "gohexis-source"
	bucketDst := "gohexis-destination"
	imgName := "under_the_sun.jpg"

	var variations []resize.Dimension
	for k, v := range request.QueryStringParameters {
		if k == "dim" {
			width, err := strconv.Atoi(strings.Split(v, "x")[0])
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}
			height, err := strconv.Atoi(strings.Split(v, "x")[1])
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			variations = append(variations, resize.Dimension{
				Width:  width,
				Height: height,
			})
		}
	}

	if len(variations) == 0 {
		return events.APIGatewayProxyResponse{}, errors.New("Dimension params are missing!")
	}

	// todo document s3 configuration
	// returns new Session based on ~/.aws/config & ~/.aws/credentials
	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	})

	// Get image from s3 bucket according to it's name
	imgFile, err := bucket.GetImageFromS3(svc, bucketSrc, imgName)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Resize image
	if _, err = resize.Resize(&imgFile, variations); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Put image on destination bucket
	if err = bucket.PutObjectToS3(svc, bucketDst, imgFile.FullPath); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("not an error"),
		StatusCode: 200,
	}, nil
}

func main() {

	lambda.Start(handler)
}
