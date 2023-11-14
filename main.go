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
		inetGwResource, createIgwErr := network.CreateInternetGateway(ctx, vpcId, projectName, "1")
		if createIgwErr != nil {
			return createIgwErr
		}

		inetGwId := inetGwResource.ID()

		//var subnetResource *ec2.Subnet
		// Create VPC Subnets
		for subnetName, cidr := range subnetList {
			var subnetType string
			var gatewayType string

			if strings.Contains(subnetName, "private") {
				subnetType = "private"
				gatewayType = "natgw"
			} else {
				subnetType = "public"
				gatewayType = "igw"
			}

			var createSubnetErr error
			var currentSubnet *ec2.Subnet

			// create subnets
			currentSubnet, createSubnetErr = network.CreateSubnet(ctx, vpcId, subnetType, subnetName, cidr)
			if createSubnetErr != nil {
				return createSubnetErr
			}

			currentSubnetId := currentSubnet.ID()
			indexNum := subnetName[len(subnetName)-1:]

			if subnetType == "private" {
				// TODO: do we really need to create a route table per subnet - maybe create one per public/private type
				// Create a NAT Gateway for each private subnet
				var currentNatGateway *ec2.NatGateway
				var createNatGwErr error

				currentNatGateway, createNatGwErr = network.CreateNatGateway(ctx, vpcId, projectName, indexNum, currentSubnetId, inetGwResource)
				if createNatGwErr != nil {
					return createNatGwErr
				}

				currentNatGwId := currentNatGateway.ID()

				routeTable, createRouteTableErr := network.CreateRouteTable(ctx, projectName, indexNum, vpcId, gatewayType, subnetType, currentNatGwId, cidr)
				if createRouteTableErr != nil {
					return createRouteTableErr
				}

				routeTableId := routeTable.ID()

				_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, routeTableId, currentSubnetId, subnetType)
				if associateRouteTableErr != nil {
					return associateRouteTableErr
				}
			}

			if subnetType == "public" {
				// TODO: do we really need to create a route table per subnet - maybe create one per public/private type
				routeTable, createRouteTableErr := network.CreateRouteTable(ctx, projectName, indexNum, vpcId, gatewayType, subnetType, inetGwId, cidr)
				if createRouteTableErr != nil {
					return createRouteTableErr
				}

				routeTableId := routeTable.ID()

				_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, routeTableId, currentSubnetId, subnetType)
				if associateRouteTableErr != nil {
					return associateRouteTableErr
				}
			}
		}

		// TODO : check what to do with exports and if we need them at all
		ctx.Export("vpcId", vpcId)

		return nil
	})
}
