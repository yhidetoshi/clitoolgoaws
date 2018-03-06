package clitoolgoaws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	ROUTE53 = "route53"
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
		}
		allHostZones = append(allHostZones, zones)
	}
	OutputFormat(allHostZones, ROUTE53)
}
