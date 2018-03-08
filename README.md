# goImage

## Installation

`go get -d github.com/hexis-hr/goImage...`

## AWS Configuration

And official aws go library used https://github.com/aws/aws-lambda-go

Here is official aws doc: [Programming model for Golang](https://docs.aws.amazon.com/lambda/latest/dg/programming-model-v2.html)

You can generate policies on the [policy generator](http://awspolicygen.s3.amazonaws.com/policygen.html) by selecting type of policy from the drop down and then selecting permissions you need.

#### Lambda function

The Lambda execution role (the one selected on the Configuration tab of the Lambda Console) will need read/write access to call *`S3`* *&* *`CloudWatch Logs`*.

**Limits:** https://docs.aws.amazon.com/lambda/latest/dg/limits.html

#### S3 buckets

Go to AWS S3 console and create source & destination buckets. 

**Note**: Serverless deploy will create one bucket to store application zip and compiled *cloudformation-template.json*.

#### IAM 

Create user with access to: *`AmazonS3FullAccess`*, *`CloudFormation`*, *`CloudWatch Logs`*, *`IAM`*, *`Lambda`* *&* *`STS`*.

#### AWS API Gateway

Go to AWS Lambda console https://console.aws.amazon.com/lambda.

Add triggers on the left sidebar. Than you need to create API Gateway that will integrate with Lambda. It needs to be _Edge Optimized_. Than make `POST` method in _Actions/Create Method_. Select existing Lambda function and check `Use Lambda Proxy integration` to true.

You can make a test call to your method with the provided input.

**Limits:** https://docs.aws.amazon.com/apigateway/latest/developerguide/limits.html

## Deploy

For deploying lambda function you can use one of the following: 

- [Serverless](https://serverless.com)
- [AWS SAM](https://docs.aws.amazon.com/lambda/latest/dg/serverless_app.html)
- [AWS CloudFormation](https://aws.amazon.com/cloudformation/)
- ...

Official documentation on how to deploy Lambda apps [Deploying Lambda Apps]( https://docs.aws.amazon.com/lambda/latest/dg/deploying-lambda-apps.html).

#### With Serverless

- `npm install serverless -g` 
- `make build` 
- deploy app with `serverless deploy` or `sls deploy`

## Usage

Make a POST request on API with following params:

| name of parameter | what is it?                                      |
| ----------------- | :----------------------------------------------- |
| name              | name of picture from source bucket to be resized |
| bucketSrc         | name of source bucket                            |
| bucketDst         | name of destination bucket                       |
| lib               | lib keyword (look at [Algorithm](#algorithms))          |
| filter            | filter keyword (look at [Algorithm](#algorithms))       |

Example of API request:

```
http://[url]/name=under_the_sun.jpg&bucketSrc=gohexis-source&bucketDst=gohexis-destination&alg=imaging&filter=nn
```

Array of variations needs to be sent as JSON to request body:

```
{
    "dim": [
        {
            "w":350,
            "h":300
        },
        {
            "w":1200,
            "h":750
        }
    ]
}
```

If function is successful, it will return status 200 and list of paths:
```
[
  "Thumbnails/350x300/under_the_sun.jpg",
  "Thumbnails/1200x750/under_the_sun.jpg"
]
```

If not, it will return status 4xx, 5xx and message of error:

```
ErrCodeNoSuchKey occurred: NoSuchKey: The specified key does not exist.
status code: 404, request id: 04BC7EA67C9B5407, host id: 4Uz26ywSVspr9BobIF/5yJS9q+8/xq2gP1IpTojZE+wMcecx27ajDhF3EAxXXEqkJs4qa3Quchw=
```

## Algorithms

| library                | query param value |
| ---------------------- | ----------------- |
| disintegration/imaging | imaging           |
| nfnt                   | nfnt              |

| algorithms for disintegration/imaging | query param value |
| :------------------------------------ | :---------------- |
| NearestNeighbor                       | nn                |
| Box                                   | box               |
| Linear                                | linear            |
| MitchellNetravali                     | mn                |
| CatmullRom                            | cr                |
| Gaussian                              | gaussian          |
| Lanczos                               | lan               |

| algorithms for nfnt | query param value |
| :------------------ | :---------------- |
| NearestNeighbor     | nn                |
| Bilinear            | bil               |
| Bicubic             | bic               |
| MitchellNetravali   | mn                |
| Lanczos2            | lan2              |
| Lanczos3            | lan3              |

#### Description

 -  **Lanczos**
    Probably the best resampling filter for photographic images yielding sharp results, but it's slower than cubic filters (see below).

 -  **CatmullRom**
    A sharp cubic filter. It's a good filter for both upscaling and downscaling if sharp results are needed.

- **MitchellNetravali**
    A high quality cubic filter that produces smoother results with less ringing than CatmullRom.

- **Linear**
    Bilinear interpolation filter, produces reasonably good, smooth output. It's faster than cubic filters.

- **Box**
    Simple and fast resampling filter appropriate for downscaling.
    When upscaling it's similar to NearestNeighbor.

- **NearestNeighbor**
    Fastest resample filter, no antialiasing at all. Rarely used.

    ​
