package clitoolgoaws

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func OutputFormat(data [][]string, resourceType string) {
	table := tablewriter.NewWriter(os.Stdout)

	switch resourceType {
	case EC2:
		table.SetHeader([]string{"tag:Name", "InstanceId", "InstanceType", "AZ", "PublicIp", "PrivateIp", "Status", "VPCID", "SubnetId", "SG_Name", "SG_Id", "DeviceType", "KeyName"})
	case AMI:
		table.SetHeader([]string{"Name", "ImageId", "OwnerId", "Public", "State", "CreationDate"})
	case RDS:
		table.SetHeader([]string{"DBName", "InstanceType", "Status", "Engine", "EngineVersion", "MasterUsername", "DBName", "AvailabilityZone"})
	case ELB:
		table.SetHeader([]string{"ELB_Name", "Scheme", "VPCId", "DNSName"})
	case ELB_INS:
		table.SetHeader([]string{"BackEnd_INstance"})
	case CLOUDWATCH:
		table.SetHeader([]string{"Cloudwatch_Alerm", "MetricName", "Namespace", "Dimensions", "Period", "THRESHOLD", "Statistic", "AlarmActions", "State"})
	case CLOUDWATCH_BILLING:
		table.SetHeader([]string{"BILLING_(USD)"})
	case KINESIS:
		table.SetHeader([]string{"Stream_Name"})
	case S3BUCKETLIST:
		table.SetHeader([]string{"Bucket_Name"})
	case S3OBJECT:
		table.SetHeader([]string{"Object_Name", "SIZE(KiB)", "StorageClass"})
	case S3BUCKETSIZE:
		table.SetHeader([]string{"Total_SIZE(KiB)"})
	case S3BUCKETSTATUS:
		table.SetHeader([]string{"Bucket", "Status"})
	case IAMUSER:
		table.SetHeader([]string{"username"})
	case IAMGROUP:
		table.SetHeader([]string{"groupname"})
	case EIP:
		table.SetHeader([]string{"publicip", "AllocationId", "InstanceId", "PrivateIpAddress", "Domain", "NetworkInterfaceId", "NetworkInterfaceOwnerId"})
	case SG:
		table.SetHeader([]string{"GroupName", "GroupId", "VpcId"})
	case AS:
		table.SetHeader([]string{"Group_Name", "LaunchConfig_Name", "Instance", "desire", "Min", "Max", "Cooldown", "healthcheck_time", "healthcheck_type"})
	case ROUTE53:
		table.SetHeader([]string{"name", "id", "RecordSetCount", "PrivateZone"})
	case ROUTE53RECORD:
		table.SetHeader([]string{"type", "name", "value", "ttl"})
	}

	for _, value := range data {
		table.Append(value)
	}

	table.Render()
}
