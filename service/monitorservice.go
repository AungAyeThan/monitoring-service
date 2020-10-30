package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"monitoring-service/utils"
	"net/http"
)

type lambdaService struct {
	cloudWatchSrv 	*cloudwatch.CloudWatch
}

type LambdaService interface {
	StartService()
	MonitorWebsite(url string) (int, error)
	SendMetric(respStatus int, url string) error
}

func NewLambdaService(srv *cloudwatch.CloudWatch) LambdaService{
	return &lambdaService{
		cloudWatchSrv: srv,
	}
}

func (service *lambdaService) StartService() {

	websiteUrls := []string {
		"http://www.google.com",
		"https://www.youtube.com",
	}

	for _, url := range websiteUrls {
		respCode, err := service.MonitorWebsite(url)
		sendMetricErr := service.SendMetric(respCode, url)
		if sendMetricErr != nil {
			fmt.Println(sendMetricErr)
		}
		if respCode == http.StatusOK || respCode == http.StatusCreated {
			fmt.Println("website", url , "is up")
		} else {
			if err != nil {
				fmt.Println("website", url, "is down, reason :", err.Error())
			} else {
				fmt.Println("website", url, "is down due to unexpected error")
			}
		}
	}
}

func (service *lambdaService) MonitorWebsite(url string) (int, error) {

	request, requestErr := http.Get(url)
	if request == nil {
		return http.StatusInternalServerError, requestErr
	}
	if requestErr != nil {
		return request.StatusCode, requestErr
	}

	return request.StatusCode, nil

}

func (service *lambdaService) SendMetric(respStatus int, url string) error {
	_, err := service.cloudWatchSrv.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("Website Status"),
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String(utils.MetricName),
				Value:      aws.Float64(float64(respStatus)),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String(url),
						Value: aws.String("Website Status code"),
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}



}



