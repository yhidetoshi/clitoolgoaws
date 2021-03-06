package clitoolgoaws

import (
	"fmt"
	"os"
	"strconv"
	"sync"

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
	CHUNKSIZE      = 2
)

// client
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

// バケットの一覧表示
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

// 指定したバケットのオブジェクトを表示
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

// パブリックアクセス可能なバケットを表示
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

// 指定したバケットのサイズを表示
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

// Publicアクセス可能なバケットを返す
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

// バケットの一覧を返す
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
		//if location == LOCATION {
		if location == LOCATION {
			bucket = append(bucket, *resInfo.Name)
		}
	}
	return bucket
}

//  バケットのロケーションを返す
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

	// us-east-1(バージニア北部)の対処
	if location.LocationConstraint == nil {
		region = "NULL"
	} else {
		region = *location.LocationConstraint
	}
	//region = *location.LocationConstraint

	return region
}

// 指定したバケットを削除 (empty is ok)
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

// 指定したオブジェクトを削除
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

// 指定したバケット内のオブジェクトを削除
func DeleteAllObjects(S3Client *s3.S3, bucketname *string) {
	var buffObject *string
	_objectlist := GetS3Objects(S3Client, bucketname)

	for i := 0; i < len(_objectlist); i++ {
		buffObject = &_objectlist[i]
		DeleteObject(S3Client, bucketname, buffObject)
	}
	fmt.Println("Finish!!")
}

// 指定したバケットのオブジェクト一覧を返す
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

// 指定したBucketのサイズを返す(通常)
// この関数を並列処理する。バケットの配列を受けれるよおうに修正 goroutine用に 新規で作る
func CalcThreadBucketSize(S3Client *s3.S3, bucketname *string, ch chan int64, wg *sync.WaitGroup) {
	//var buffBucket *string
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
	ch <- sumObjectSize
	wg.Done()
	wg.Wait()
}

/*
func Reciver(S3Client *s3.S3, s []string, ch chan int64, wg *sync.WaitGroup) {
	//var buffBucket *string
	//var sum int64

	for i := 0; i < len(s); i++ {
		buffBucket = &s[i]
		// マルチスレッドで呼び出すが先で合計処理しているのはシングルスレッドだから遅い？
		sum += CalcBucketSize(S3Client, buffBucket)
	}
	//ch <- sum
	//wg.Done()
	//wg.Wait()
}
*/

// リージョン内の全バケット合計値を表示
func TotalGetBucketSizeSingle(S3Client *s3.S3) {

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

// リージョン内の全バケット合計値を表示
func TotalGetBucketSize(S3Client *s3.S3) {

	// ------- 並列処理実装中 ------------------ //
	/*
		var chunkNum int // threadの数になる
		var chunkMod int // CHUNKSIZEでの余り
	*/
	var result int64

	//var _size int64
	var buffBucket *string
	totalSize := [][]string{}
	allBuckets := ListS3Buckets(S3Client)

	// -- 並列処理実装中 -- //
	wg := new(sync.WaitGroup)
	ch := make(chan int64, len(allBuckets))
	/*
		chunkNum = len(allBuckets) / CHUNKSIZE
		chunkMod = len(allBuckets) % CHUNKSIZE

		if chunkMod == 0 {
			chunkNum = chunkNum
		} else {
			chunkNum = chunkNum + 1
		}
		wg.Add(chunkNum)
		for i := CHUNKSIZE; len(allBuckets) > 0; {
			if len(allBuckets) < CHUNKSIZE {
				i = len(allBuckets)
			}

			target := allBuckets[:i]
			allBuckets = allBuckets[i:]
			fmt.Println(target)
			go Reciver(S3Client, target, ch, wg)
		}
	*/
	wg.Add(len(allBuckets))
	for i := 0; i < len(allBuckets); i++ {
		buffBucket = &allBuckets[i]
		go CalcThreadBucketSize(S3Client, buffBucket, ch, wg)
	}

	// 合計を受け取りし合算
	for j := 0; j < len(allBuckets); j++ {
		result += <-ch
	}
	//time.Sleep(1 * time.Second)
	close(ch)
	// ------------------------------------------- //

	// ---  並列処理で置き換える箇所 -----//
	/*
		for i := 0; i < len(allBuckets); i++ {
			buffBucket = &allBuckets[i]
			_size += CalcBucketSize(S3Client, buffBucket)
		}
	*/
	// ----------------------------- //

	size := strconv.FormatInt(result, 10)
	_totalSize := []string{
		size,
	}
	totalSize = append(totalSize, _totalSize)
	OutputFormat(totalSize, S3BUCKETSIZE)
}
