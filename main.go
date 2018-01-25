package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yhidetoshi/clitoolgoaws"
)

const (
	BILLING_REGION = "us-east-1"
)

var (
	argResource   = flag.String("resource", "", "select resource")
	argProfile    = flag.String("profile", "", "slect profile.")
	argRegion     = flag.String("region", "ap-northeast-1", "slect Region")
	argInstances  = flag.String("instances", "", " slect Instance ID or Instance Tag:Name or RDSinstanceName ")
	argELBName    = flag.String("elbname", "", "input elbname")
	argAmiName    = flag.String("aminame", "", "input ami name")
	argAmiId      = flag.String("amiid", "", "input ami id")
	argAMI        = flag.Bool("ami", false, "create ami")
	argAMIList    = flag.Bool("amilist", false, "list ami")
	argStop       = flag.Bool("stop", false, "Instance stop")
	argStart      = flag.Bool("start", false, "Instance start")
	argShow       = flag.Bool("show", false, "show ELB backendend Instances")
	argBilling    = flag.Bool("billing", false, "get billing info")
	argsTerminate = flag.Bool("terminate", false, "Instance terminate")
	argRegister   = flag.Bool("register", false, "Register Instances to ELB")
	argDeregister = flag.Bool("deregister", false, "Deregister Instances to ELB")
)

func main() {
	flag.Parse()

	ec2Client := clitoolgoaws.AwsEC2Client(*argProfile, *argRegion)
	rdsClient := clitoolgoaws.AwsRDSClient(*argProfile, *argRegion)
	elbClient := clitoolgoaws.AwsELBClient(*argProfile, *argRegion)
	cloudwatchClient := clitoolgoaws.AwsCloudwatchClient(*argProfile, *argRegion)
	kinesisClient := clitoolgoaws.AwsKinesisClient(*argProfile, *argRegion)
	iamClient := clitoolgoaws.AwsIAMClient(*argProfile, *argRegion)

	// EC2のコマンド
	var ec2Instances []*string
	var ec2InstancesAMI *string
	exeFlag := true
	if *argResource == "ec2" {
		if *argAMIList {
			clitoolgoaws.ListAMI(ec2Client, nil)
			exeFlag = false
		} else if *argDeregister {
			clitoolgoaws.DeregisterAMI(ec2Client, argAmiId)
			exeFlag = false
		}

		if *argInstances != "" {
			ec2Instances = clitoolgoaws.GetEC2InstanceIds(ec2Client, *argInstances)
			if *argStart {
				clitoolgoaws.ControlEC2Instances(ec2Client, ec2Instances, "start")
			} else if *argStop {
				clitoolgoaws.ControlEC2Instances(ec2Client, ec2Instances, "stop")
			} else if *argsTerminate {
				clitoolgoaws.ControlEC2Instances(ec2Client, ec2Instances, "terminate")
			} else if *argAMI {
				ec2InstancesAMI = clitoolgoaws.GetEC2InstanceIdsAMI(ec2Client, *argInstances)
				clitoolgoaws.RegisterAMI(ec2Client, argAmiName, ec2InstancesAMI)
			} else {
				fmt.Println("`-start` or `-stop` or `-terminate` or `-ami` slect option")
				os.Exit(1)
			}
		} else if exeFlag {
			clitoolgoaws.ListEC2Instances(ec2Client, nil)
		}
	}

	// RDSのコマンド
	var rdsInstances *string
	if *argResource == "rds" {
		if *argInstances != "" {
			rdsInstances = clitoolgoaws.GetRDSInstanceName(rdsClient, *argInstances)
			if *argStart {
				clitoolgoaws.ControlRDSInstances(rdsClient, rdsInstances, "start")
			} else if *argStop {
				clitoolgoaws.ControlRDSInstances(rdsClient, rdsInstances, "stop")
			} else {
				fmt.Println("`-start` or `-stop` slect option")
				os.Exit(1)
			}
		} else {
			clitoolgoaws.ListRDSInstances(rdsClient, nil)
		}

	}

	// ELBのコマンド
	var elasticLoadbalancers []*string
	if *argResource == "elb" {
		if *argELBName != "" {
			elasticLoadbalancers = clitoolgoaws.GetELBInfo(elbClient, *argELBName) //ポインタ
			if *argShow {
				clitoolgoaws.ListELBBackendInstances(elbClient, elasticLoadbalancers, "show")
			} else if *argRegister && *argInstances != "" {
				clitoolgoaws.ControlELB(elbClient, *argELBName, *argInstances, "register")
				clitoolgoaws.ListELBBackendInstances(elbClient, elasticLoadbalancers, "show")
			} else if *argDeregister && *argInstances != "" {
				clitoolgoaws.ControlELB(elbClient, *argELBName, *argInstances, "deregister")
				clitoolgoaws.ListELBBackendInstances(elbClient, elasticLoadbalancers, "show")
			} else {
				fmt.Println("`-show` slect option")
				os.Exit(1)
			}
		} else {
			clitoolgoaws.ListELB(elbClient, nil)
		}
	}

	// Cloudwatchのコマンド
	if *argResource == "cloudwatch" {
		if *argBilling {
			clitoolgoaws.GetBilling(*argProfile, BILLING_REGION)
		} else {
			clitoolgoaws.ListCloudwatch(cloudwatchClient, nil)
		}
	}

	// Kinesisのコマンド
	if *argResource == "kinesis" {
		clitoolgoaws.ListKinesis(kinesisClient, nil)
	}

	// IAM-Userコマンド
	if *argResource == "iam-user" {
		clitoolgoaws.ListIAMUser(iamClient, nil)
	}
	// IAM-Groupコマンド
	if *argResource == "iam-group" {
		clitoolgoaws.ListIAMGroup(iamClient, nil)
	}

}
