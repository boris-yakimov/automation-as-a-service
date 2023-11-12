package main

import (
	//"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"automation-as-a-service/network"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		projectName := "pulumi-test"
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
		natGw, createNatGwErr := network.CreateNatGateway(ctx, vpcId, projectName, subnetId, igwResource)
		if createNatGwErr != nil {
			return createNatGwErr
		}

		natGwId := natGw.ID()

		// TODO: get actual CIDR from map above
		var tempCidrRange = pulumi.StringInput("10.0.0.0/20")

		// Create Route Table
		_, createRouteTableErr := network.CreateRouteTable(ctx, projectName, vpcId, "NATGW", natGwId, tempCidrRange)
		if createRouteTableErr != nil {
			return createRouteTableErr
		}

		// TODO : check what to do with exports and if we need them at all
		ctx.Export("vpcId", vpcId)

		return nil
	})
}
