package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"monitoring-service/service"
	"os"
)



func main() {
	cred := credentials.NewEnvCredentials()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: cred,
	})

	if err != nil {
		fmt.Printf("error in getting session %s ",err.Error())
	}
	svc := cloudwatch.New(sess)

	lambdaService := service.NewLambdaService(svc)

	//uncomment below line to run in local or apart from lambda
	lambdaService.StartService()

	//uncomment below like to run in lambda
	//lambda.Start(lambdaService.StartService)
}