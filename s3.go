package clitoolgoaws

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	S3           = "s3"
	S3OBJECT     = "s3object"
	S3BUCKETSIZE = "s3bucketsize"
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

// バケットリスト取得
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

// オブジェクトリスト取得
func ListObjects(S3Client *s3.S3, bucketname *string, operation string) {
	params := &s3.ListObjectsInput{
		//Bucket:  aws.String(bucketname),
		Bucket:  bucketname,
		MaxKeys: aws.Int64(100),
	}
	res, err := S3Client.ListObjects(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//fmt.Println(res)
	allObjects := [][]string{}
	var sumObjectSize int64
	for _, resInfo := range res.Contents {
		sumObjectSize += *resInfo.Size
		object := []string{
			*resInfo.Key,
			strconv.FormatInt(*resInfo.Size, 10),
			*resInfo.StorageClass,
		}
		allObjects = append(allObjects, object)
	}
	OutputFormat(allObjects, S3OBJECT)
	// 合計 KiB
	fmt.Println(sumObjectSize)
}

func CalcBucketSize(S3Client *s3.S3, bucketname *string) {
	params := &s3.ListObjectsInput{
		Bucket:  bucketname,
		MaxKeys: aws.Int64(100),
	}
	res, err := S3Client.ListObjects(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	totalSize := [][]string{}
	var sumObjectSize int64

	for _, resInfo := range res.Contents {
		sumObjectSize += *resInfo.Size
	}
	size := strconv.FormatInt(sumObjectSize, 10)
	result := []string{
		size,
	}
	totalSize = append(totalSize, result)

	OutputFormat(totalSize, S3BUCKETSIZE)
}
