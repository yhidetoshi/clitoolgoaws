
package clitoolgoaws

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-sdk-go/service/kinesis"
)

const (
	KINESIS = "kinesis"
)

func AwsKinesisClient(profile string, region string) *kinesis.Kinesis {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else{
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	kinesisClient := kinesis.New(ses)

	return kinesisClient

}

func ListKinesis(kinesisClient *kinesis.Kinesis, kinesisName *string) {
	params := &kinesis.ListStreamsInput {
		ExclusiveStartStreamName: kinesisName,
	}
	res, err := kinesisClient.ListStreams(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	allstreams := [][]string{}

	for _, resInfo := range res.StreamNames {
		streams := []string{
			*resInfo,
		}
		allstreams = append(allstreams, streams)
	}
	OutputFormat(allstreams, KINESIS)
}
