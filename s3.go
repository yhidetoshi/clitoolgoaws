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
	S3BUCKETLIST = "s3bucketlist"
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
func ListS3Buckets(S3Client *s3.S3) []string {
	params := &s3.ListBucketsInput{}
	res, err := S3Client.ListBuckets(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var bucket []string
	for _, resInfo := range res.Buckets {
		bucket = append(bucket, *resInfo.Name)
	}
	return bucket
}

// 追加中 ===================
// APIconfigで指定している同一リージョンかどうか判定
/*
func CheckRegion(S3CustomClient *s3.S3) {
	bucketlist := ListS3Buckets(S3CustomClient)
	ctx := context.Background()
	for i := 0; i < len(bucketlist); i++ {
		region, err := s3manager.GetBucketRegion(ctx, S3CustomClient, bucketlist[i], "ap-northeast-1")
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
				fmt.Fprintf(os.Stderr, "unable to find bucket %s's region not found\n", bucketlist[i])
			}
		}
		fmt.Printf("Bucket %s is in %s region\n", bucketlist[i])

	}
}
*/

// mainからの呼び出し、結果を出力
func ShowBuckets(S3Client *s3.S3) {
	_bucketlist := ListS3Buckets(S3Client)
	allBuckets := [][]string{}
	for i := 0; i < len(_bucketlist); i++ {
		bucketlist := []string{
			_bucketlist[i],
		}
		allBuckets = append(allBuckets, bucketlist)
	}
	OutputFormat(allBuckets, S3BUCKETLIST)
}

// bucketを指定してオブジェクトリストを出力
func ShowObjects(S3Client *s3.S3, bucketname *string) {
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
	allObjects := [][]string{}
	for _, resInfo := range res.Contents {
		object := []string{
			*resInfo.Key,
			strconv.FormatInt(*resInfo.Size, 10),
			*resInfo.StorageClass,
		}
		allObjects = append(allObjects, object)
	}
	OutputFormat(allObjects, S3OBJECT)
	// 合計 KiB
}

func ShowBucketSize(S3Client *s3.S3, bucketname *string) {
	totalSize := [][]string{}

	_size := CalcBucketSize(S3Client, bucketname)

	size := strconv.FormatInt(_size, 10)
	result := []string{
		size,
	}
	totalSize = append(totalSize, result)
	OutputFormat(totalSize, S3BUCKETSIZE)
}

// Bucketサイズを計算する
func CalcBucketSize(S3Client *s3.S3, bucketname *string) int64 {
	var sumObjectSize int64
	params := &s3.ListObjectsInput{
		Bucket:  bucketname,
		MaxKeys: aws.Int64(100),
	}
	res, err := S3Client.ListObjects(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, resInfo := range res.Contents {
		sumObjectSize += *resInfo.Size
	}
	return sumObjectSize
}

func TotalGetBucketSize(S3Client *s3.S3) {

	var _size int64
	var buffBucket *string
	totalSize := [][]string{}
	allBuckets := ListS3Buckets(S3Client)
	for i := 0; i < len(allBuckets); i++ {
		buffBucket = &allBuckets[i]
		_size += CalcBucketSize(S3Client, buffBucket)
	}
	size := strconv.FormatInt(_size, 10)
	_totalSize := []string{
		size,
	}
	totalSize = append(totalSize, _totalSize)
	OutputFormat(totalSize, S3BUCKETSIZE)
}
