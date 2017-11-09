import clitoolgoaws

func OutputFormat(data [][]string, resourceType string) {
	table := tablewriter.NewWriter(os.Stdout)

	if resourceType == EC2 {
		table.SetHeader([]string{"tag:Name", "InstanceId", "InstanceType", "AvailabilityZone", "PrivateIp", "PublicIp", "Status"})
	}

	if resourceType == RDS {
		table.SetHeader([]string{"DBName", "InstanceType", "Status", "Engine", "EngineVersion", "MasterUsername", "DBName", "AvailabilityZone"})
	}

	for _, value := range data {
		table.Append(value)
	}
	table.Render()
}
