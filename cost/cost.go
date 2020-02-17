package cost

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"time"
)

func DailyCost(sess *session.Session) string {

	svc := costexplorer.New(sess)

	param := dailyCostParam()

	costRes, err := svc.GetCostAndUsage(param)

	if err != nil {
		data := map[string]string{"Error": err.Error()}
		jsonData, _ := json.Marshal(data)

		return string(jsonData)
	} else {
		jsonData, _ := json.Marshal(costRes)
		return string(jsonData)
	}

}

func MonthlyCost(sess *session.Session) string {

	svc := costexplorer.New(sess)

	param := monthlyCostParam()

	costRes, err := svc.GetCostAndUsage(param)

	if err != nil {
		data := map[string]string{"Error": err.Error()}
		jsonData, _ := json.Marshal(data)

		return string(jsonData)
	} else {
		jsonData, _ := json.Marshal(costRes)
		return string(jsonData)
	}

}

func TagFilteredMonthlyCost(sess *session.Session, tagName string, tagValue string) string {

	svc := costexplorer.New(sess)

	param := tagFilteredMothlyCostParam(tagName, tagValue)

	costRes, err := svc.GetCostAndUsage(param)

	if err != nil {
		data := map[string]string{"Error": err.Error()}
		jsonData, _ := json.Marshal(data)

		return string(jsonData)
	} else {
		jsonData, _ := json.Marshal(costRes)
		return string(jsonData)
	}

}

func tagFilteredMothlyCostParam(tagName,tagValue string) *costexplorer.GetCostAndUsageInput {
	granularity := "MONTHLY"
	//metric1 := "BlendedCost"
	metric2 := "UnblendedCost"

	//metrics := []*string{&metric1, &metric2}
	metrics := []*string{&metric2}

	t := time.Now()
	//endDate := t.AddDate(0, 0, -t.Day()+1).Format("2006-01-02")     // firstDayOfMonth
	endDate := t.Format("2006-01-02")
	startDate := t.AddDate(0, -12, -t.Day()+1).Format("2006-01-02") // 6 months  ago

	dateInterval := costexplorer.DateInterval{}
	dateInt := &dateInterval
	dateInt = dateInterval.SetEnd(endDate)
	dateInt = dateInterval.SetStart(startDate)
	param := costexplorer.GetCostAndUsageInput{
		Granularity: &granularity,
		Metrics:     metrics,
		TimePeriod:  dateInt,

		Filter: &costexplorer.Expression{
			Tags: &costexplorer.TagValues{
				Key:    aws.String(tagName),
				Values: []*string{aws.String(tagValue)},
			},
		},
	}

	return &param
}

func monthlyCostParam() *costexplorer.GetCostAndUsageInput {
	granularity := "MONTHLY"
	//metric1 := "BlendedCost"
	metric2 := "UnblendedCost"

	//metrics := []*string{&metric1, &metric2}
	metrics := []*string{&metric2}

	t := time.Now()
	endDate := t.AddDate(0, 0, -t.Day()+1).Format("2006-01-02")    // firstDayOfMonth
	startDate := t.AddDate(0, -6, -t.Day()+1).Format("2006-01-02") // 6 months  ago

	dateInterval := costexplorer.DateInterval{}
	dateInt := &dateInterval
	dateInt = dateInterval.SetEnd(endDate)
	dateInt = dateInterval.SetStart(startDate)
	param := costexplorer.GetCostAndUsageInput{
		Granularity: &granularity,
		Metrics:     metrics,
		TimePeriod:  dateInt,
	}

	return &param
}

func dailyCostParam() *costexplorer.GetCostAndUsageInput {
	granularity := "DAILY"
	//metric1 := "BlendedCost"
	metric2 := "UnblendedCost"
	t := time.Now()

	now := time.Now()

	//metrics := []*string{&metric1, &metric2}
	metrics := []*string{&metric2}

	endDate := now.Format("2006-01-02")
	startDate := t.AddDate(0, 0, -t.Day()+1).Format("2006-01-02")
	dateInterval := costexplorer.DateInterval{}
	dateInt := &dateInterval
	dateInt = dateInterval.SetEnd(endDate)
	dateInt = dateInterval.SetStart(startDate)
	param := costexplorer.GetCostAndUsageInput{
		Granularity: &granularity,
		Metrics:     metrics,
		TimePeriod:  dateInt,
	}

	return &param

}
