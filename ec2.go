package clitoolgoaws

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	EC2 = "ec2"
	AMI = "ami"
	EIP = "eip"
	SG  = "sg"
)

// EC2リソース接続用
func AwsEC2Client(profile string, region string) *ec2.EC2 {
	//ec2Client := ec2.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	ses := session.New(&config)
	ec2Client := ec2.New(ses)

	return ec2Client
}

func ShowElasticIP(ec2Client *ec2.EC2) {
	params := &ec2.DescribeAddressesInput{}
	res, err := ec2Client.DescribeAddresses(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	alleiplist := [][]string{}
	for _, resInfo := range res.Addresses {
		// ValueがNULLだった場合の例外処理
		if resInfo.InstanceId == nil {
			resInfo.InstanceId = aws.String("NULL")
		}
		if resInfo.PrivateIpAddress == nil {
			resInfo.PrivateIpAddress = aws.String("NULL")
		}
		if resInfo.NetworkInterfaceId == nil {
			resInfo.NetworkInterfaceId = aws.String("NULL")
		}
		if resInfo.NetworkInterfaceOwnerId == nil {
			resInfo.NetworkInterfaceOwnerId = aws.String("NULL")
		}
		eiplist := []string{
			*resInfo.PublicIp,
			*resInfo.AllocationId,
			*resInfo.InstanceId,
			*resInfo.PrivateIpAddress,
			*resInfo.Domain,
			*resInfo.NetworkInterfaceId,
			*resInfo.NetworkInterfaceOwnerId,
		}
		alleiplist = append(alleiplist, eiplist)
	}
	OutputFormat(alleiplist, EIP)
}

func DeleteElasticIP(ec2Client *ec2.EC2, allocationid *string) {
	params := &ec2.ReleaseAddressInput{
		AllocationId: allocationid,
	}
	_, err := ec2Client.ReleaseAddress(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Released!!")
}

func ListSecurityGroup(ec2Client *ec2.EC2) {
	params := &ec2.DescribeSecurityGroupsInput{}
	res, err := ec2Client.DescribeSecurityGroups(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	allsglist := [][]string{}
	for _, resInfo := range res.SecurityGroups {
		sglist := []string{
			*resInfo.GroupName,
			*resInfo.GroupId,
			*resInfo.VpcId,
		}
		allsglist = append(allsglist, sglist)
	}
	OutputFormat(allsglist, SG)
}

func ShowSecurityGroup(ec2Client *ec2.EC2, sgid []*string) {
	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: sgid,
	}
	res, err := ec2Client.DescribeSecurityGroups(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(res)

}

func RegisterAMI(ec2Client *ec2.EC2, ec2AMIName *string, ec2Instances *string) {
	reboot := true
	params := &ec2.CreateImageInput{
		InstanceId:  ec2Instances,
		Name:        ec2AMIName,
		NoReboot:    &reboot,
		Description: ec2AMIName,
	}
	res, err := ec2Client.CreateImage(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("success! creating... %s\n", *res.ImageId)
}

func DeregisterAMI(ec2Client *ec2.EC2, ec2AMIid *string) {
	params := &ec2.DeregisterImageInput{
		ImageId: ec2AMIid,
	}
	_, err := ec2Client.DeregisterImage(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Success!!")
}

func ListAMI(ec2Client *ec2.EC2, images []*string) {
	var owner []*string
	var _owner []string = []string{"self"}
	// Convert []string to []*string
	owner = aws.StringSlice(_owner)

	params := &ec2.DescribeImagesInput{
		ImageIds: images,
		Owners:   owner,
	}
	res, err := ec2Client.DescribeImages(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//fmt.Println(res)
	allAmiInfo := [][]string{}
	for _, resInfo := range res.Images {
		amiInfo := []string{
			*resInfo.Name,
			*resInfo.ImageId,
			*resInfo.OwnerId,
			strconv.FormatBool(*resInfo.Public),
			*resInfo.State,
			*resInfo.CreationDate,
		}
		allAmiInfo = append(allAmiInfo, amiInfo)
	}
	OutputFormat(allAmiInfo, EIP)
}

func StopEC2Instances(ec2Client *ec2.EC2, ec2Instances []*string) {
	params := &ec2.StopInstancesInput{
		InstanceIds: ec2Instances,
	}
	res, err := ec2Client.StopInstances(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, r := range res.StoppingInstances {
		fmt.Printf("%s stopped", *r.InstanceId)
	}
}

func StartEC2Instances(ec2Client *ec2.EC2, ec2Instances []*string) {
	params := &ec2.StartInstancesInput{
		InstanceIds: ec2Instances,
	}
	res, err := ec2Client.StartInstances(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, r := range res.StartingInstances {
		fmt.Printf("%s started", *r.InstanceId)
	}
}

func TerminateEC2Instances(ec2Client *ec2.EC2, ec2Instances []*string) {
	params := &ec2.TerminateInstancesInput{
		InstanceIds: ec2Instances,
	}
	res, err := ec2Client.TerminateInstances(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, r := range res.TerminatingInstances {
		fmt.Printf("%s terminated", *r.InstanceId)
	}
}

func ListEC2Instances(ec2Client *ec2.EC2, ec2Instances []*string) {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: ec2Instances,
	}
	res, err := ec2Client.DescribeInstances(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	allInstances := [][]string{}

	for _, resInfo := range res.Reservations {
		for _, instanceInfo := range resInfo.Instances {
			var tagName string
			for _, tagInfo := range instanceInfo.Tags {
				if *tagInfo.Key == "Name" {
					tagName = *tagInfo.Value
				}
			}
			// PublicIpAddressがNULLの場合の例外処理
			if instanceInfo.PublicIpAddress == nil {
				instanceInfo.PublicIpAddress = aws.String("NULL")
			}
			// PrivateIpAddressがNULLの場合の例外処理
			if instanceInfo.PrivateIpAddress == nil {
				instanceInfo.PrivateIpAddress = aws.String("NULL")
			}
			instance := []string{
				tagName,
				*instanceInfo.InstanceId,
				*instanceInfo.InstanceType,
				*instanceInfo.Placement.AvailabilityZone,
				*instanceInfo.PublicIpAddress,
				*instanceInfo.PrivateIpAddress,
				*instanceInfo.State.Name,
				*instanceInfo.VpcId,
				*instanceInfo.SubnetId,
				*instanceInfo.RootDeviceType,
				*instanceInfo.KeyName,
			}
			allInstances = append(allInstances, instance)
		}
	}
	OutputFormat(allInstances, EC2)
}

func ControlEC2Instances(ec2Client *ec2.EC2, ec2Instances []*string, operation string) {
	ListEC2Instances(ec2Client, ec2Instances)

	fmt.Print("Do you control Instance ?")
	var stdin string
	fmt.Scan(&stdin)

	switch stdin {
	case "y", "Y", "yes":
		switch operation {
		case "start":
			fmt.Println("Start EC2 instance")
			StartEC2Instances(ec2Client, ec2Instances)
		case "stop":
			fmt.Println("Stop EC2 instance")
			StopEC2Instances(ec2Client, ec2Instances)
		case "terminate":
			fmt.Println("Terminate EC2 instance")
			TerminateEC2Instances(ec2Client, ec2Instances)
		}
	case "n", "N", "no":
		fmt.Println("Exit ...!")
		os.Exit(0)
	default:
		fmt.Println("Exit ...!")
		os.Exit(0)
	}
}

func GetEC2InstanceIds(ec2Client *ec2.EC2, ec2Instances string) []*string {
	splitedInstances := strings.Split(ec2Instances, ",")
	res, err := ec2Client.DescribeInstances(nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var instanceIds []*string
	for _, s := range splitedInstances {
		for _, r := range res.Reservations {
			for _, i := range r.Instances {
				for _, t := range i.Tags {
					if *t.Key == "Name" {
						if *t.Value == s {
							instanceIds = append(instanceIds, aws.String(*i.InstanceId))
						}
					}
				}
				if *i.InstanceId == s {
					instanceIds = append(instanceIds, aws.String(*i.InstanceId))
				}
			}
		}
	}
	return instanceIds
}

func GetEC2InstanceIdsAMI(ec2Client *ec2.EC2, ec2Instances string) *string {
	splitedInstances := strings.Split(ec2Instances, ",")
	res, err := ec2Client.DescribeInstances(nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var instanceIds *string
	for _, s := range splitedInstances {
		for _, r := range res.Reservations {
			for _, i := range r.Instances {
				for _, t := range i.Tags {
					if *t.Key == "Name" {
						if *t.Value == s {
							instanceIds = aws.String(*i.InstanceId)
						}
					}
				}
				if *i.InstanceId == s {
					instanceIds = aws.String(*i.InstanceId)
				}
			}
		}
	}
	return instanceIds
}
