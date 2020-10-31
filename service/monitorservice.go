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
	MonitorWebsite(url string) (int, error)
	SendMetric(respStatus int, url string) error
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
		respCode, err := service.MonitorWebsite(url)
		sendMetricErr := service.SendMetric(respCode, url)
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
	startTIme := time.Now()
	websiteUrls := []string {
		"https://www.youtube.com",
		"http://google.com",
		"https://www.facebook.com",
		"https://www.gmail.com",
		"https://www.instagram.com",
		"https://oway.com.mm",
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

	endTime := time.Since(startTIme)
	fmt.Println("time taken", endTime)
}

func (service *lambdaService) MonitorWebsite(url string) (int, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		} }

	resp, err := client.Get(url)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
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
	return err
}



