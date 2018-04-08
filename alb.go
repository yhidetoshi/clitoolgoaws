package clitoolgoaws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func AwsALBClient(profile string, region string) *elbv2.ELBV2 {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	albClient := elbv2.New(ses)

	return albClient
}

// ALBのインスタンスの状態を取得
func GetALBInstanceInfo(albClient *elbv2.ELBV2, targetgroup *string) {
	var instanceHealthyCount int64
	params := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: targetgroup,
	}
	res, err := albClient.DescribeTargetHealth(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, resInfo := range res.TargetHealthDescriptions {
		if *resInfo.TargetHealth.State != "healthy" {
			instanceHealthyCount++
		}
	}
	fmt.Println(instanceHealthyCount)

}
