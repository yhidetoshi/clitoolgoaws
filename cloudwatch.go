
package clitoolgoaws

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-sdk-go/service/cloudwatch"
)

const (
	CLOUDWATCH = "cloudwatch"
)

func AwscloudwatchClient(profile string, region string) *cloudwatch.CloudWatch{
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else{
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	cloudwatchClient := cloudwatch.New(ses)

	return cloudwatchClient

}

func ListCloudwatch(cloudwatchClient *cloudwatch.CloudWatch, cloudwatchName []*string) {
	params := &cloudwatch.DescribeAlarmsInput {
		AlarmNames: cloudwatchName,
	}
	res, err := cloudwatchClient.DescribeAlarms(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}


	allAlerm := [][]string{}
	var dimensionsInfo string
	var alarmactionsInfo string

	for _, resInfo := range res.MetricAlarms {

		for _, alarmactions := range resInfo.AlarmActions {
			alarmactionsInfo = *alarmactions
		}

		for _, dimensions := range resInfo.Dimensions {
			switch *dimensions.Name {
			case "InstanceId":
				dimensionsInfo = *dimensions.Value
			case "DBInstanceIdentifier":
				dimensionsInfo = *dimensions.Value
			case "StreamName":
				dimensionsInfo = *dimensions.Value
			case "LoadBalancerName":
				dimensionsInfo = *dimensions.Value
			}
		}

		stream := []string {
			*resInfo.AlarmName,
			*resInfo.MetricName,
			*resInfo.Namespace,
			dimensionsInfo,
			alarmactionsInfo,
			*resInfo.StateValue,
		}
		allAlerm = append(allAlerm, stream)
	}
	OutputFormat(allAlerm, CLOUDWATCH)
}
