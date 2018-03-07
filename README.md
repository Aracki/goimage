# goImage

And official aws go library used https://github.com/aws/aws-lambda-go

Here is official aws doc: [Programming model for Golang](https://docs.aws.amazon.com/lambda/latest/dg/programming-model-v2.html) 

## AWS Configuration

#### Lambda function

#### S3 buckets

#### AWS API Gateway

Go to AWS Lambda console  https://console.aws.amazon.com/lambda.

Add triggers on the left sidebar. Than you need to integrate that API Gateway with Lambda.

Than make _POST_  in Actions/Create Method.
Select existing Lambda function. <br>
Check `Use Lambda Proxy integration` to true.

## Deploy

For deploying lambda function you can use one of the following: 

- [Serverless](https://serverless.com)
- [AWS SAM](https://docs.aws.amazon.com/lambda/latest/dg/serverless_app.html)
- [AWS CloudFormation](https://aws.amazon.com/cloudformation/)
- ...

Official documentation on how to deploy Lambda apps [Deploying Lambda Apps]( https://docs.aws.amazon.com/lambda/latest/dg/deploying-lambda-apps.html).

#### Deploy via Serverless

- first install serverless `npm install serverless -g`
- make binaries from cmd/main.go `make build`
- deploy app with `serverless deploy` or `sls deploy`

## Usage

Make _POST_ request on API url with following params:
| name of parameter | about                                            |
| ----------------- | ------------------------------------------------ |
| name              | name of picture from source bucket to be resized |
| bucketSrc         | name of source bucket                            |
| bucketDst         | name of destination bucket                       |
| lib               | lib keyword (given in Algorithm section)         |
| filter            | filter keyword (given in ALgorithm section)      |

Array of variations needs to be sent as JSON to request body. Example:
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

If successful function respond with 200 status and list of paths. Example:
```
[
  "Thumbnails/350x300/under_the_sun.jpg",
  "Thumbnails/1200x750/under_the_sun.jpg"
]
```



### Algorithms 

 -  **Lanczos** (Keyword: l)
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