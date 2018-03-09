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
	Lib        string
	Filter     string
	Subtype    string
}

type ReqBody struct {
	Dim []Dimension `json:"dim"`
}

type Dimension struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

func processRequired(queryParams map[string]string, body string, p *Params) (err error) {

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

	if v, ok := queryParams["subtype"]; ok {
		p.Subtype = v
	} else {
		return fmt.Errorf("missing subtype param")
	}

	if body != "" {
		reqbody := ReqBody{}
		json.Unmarshal([]byte(body), &reqbody)

		var allDim []Dimension
		for _, d := range reqbody.Dim {
			allDim = append(allDim, d)
		}
		p.Dimensions = allDim
	} else {
		return fmt.Errorf("missing body")
	}

	return nil
}

func processOptional(queryParams map[string]string, p *Params) {

	if v, ok := queryParams["lib"]; ok {
		p.Lib = v
	}

	if v, ok := queryParams["filter"]; ok {
		p.Filter = v
	}
}

// Process checks if there are required and optional params in request.
// Returns an error if there are some missing required parameters.
func Process(request events.APIGatewayProxyRequest, p *Params) (err error) {

	if err := processRequired(request.QueryStringParameters, request.Body, p); err != nil {
		return err
	}

	processOptional(request.QueryStringParameters, p)
	return nil
}
