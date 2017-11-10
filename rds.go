package clitoolgoaws

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

const (
	RDS = "rds"
)

// RDSリソース接続用
func AwsRDSClient(profile string, region string) *rds.RDS {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	rdsClient := rds.New(ses)

	return rdsClient
}

func StopRDSInstances(rdsClient *rds.RDS, rdsInstances *string) {
	params := &rds.StopDBInstanceInput{
		DBInstanceIdentifier: rdsInstances,
	}
	_, err := rdsClient.StopDBInstance(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Success...!")
	}
}

func StartRDSInstances(rdsClient *rds.RDS, rdsInstances *string) {
	params := &rds.StartDBInstanceInput{
		DBInstanceIdentifier: rdsInstances,
	}
	_, err := rdsClient.StartDBInstance(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Success...!")
	}
}

func ListRDSInstances(rdsClient *rds.RDS, rdsInstances *string) {
	params := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: rdsInstances,
	}
	res, err := rdsClient.DescribeDBInstances(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	allDBInstances := [][]string{}

	for _, resInfo := range res.DBInstances {
		instance := []string{
			*resInfo.DBInstanceIdentifier,
			*resInfo.DBInstanceClass,
			*resInfo.DBInstanceStatus,
			*resInfo.Engine,
			*resInfo.EngineVersion,
			*resInfo.MasterUsername,
			*resInfo.DBName,
			*resInfo.AvailabilityZone,
		}
		allDBInstances = append(allDBInstances, instance)
	}
	//fmt.Println(allDBInstances)
	OutputFormat(allDBInstances, RDS)
}

func ControlRDSInstances(rdsClient *rds.RDS, rdsInstances *string, operation string) {
	ListRDSInstances(rdsClient, rdsInstances)

	fmt.Print("Do you control DBInstance ?")
	var stdin string
	fmt.Scan(&stdin)

	switch stdin {
	case "y", "Y", "yes":
		switch operation {
		case "start":
			fmt.Println("start RDS instance")
			StartRDSInstances(rdsClient, rdsInstances)
		case "stop":
			StopRDSInstances(rdsClient, rdsInstances)
		}
	case "n", "N", "no":
		fmt.Println("Exit ...!")
		os.Exit(0)
	default:
		fmt.Println("Exit ...!")
		os.Exit(0)
	}
}

func GetRDSInstanceName(rdsClient *rds.RDS, rdsInstances string) *string {
	splitedInstances := strings.Split(rdsInstances, ",")
	res, err := rdsClient.DescribeDBInstances(nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var instanceName *string
	for _, s := range splitedInstances {
		for _, res := range res.DBInstances {
			if *res.DBInstanceIdentifier == s {
				instanceName = aws.String(rdsInstances)
			}
		}
	}
	return instanceName
}
