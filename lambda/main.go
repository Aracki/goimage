package main

import (
	"fmt"
	"log"

	"github.com/aracki/gohexis/gohexis/api"
	"github.com/aracki/gohexis/gohexis/bucket"
	"github.com/aracki/gohexis/gohexis/resize"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Proccessing Lambda request: %s\n", request.RequestContext.RequestID)

	p := api.Params{}
	if err := api.Process(request, &p); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	})

	// Get image from s3 bucket according to it's name and download it to /tmp/ folder
	imgFile, err := bucket.GetImageFromS3(svc, p.BucketSrc, p.ImgName)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Resize image with proper algorithm and create local files under /tmp/ folder
	filePaths, err := resize.Resize(&imgFile, p.Dimensions)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Put multiple images on destination bucket to proper paths
	if err = bucket.PutObjectToS3(svc, p.BucketDst, filePaths); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Lambda function successfully executed!"),
		StatusCode: 200,
	}, nil
}

func main() {

	lambda.Start(handler)
}
