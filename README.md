# monitoring-service


## Table of contents
* [About](#about)
* [Prerequisite](#prerequisite)
* [Technologies](#technologies)
* [Setup](#setup)


## About

Mircoervice which can be run on AWS Lambda to monitor and check websites and metric value of website status code and response time taken per request is sent to AWS CloudWatch. Setup Alarm and SNS as per your preferences can also be done in AWS.

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

To setup in lambda, follow the below steps,

1. Uncomment this line on main.go in line number 33.
```
lambda.Start(lambdaService.StartService)
```

2. Build the project with below command to create a zip file which is necessary to deploy on lambda

```
1. GOOS=linux go build -o main main.go
2. zip monitor-service.zip main 
```

3. Setup the lambda function

- Create a function as follow (with the policy defined under ([policy](#policy)) section) 

![lambda-setup](https://github.com/AungAyeThan/monitoring-service/blob/assets/lambda-setup.png )

- Add the zip file which is generated from earlier in the step 2.

![lambda-deploy](https://github.com/AungAyeThan/monitoring-service/blob/assets/lambda-function.png )

- Add the trigger using eventbridge, define the frequent time to trigger (eg every 5 min, 1 hour etc)

![lambda-trigger](https://github.com/AungAyeThan/monitoring-service/blob/assets/lambda-service.png )

- After that, check the metric log in cloudwatch, under metric section, you will see the metric sent from lambda

![cloudwatch-metric](https://github.com/AungAyeThan/monitoring-service/blob/assets/metric.png )

- here is the sample graph from cloudwatch metric

![cloudwatch-metric-sample](https://github.com/AungAyeThan/monitoring-service/blob/assets/metric-sample.png )

- Add the alarm as per your preferences. You may add SNS together with Alarm

![alarm-setup](https://github.com/AungAyeThan/monitoring-service/blob/assets/alarm-metric.png )





