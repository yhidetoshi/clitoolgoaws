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
	S3BUCKETLIST   = "s3bucketlist"
	S3OBJECT       = "s3object"
	S3BUCKETSIZE   = "s3bucketsize"
	S3BUCKETSTATUS = "s3bucketstatus"
	LIMIT          = 100
	LOCATION       = "ap-northeast-1"
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

func JudgeS3PublicBucket(S3Client *s3.S3, bucketname *string) *string {
	params := &s3.GetBucketAclInput{
		Bucket: bucketname,
	}
	res, err := S3Client.GetBucketAcl(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// バケットがpublicであるかの判定用
	judgePublic := false
	for _, resInfo := range res.Grants {
		if resInfo.Grantee.URI == nil {
			resInfo.Grantee.URI = aws.String("NULL")
		} else if *resInfo.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
			judgePublic = true
		}
	}
	if judgePublic {
		return aws.String("Public")
	} else {
		return aws.String("Private")
	}

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
		// region情報を取得して、ap-northeast-1だけ絞る
		location := GetS3BucketLocation(S3Client, resInfo.Name)
		if location == LOCATION {
			bucket = append(bucket, *resInfo.Name)
		}
	}
	return bucket
}

//  バケットのロケーション取得
func GetS3BucketLocation(S3Client *s3.S3, bucketname *string) string {
	params := &s3.GetBucketLocationInput{
		Bucket: bucketname,
	}
	location, err := S3Client.GetBucketLocation(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var region string
	region = *location.LocationConstraint

	return region
}

// バケット削除 (empty is ok)
func DeleteBucket(S3Client *s3.S3, bucketname *string) {
	params := &s3.DeleteBucketInput{
		Bucket: bucketname,
	}
	_, err := S3Client.DeleteBucket(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Success!!")
}

func DeleteObject(S3Client *s3.S3, bucketname *string, objectname *string) {
	params := &s3.DeleteObjectInput{
		Bucket: bucketname,
		Key:    objectname,
	}
	_, err := S3Client.DeleteObject(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//fmt.Println("Success!!")
}

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

func DeleteAllObjects(S3Client *s3.S3, bucketname *string) {
	var buffObject *string
	_objectlist := GetS3Objects(S3Client, bucketname)

	for i := 0; i < len(_objectlist); i++ {
		buffObject = &_objectlist[i]
		DeleteObject(S3Client, bucketname, buffObject)
	}
	fmt.Println("Finish!!")

}

// bucketを指定してオブジェクトリストと関連する情報を表示
func ShowObjects(S3Client *s3.S3, bucketname *string) {
	pageNum := 0
	var objects []string
	allObjects := [][]string{}

	params := &s3.ListObjectsInput{
		Bucket:  bucketname,
		MaxKeys: aws.Int64(LIMIT),
	}

	S3Client.ListObjectsPages(params, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		pageNum++

		for _, resInfo := range page.Contents {
			objects = []string{
				*resInfo.Key,
				strconv.FormatInt(*resInfo.Size, 10),
				*resInfo.StorageClass,
			}

			allObjects = append(allObjects, objects)
		}
		return true
	})
	OutputFormat(allObjects, S3OBJECT)
	// 合計 KiB
}

func GetS3Objects(S3Client *s3.S3, bucketname *string) []string {
	params := &s3.ListObjectsInput{
		//Bucket:  aws.String(bucketname),
		Bucket:  bucketname,
		MaxKeys: aws.Int64(LIMIT),
	}
	pageNum := 0
	var objects []string

	S3Client.ListObjectsPages(params, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		pageNum++

		for _, resInfo := range page.Contents {
			objects = append(objects, *resInfo.Key)
		}
		return true

	})
	return objects
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
		MaxKeys: aws.Int64(LIMIT),
	}

	pageNum := 0

	S3Client.ListObjectsPages(params, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		pageNum++

		for _, resInfo := range page.Contents {
			sumObjectSize += *resInfo.Size
		}
		return true
	})
	return sumObjectSize
}

func ShowPublicBucket(S3Client *s3.S3) {
	var buffBucket *string
	//var publicBuckets []string
	result := [][]string{}
	allBuckets := ListS3Buckets(S3Client)

	for i := 0; i < len(allBuckets); i++ {
		buffBucket = &allBuckets[i]
		bucketStatus := JudgeS3PublicBucket(S3Client, buffBucket)
		judgedBuckets := []string{
			*buffBucket,
			*bucketStatus,
		}
		result = append(result, judgedBuckets)
	}
	OutputFormat(result, S3BUCKETSTATUS)
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
