# goImage

## Installation

`go get -d github.com/hexis-hr/goImage...`

## AWS Configuration

You will need to configure AWS services.
Look at [How to configure AWS](https://github.com/hexis-hr/goImage/tree/master/bucket)?

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
|  subtype			     | which type of image transformation is used (look at [Subtypes](#subtypes))	|
| lib               | which library is used (look at [Algorithm](#algorithms)) |
| filter            | which filter/algorithm is used (look at [Algorithm](#algorithms)) |
| quality           | ranges from 1 to 100 inclusive, higher is better. |

Example of API request:

```
http://[url]/name=under_the_sun.jpg&bucketSrc=gohexis-source&bucketDst=gohexis-destination&subtype=resize&alg=imaging&filter=nn
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

## Subtypes

| subtype (type of image transformation) | query param value |
| -------------------------------------- | ----------------- |
| [resize](#resize)                      | resize            |
| [smart crop](#smart-crop)              | smart_crop        |
| [crop](#crop)                          | crop              |

For all subtypes there are required params: 

- name
- bucketSrc
- bucketDst

and required body in JSON format:

- dim array

### Resize

#### Algorithms

| library                                                      | query param value |
| ------------------------------------------------------------ | ----------------- |
| [disintegration/imaging](https://github.com/disintegration/imaging) | imaging           |
| [nfnt](https://github.com/nfnt/resize)                       | nfnt              |

| algorithms for disintegration/imaging | query param value |
| ------------------------------------- | ----------------- |
| NearestNeighbor                       | nn                |
| Box                                   | box               |
| Linear                                | linear            |
| MitchellNetravali                     | mn                |
| CatmullRom                            | cr                |
| Gaussian                              | gaussian          |
| Lanczos                               | lan               |

| algorithms for nfnt | query param value |
| ------------------- | ----------------- |
| NearestNeighbor     | nn                |
| Bilinear            | bil               |
| Bicubic             | bic               |
| MitchellNetravali   | mn                |
| Lanczos2            | lan2              |
| Lanczos3            | lan3              |

### Smart Crop 

Smartcrop finds good image crops for arbitrary sizes. It is using https://github.com/muesli/smartcrop.

### Crop

Crop resize and crop the image to fill the Width and Height area.
