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

var (
	imgName    string
	dimensions []resize.Dimension
	bucketSrc  string
	bucketDst  string
)

// ProcessParams checks if there are proper dim params (eg. ?dim=200x200&dim350x350).
// Returns an array of Dimension struct.
func processParams(params map[string]string) (err error) {

	if v, ok := params["name"]; ok {
		imgName = v
	} else {
		return errors.New(fmt.Sprintf("Missing imgName param"))
	}

	if v, ok := params["bucketSrc"]; ok {
		bucketSrc = v
	} else {
		return errors.New(fmt.Sprintf("Missing bucketSrc param"))
	}

	if v, ok := params["bucketDst"]; ok {
		bucketDst = v
	} else {
		return errors.New(fmt.Sprintf("Missing bucketDst param"))
	}

	if v, ok := params["dim"]; ok {

		width, err := strconv.Atoi(strings.Split(v, "x")[0])
		if err != nil {
			return errors.New("Width is not a number")
		}
		height, err := strconv.Atoi(strings.Split(v, "x")[1])
		if err != nil {
			return errors.New("Height is not a number")
		}

		dimensions = append(dimensions, resize.Dimension{
			Width:  width,
			Height: height,
		})
	} else {
		return errors.New(fmt.Sprintf("Missing dim param"))
	}

	return nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Proccessing Lambda request %s\n", request.RequestContext.RequestID)

	if err := processParams(request.QueryStringParameters); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

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
	if _, err = resize.Resize(&imgFile, dimensions); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Put image on destination bucket
	if err = bucket.PutObjectToS3(svc, bucketDst, imgFile.FullPath); err != nil {
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
