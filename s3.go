package clitoolgoaws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	S3 = "s3"
)

func AwsS3Client(profile string, region string) *s3.S3 {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	S3Client := s3.New(ses)

	return S3Client
}

// バケットリスト
func ListS3Buckets(S3Client *s3.S3, bucketlist *string) {
	params := &s3.ListBucketsInput{}
	res, err := S3Client.ListBuckets(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//fmt.Println(res)
	allBuckets := [][]string{}
	for _, resInfo := range res.Buckets {
		bucket := []string{
			*resInfo.Name,
		}
		allBuckets = append(allBuckets, bucket)
	}
	OutputFormat(allBuckets, S3)

}
