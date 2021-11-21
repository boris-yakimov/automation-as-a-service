package main

import (
	//"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
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
		vpcResource, createVpcErr := network.CreateVPC(ctx, projectName, vpcCidrRange)
		if createVpcErr != nil {
			return createVpcErr
		}

		vpcId := vpcResource.ID()

		// Create AWS Internet Gateway
		igwResource, createIgwErr := network.CreateInternetGateway(ctx, vpcId, projectName)
		if createIgwErr != nil {
			return createIgwErr
		}

		var subnetResource *ec2.Subnet
		// Create AWS Subnets
		for subnetName, cidr := range subnetList {
			var subnetType string
			if strings.Contains(subnetName, "private") {
				subnetType = "private"
			} else {
				subnetType = "public"
			}

			var createSubnetErr error
			// TODO: create a map of subnet ids/names to use in later associations
			if subnetName == "public-subnet1" {
				subnetResource, createSubnetErr = network.CreateSubnet(ctx, vpcId, subnetType, subnetName, cidr)
			} else {
				_, createSubnetErr = network.CreateSubnet(ctx, vpcId, subnetType, subnetName, cidr)
			}
			if createSubnetErr != nil {
				return createSubnetErr
			}
		}
		// TODO: remove after converting this to a map of IDs
		subnetId := subnetResource.ID()

		// Create NAT Gateway
		// TODO: optional configure of how many NATGWs we want - specify cost implications
		_, createNatGwErr := network.CreateNatGateway(ctx, vpcId, projectName, subnetId, igwResource)
		if createNatGwErr != nil {
			return createNatGwErr
		}

		return nil
	})
}
