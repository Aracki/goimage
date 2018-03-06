package main

import (
	"encoding/json"
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

// Err returns response with error message in body (if error is nil).
// If returns error than body will be 'Internal server error' with status 500.
func Err(e error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       e.Error(),
	}, nil
}

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
		return Err(err)
	}

	// Resize image with proper algorithm and create local files under /tmp/ folder
	filePaths, err := resize.Resize(&imgFile, p.Dimensions)
	if err != nil {
		return Err(err)
	}

	// Put multiple images on destination bucket to proper paths
	if err = bucket.PutObjectToS3(svc, p.BucketDst, filePaths); err != nil {
		return Err(err)
	}

	// Make filePaths json for response
	jsonResp, err := json.Marshal(filePaths)
	if err != nil {
		return Err(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonResp),
		StatusCode: 200,
	}, nil
}

/*
* handler must be a function
* handler may take between 0 and two arguments.
* if there are two arguments, the first argument must implement "context.Context".
* handler may return between 0 and two arguments.
* if there are two return values, the second argument must implement "error".
* if there is one return value it must implement "error".
 */
func main() {
	lambda.Start(handler)
}
