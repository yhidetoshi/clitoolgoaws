package clitoolgoaws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	IAM = "iam"
)

func AwsIamClient(profile string, region string) *iam.IAM {

	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	iamClient := iam.New(ses)

	return iamClient
}

func ListIamUser(iamClient *iam.IAM, userNmaeList *string) {

	res, err := iamClient.ListUsers(&iam.ListUsersInput{
		MaxItems: aws.Int64(10),
	})

	//res, err := iamClient.ListUsers(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	allusers := [][]string{}

	for _, userInfo := range res.Users {
		users := []string{
			*userInfo.UserName,
		}
		allusers = append(allusers, users)
	}
	OutputFormat(allusers, IAM)
}
