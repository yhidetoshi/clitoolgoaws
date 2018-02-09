package clitoolgoaws

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

const (
	AS = "as"
)

// client
func AwsASClient(profile string, region string) *autoscaling.AutoScaling {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	asClient := autoscaling.New(ses)

	return asClient
}

// ASGの一覧取得
func ShowAutoScaling(asClient *autoscaling.AutoScaling) {
	allAutoScalingInfo := [][]string{}
	var _instanceCount int64

	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	res, err := asClient.DescribeAutoScalingGroups(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, resInfo := range res.AutoScalingGroups {
		maxsize := strconv.FormatInt(*resInfo.MaxSize, 10)
		minsize := strconv.FormatInt(*resInfo.MinSize, 10)
		desiresize := strconv.FormatInt(*resInfo.DesiredCapacity, 10)
		cooldown := strconv.FormatInt(*resInfo.DefaultCooldown, 10)
		healthchecktime := strconv.FormatInt(*resInfo.HealthCheckGracePeriod, 10)

		for _, instanceInfo := range resInfo.Instances {
			if *instanceInfo.InstanceId != "" {
				_instanceCount++
			}
		}
		instanceCount := strconv.FormatInt(_instanceCount, 10)
		autoscalingInfo := []string{
			*resInfo.AutoScalingGroupName,
			*resInfo.LaunchConfigurationName,
			instanceCount,
			desiresize,
			minsize,
			maxsize,
			cooldown,
			healthchecktime,
			*resInfo.HealthCheckType,
		}
		allAutoScalingInfo = append(allAutoScalingInfo, autoscalingInfo)
		_instanceCount = 0
	}
	OutputFormat(allAutoScalingInfo, AS)
}

func ChangeMaxSizeInstances(asClient *autoscaling.AutoScaling, asgname *string, maxnum *int64) {
	params := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: asgname,
		MaxSize:              maxnum,
	}
	_, err := asClient.UpdateAutoScalingGroup(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Success!!")
}

func ChangeMinSizeInstances(asClient *autoscaling.AutoScaling, asgname *string, minnum *int64) {
	params := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: asgname,
		MinSize:              minnum,
	}
	_, err := asClient.UpdateAutoScalingGroup(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Success!!")
}

func ChangeDesireSizeInstances(asClient *autoscaling.AutoScaling, asgname *string, desirenum *int64) {
	params := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: asgname,
		DesiredCapacity:      desirenum,
	}
	_, err := asClient.UpdateAutoScalingGroup(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Success!!")
}
