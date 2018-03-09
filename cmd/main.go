package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hexis-hr/goImage/api"
	"github.com/hexis-hr/goImage/bucket"
	"github.com/hexis-hr/goImage/pic"
)

// InitS3 initialize new s3 client.
// You need to set policies for lambda and s3 buckets.
func initS3() *s3.S3 {
	sess := session.Must(session.NewSession())
	return s3.New(sess, &aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	})
}

// Err returns response with error message in body (if error is nil).
// If returns error other than nil, than body will be 'Internal server error' with status 500.
func Err(e error, status int) (events.APIGatewayProxyResponse, error) {
	log.Println(e)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       e.Error(),
	}, nil
}

/*
	- handler must be a function
	- handler may take between 0 and two arguments.
	- if there are two arguments, the first argument must implement "context.Context".
	- handler may return between 0 and two arguments.
	- if there are two return values, the second argument must implement "error".
	- if there is one return value it must implement "error".
*/
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Proccessing Lambda request: %s\n", request.RequestContext.RequestID)

	// Process request and populate Params struct with all query path parameters
	p := api.Params{}
	if err := api.Process(request, &p); err != nil {
		return Err(err, http.StatusBadRequest)
	}

	svc := initS3()

	// Get image from s3 bucket according to it's name and download it to /tmp/ folder
	img, err := bucket.GetImageFromS3(svc, p.BucketSrc, p.ImgName)
	if err != nil {
		return Err(err, http.StatusInternalServerError)
	}

	// Resize image with proper library/filter and create local files under /tmp/ folder
	filePaths, err := pic.Transform(img, p.ImgName, p.Dimensions, p.Subtype, p.Lib, p.Filter)
	if err != nil {
		return Err(err, http.StatusInternalServerError)
	}

	// Put multiple images on destination bucket to proper paths
	keys, err := bucket.UploadAllToS3(svc, p.BucketDst, filePaths)
	if err != nil {
		return Err(err, http.StatusInternalServerError)
	}

	// Marshal all keys to json for response
	jsonResp, err := json.Marshal(keys)
	if err != nil {
		return Err(err, http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonResp),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
