package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"monitoring-service/utils"
	"net/http"
	"sync"
	"time"
)

type lambdaService struct {
	cloudWatchSrv 	*cloudwatch.CloudWatch
}

type LambdaService interface {
	StartService()
	MonitorWebsite(url string) (int, error, float64)
	SendMetric(respStatus int, url string, timeTaken float64) error
}

func NewLambdaService(srv *cloudwatch.CloudWatch) LambdaService{
	return &lambdaService{
		cloudWatchSrv: srv,
	}
}

func (service *lambdaService) worker(input chan string, output chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	var result string
	for url := range input {
		respCode, err, timeTaken := service.MonitorWebsite(url)
		sendMetricErr := service.SendMetric(respCode, url, timeTaken)
		if sendMetricErr != nil {
			result = sendMetricErr.Error()
		}
		if err != nil {
			result = fmt.Sprintf("connectivity error in %s, reason : %s", url, err.Error())
		}
		result = fmt.Sprintf("website name %s, status %d", url, respCode)

		output <- result
	}
}

func (service *lambdaService) StartService() {
	//add url that you would like to check
	//make sure to add http scheme otherwise it will return error
	websiteUrls := []string {
		"http://google.com",
		"https://www.facebook.com",
		"https://www.gmail.com",
	}
	var wg sync.WaitGroup

	input := make(chan string, len(websiteUrls))
	output := make(chan string, len(websiteUrls))
	workers := utils.WorkerLimit

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go service.worker(input, output, &wg)
	}

	for _, job := range websiteUrls {
		input <- job
	}

	close(input)
	wg.Wait()
	close(output)

	// Read from output channel
	for result := range output {
		fmt.Println(result)
	}
}

func (service *lambdaService) MonitorWebsite(url string) (int, error, float64) {
	callStartTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, err, 0
	}
	callEndTime := time.Since(callStartTime).Seconds()
	fmt.Println("url - ", url, "call time taken -", callEndTime, "sec")

	return resp.StatusCode, nil, callEndTime
}

func (service *lambdaService) SendMetric(respStatus int, url string, timeTaken float64) error {
	_, err := service.cloudWatchSrv.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String(utils.MetricStatusNameSpace),
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
			&cloudwatch.MetricDatum{
				MetricName: aws.String(utils.MetricDurationNameSpace),
				Value:      aws.Float64(timeTaken),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String(url),
						Value: aws.String("Response Time"),
					},
				},
			},
		},
	})
	return err
}



