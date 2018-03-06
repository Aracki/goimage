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

type Params struct {
	bucketSrc  string
	bucketDst  string
	imgName    string
	dimensions []resize.Dimension
}

// ProcessParams checks if there are proper dim params (eg. ?dim=200x200&dim350x350).
// Returns an array of Dimension struct.
func processRequest(request events.APIGatewayProxyRequest, p *Params) (err error) {

	queryParams := request.QueryStringParameters

	if v, ok := queryParams["name"]; ok {
		p.imgName = v
	} else {
		return errors.New(fmt.Sprintf("Missing imgName param"))
	}

	if v, ok := queryParams["bucketSrc"]; ok {
		p.bucketSrc = v
	} else {
		return errors.New(fmt.Sprintf("Missing bucketSrc param"))
	}

	if v, ok := queryParams["bucketDst"]; ok {
		p.bucketDst = v
	} else {
		return errors.New(fmt.Sprintf("Missing bucketDst param"))
	}

	// TODO make possible for multiple dim params
	if v, ok := queryParams["dim"]; ok {

		var d []resize.Dimension

		width, err := strconv.Atoi(strings.Split(v, "x")[0])
		if err != nil {
			return errors.New("Width is not a number")
		}
		height, err := strconv.Atoi(strings.Split(v, "x")[1])
		if err != nil {
			return errors.New("Height is not a number")
		}
		d = append(d, resize.Dimension{
			Width:  width,
			Height: height,
		})

		p.dimensions = d
	} else {
		return errors.New(fmt.Sprintf("Missing dim param"))
	}

	return nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Proccessing Lambda request: %s\n", request.RequestContext.RequestID)

	p := Params{}
	if err := processRequest(request, &p); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	})

	// Get image from s3 bucket according to it's name
	imgFile, err := bucket.GetImageFromS3(svc, p.bucketSrc, p.imgName)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Resize image
	if _, err = resize.Resize(&imgFile, p.dimensions); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Put image on destination bucket
	if err = bucket.PutObjectToS3(svc, p.bucketDst, imgFile.FullPath); err != nil {
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
