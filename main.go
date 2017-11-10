package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yhidetoshi/clitoolgoaws"
)


var (
	argResource    = flag.String("resource", "", "select resource")
	argProfile     = flag.String("profile", "", "slect profile.")
	argRegion      = flag.String("region", "ap-northeast-1", "slect Region")
	argInstances   = flag.String("instances", "", " slect Instance ID or Instance Tag:Name or RDSinstanceName ")
	argELBName    = flag.String("elbname", "", "input elbname")
	argStop        = flag.Bool("stop", false, "Instance stop")
	argStart       = flag.Bool("start", false, "Instance start")
	argShow         = flag.Bool("show", false,"show ELB backendend Instances")
	argsTerminate  = flag.Bool("terminate", false, "Instance terminate")
)

func main() {
	flag.Parse()

	ec2Client := clitoolgoaws.AwsEC2Client(*argProfile, *argRegion)
	rdsClient := clitoolgoaws.AwsRDSClient(*argProfile, *argRegion)
	elbClient := clitoolgoaws.AwsELBClient(*argProfile, *argRegion)

	// EC2のコマンド
	var ec2Instances []*string
	if *argResource == "ec2" {
		if *argInstances != "" {
			ec2Instances = clitoolgoaws.GetEC2InstanceIds(ec2Client, *argInstances)
			if *argStart {
				clitoolgoaws.ControlEC2Instances(ec2Client, ec2Instances, "start")
			} else if *argStop {
				clitoolgoaws.ControlEC2Instances(ec2Client, ec2Instances, "stop")
			} else if *argsTerminate {
				clitoolgoaws.ControlEC2Instances(ec2Client, ec2Instances, "terminate")
			} else {
				fmt.Println("`-start` or `-stop` or `-terminate` slect option")
				os.Exit(1)
			}
		} else {
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
	//var elbBackendInstances string
	if *argResource == "elb" {
		if *argELBName != "" {
			elasticLoadbalancers = clitoolgoaws.GetELBInfo(elbClient, *argELBName) //ポインタ
			//fmt.Println(elasticLoadbalancers)
			if *argShow {
				clitoolgoaws.ListELBBackendInstances(elbClient, elasticLoadbalancers, "show")
				//clitoolgoaws.ListELBBackendInstances(elbClient, elasticLoadbalancers, "show")
			}else{
				fmt.Println("`-show` slect option")
				os.Exit(1)
			}
		}else{
			clitoolgoaws.ListELB(elbClient, nil)
		}
	}
}
