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
	IAMUSER  = "iam-user"
	IAMGROUP = "iam-group"
)

func AwsIAMClient(profile string, region string) *iam.IAM {

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

func ListIAMUser(iamClient *iam.IAM, userNmaeList *string) {

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
	OutputFormat(allusers, IAMUSER)
}

func ListIAMGroup(iamClient *iam.IAM, userGroupList *string) {

	res, err := iamClient.ListGroups(&iam.ListGroupsInput{
		MaxItems: aws.Int64(10),
	})

	//res, err := iamClient.ListUsers(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	allgroups := [][]string{}

	for _, groupInfo := range res.Groups {
		groups := []string{
			*groupInfo.GroupName,
		}
		allgroups = append(allgroups, groups)
	}
	OutputFormat(allgroups, IAMGROUP)
}
