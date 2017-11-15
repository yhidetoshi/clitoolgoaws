package clitoolgoaws

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func OutputFormat(data [][]string, resourceType string) {
	table := tablewriter.NewWriter(os.Stdout)

	switch resourceType {
	case EC2:
		table.SetHeader([]string{"tag:Name", "InstanceId", "InstanceType", "AvailabilityZone", "PrivateIp", "PublicIp", "Status"})
	case RDS:
		table.SetHeader([]string{"DBName", "InstanceType", "Status", "Engine", "EngineVersion", "MasterUsername", "DBName", "AvailabilityZone"})
	case ELB:
		table.SetHeader([]string{"ELB_Name", "Scheme", "VPCId", "DNSName"})
	case ELB_INS:
		table.SetHeader([]string{"BackEnd_INstance"})
	case CLOUDWATCH:
		table.SetHeader([]string{"Cloudwatch_Alerm", "MetricName", "Namespace"})
	}

	for _, value := range data {
		table.Append(value)
	}

	table.Render()
}
