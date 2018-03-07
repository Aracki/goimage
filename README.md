# goImage

## Installation



## AWS Configuration

#### Lambda function

#### S3 buckets

#### AWS API Gateway

You will need to create AWS API Gateway. Add triggers on the left side bar on https://console.aws.amazon.com/lambda.

Than you need to integrate that API Gateway with Lambda.


Than make _POST_  in Actions/Create Method.
Select existing Lambda function. <br>
Check `Use Lambda Proxy integration` to true.

## How To Use

Make _POST_ request on API url with following params:
- name
- bucketSrc
- bucketDst
- alg
- filter

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
