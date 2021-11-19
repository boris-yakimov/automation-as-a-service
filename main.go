package main

import (
	//"fmt"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"automation-as-a-service/network"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		projectName := "eduspire"
		vpcCidrRange := "10.0.0.0/16"
		subnetList := map[string]string{
			"private-subnet1": "10.0.0.0/20", // 4k IPs per subnet
			"private-subnet2": "10.0.32.0/20",
			"private-subnet3": "10.0.64.0/20",

			"public-subnet1": "10.0.128.0/20",
			"public-subnet2": "10.0.160.0/20",
			"public-subnet3": "10.0.192.0/20",
		}

		// Create AWS VPC
		vpcConfig, createVpcErr := network.CreateVPC(ctx, projectName, vpcCidrRange)
		if createVpcErr != nil {
			return createVpcErr
		}

		vpcId := vpcConfig.ID()

		// Create AWS Internet Gateway
		_, createIgwErr := network.CreateInternetGateway(ctx, vpcId, projectName)
		if createIgwErr != nil {
			return createIgwErr
		}

		// Create AWS Subnets
		for subnetName, cidr := range subnetList {
			var subnetType string
			if strings.Contains(subnetName, "private") {
				subnetType = "private"
			} else {
				subnetType = "public"
			}

			_, createSubnetErr := network.CreateSubnet(ctx, vpcId, subnetType, subnetName, cidr)
			if createSubnetErr != nil {
				return createSubnetErr
			}
		}

		return nil
	})
}
