package api

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Params struct {
	BucketSrc  string
	BucketDst  string
	ImgName    string
	Dimensions []Dimension
}

type ReqBody struct {
	Dim []Dimension `json:"dim"`
}

type Dimension struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

// ProcessParams checks if there are proper dim params (eg. ?dim=200x200&dim350x350).
// Returns an array of Dimension struct.
func Process(request events.APIGatewayProxyRequest, p *Params) (err error) {

	queryParams := request.QueryStringParameters

	if v, ok := queryParams["name"]; ok {
		p.ImgName = v
	} else {
		return fmt.Errorf("missing imgName param")
	}

	if v, ok := queryParams["bucketSrc"]; ok {
		p.BucketSrc = v
	} else {
		return fmt.Errorf("missing bucketSrc param")
	}

	if v, ok := queryParams["bucketDst"]; ok {
		p.BucketDst = v
	} else {
		return fmt.Errorf("missing bucketDst param")
	}

	if request.Body != "" {
		body := ReqBody{}
		json.Unmarshal([]byte(request.Body), &body)

		var allDim []Dimension
		for _, d := range body.Dim {
			allDim = append(allDim, d)
		}
		p.Dimensions = allDim

	} else {
		return fmt.Errorf("missing body")
	}

	return nil
}
