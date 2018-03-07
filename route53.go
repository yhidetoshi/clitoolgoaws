package clitoolgoaws

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	ROUTE53       = "route53"
	ROUTE53RECORD = "route53record"
)

func AwsRoute53Client(profile string, region string) *route53.Route53 {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}

	ses := session.New(&config)
	route53Client := route53.New(ses)

	return route53Client
}

func ShowHostedZone(route53Client *route53.Route53) {
	params := &route53.ListHostedZonesInput{}
	res, err := route53Client.ListHostedZones(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	allHostZones := [][]string{}
	for _, resInfo := range res.HostedZones {
		zones := []string{
			*resInfo.Name,
			*resInfo.Id,
			strconv.FormatInt(*resInfo.ResourceRecordSetCount, 10),
			strconv.FormatBool(*resInfo.Config.PrivateZone),
		}
		allHostZones = append(allHostZones, zones)
	}
	OutputFormat(allHostZones, ROUTE53)
}

func ShowListResourceRecordSets(route53Client *route53.Route53, zoneId *string) {
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: zoneId,
	}
	res, err := route53Client.ListResourceRecordSets(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	allRecords := [][]string{}
	var value string
	for _, resInfo := range res.ResourceRecordSets {

		for _, valuInfo := range resInfo.ResourceRecords {
			value = *valuInfo.Value
		}

		records := []string{
			*resInfo.Type,
			*resInfo.Name,
			value,
			strconv.FormatInt(*resInfo.TTL, 10),
		}
		allRecords = append(allRecords, records)
	}
	OutputFormat(allRecords, ROUTE53RECORD)
}
