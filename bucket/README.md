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
