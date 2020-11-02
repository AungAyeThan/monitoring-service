# monitoring-service


## Table of contents
* [About](#about)
* [Prerequisite](#prerequisite)
* [Technologies](#technologies)
* [Setup](#setup)


## About

Service which can be run on AWS Lambda to monitor and check websites and a metric value of website status code is sent to AWS CloudWatch.

## Prerequisite

Inorder to run the project, you need to have following,

1. golang 1.13 or above (For those who don't have go installed, check out https://golang.org/doc/install)

Optional (In order to run on lambda)
1. AWS account with IAM Role (* [policy](#policy)) with the following policy below for the Lambda Function and cloudwatch metric: (in order to run with Lambda)

### policy

Find the basic preveliage policy for the lambda and cloudwatch metric 
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowMeticAlarmCreation",
            "Effect": "Allow",
            "Action": [
                "cloudwatch:PutMetricAlarm",
                "cloudwatch:PutMetricData"
            ],
            "Resource": "*"
        },
        {
            "Sid": "AllowLogCollection",
            "Effect": "Allow",
            "Action": [
                "logs:PutLogEvents",
                "logs:CreateLogStream",
                "logs:CreateLogGroup"
            ],
            "Resource": "arn:aws:logs:*:*:*"
        }
    ]
}
```
2. AWS SNS Topic -  to send alert and notification



## Technologies

- golang 1.15
- AWS (Lambda, SNS, cloudwatch)

## Setup

To run this project locally, 

```
$ go mod vendor
$ go run main.go
```
