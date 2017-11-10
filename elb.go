package clitoolgoaws

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

const (
	ELB     = "elb"
	ELB_INS = "elb_ins"
)

func AwsELBClient(profile string, region string) *elb.ELB {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	elbClient := elb.New(ses)

	return elbClient
}

// ELBの情報を取得
func ListELB(elbClient *elb.ELB, elbName []*string) {
	params := &elb.DescribeLoadBalancersInput{
		//LoadBalancerName: aws.String(elbName),
		LoadBalancerNames: elbName,
	}
	allLoadbalancers := [][]string{}

	resELBInfo, err := elbClient.DescribeLoadBalancers(params)
	if err != nil {
		os.Exit(1)
	}
	//var backendInstances string

	for _, resInfo := range resELBInfo.LoadBalancerDescriptions {

		loadbalancers := []string{
			*resInfo.LoadBalancerName,
			*resInfo.Scheme,
			*resInfo.VPCId,
			*resInfo.DNSName,
		}

		allLoadbalancers = append(allLoadbalancers, loadbalancers)
	}
	OutputFormat(allLoadbalancers, ELB)
}

// ELB名を取得
func GetELBInfo(elbClient *elb.ELB, elbName string) []*string {
	splitedELBlist := strings.Split(elbName, ",")

	res, err := elbClient.DescribeLoadBalancers(nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var _elbName []*string
	for _, s := range splitedELBlist {
		for _, res := range res.LoadBalancerDescriptions {
			if *res.LoadBalancerName == s {
				elbName = *res.LoadBalancerName
			}
		}
	}
	return _elbName
}

// ELBのBackedInstanceを取得
func ListELBBackendInstances(elbClient *elb.ELB, elbList []*string, operation string) {
	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: elbList,
	}
	resELBInfo, err := elbClient.DescribeLoadBalancers(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	allBackendInstances := [][]string{}
	for _, resInfo := range resELBInfo.LoadBalancerDescriptions {
		for _, backendInstances := range resInfo.Instances {
			backendInstances := []string{
				*backendInstances.InstanceId,
			}
			allBackendInstances = append(allBackendInstances, backendInstances)
		}
	}
	OutputFormat(allBackendInstances, ELB_INS)
}

func RegisterELBInstances(elbClient *elb.ELB, ec2Instances string, elbList string) {
	params := &elb.RegisterInstancesWithLoadBalancerInput{
<<<<<<< HEAD
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(ec2Instances),
			},
=======
			Instances: []*elb.Instance{
				{
					InstanceId: aws.String(ec2Instances),
				},
>>>>>>> 97b3cdcf4c969319ee0d0066803761f165059bae
		},
		LoadBalancerName: aws.String(elbList),
	}
	_, err := elbClient.RegisterInstancesWithLoadBalancer(params)
	if err != nil {
		os.Exit(1)
<<<<<<< HEAD
	} else {
=======
	}else{
>>>>>>> 97b3cdcf4c969319ee0d0066803761f165059bae
		fmt.Println("Success...!")
	}
}

func DeregisterELBInstances(elbClient *elb.ELB, ec2Instances string, elbList string) {
	params := &elb.DeregisterInstancesFromLoadBalancerInput{
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(ec2Instances),
			},
		},
		LoadBalancerName: aws.String(elbList),
	}
	_, err := elbClient.DeregisterInstancesFromLoadBalancer(params)
	if err != nil {
		os.Exit(1)
<<<<<<< HEAD
	} else {
=======
	}else{
>>>>>>> 97b3cdcf4c969319ee0d0066803761f165059bae
		fmt.Println("Success...!")

	}
}

func ControlELB(elbClient *elb.ELB, elbList string, ec2Instances string, operation string) {
	//ListELBBackendInstances(elbClient, elbList)
	fmt.Print("Do you control ELB ?")
	var stdin string
	fmt.Scan(&stdin)

<<<<<<< HEAD
	switch stdin {
	case "y", "Y", "yes":
		switch operation {
=======
	switch stdin{
	case "y", "Y", "yes":
		switch operation{
>>>>>>> 97b3cdcf4c969319ee0d0066803761f165059bae
		case "register":
			fmt.Println("register instances to ELB")
			RegisterELBInstances(elbClient, ec2Instances, elbList)
		case "deregister":
<<<<<<< HEAD
			fmt.Println("deregister instances to ELB")
=======
		    fmt.Println("deregister instances to ELB")
>>>>>>> 97b3cdcf4c969319ee0d0066803761f165059bae
			DeregisterELBInstances(elbClient, ec2Instances, elbList)
		}
	case "n", "N", "no":
		fmt.Println("Exit ...!")
		os.Exit(0)
	default:
		fmt.Println("Exit ...!")
		os.Exit(0)
	}
}
<<<<<<< HEAD
=======


>>>>>>> 97b3cdcf4c969319ee0d0066803761f165059bae
