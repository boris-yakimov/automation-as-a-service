package provisioning

import (
	"automation-as-a-service/modules/network"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Network(ctx *pulumi.Context, projectName string, mainRegion string, vpcCidrRange string, subnetList map[string]string) (networkProvisioningError error) {

	// VPC
	vpcResource, createVpcErr := network.CreateVPC(ctx, projectName, vpcCidrRange)
	if createVpcErr != nil {
		return createVpcErr
	}

	// Internet Gateway
	// TODO: what should I do with this hardcoded index number
	inetGwResource, createIgwErr := network.CreateInternetGateway(ctx, projectName, "1", vpcResource)
	if createIgwErr != nil {
		return createIgwErr
	}

	// TODO: check if I can automate handling of request to increase max number of IPs in account - creating EC2 EIP: AddressLimitExceeded: The maximum number of addresses has been reached.

	// Public Subnets - NAT gateway and Route tables
	var natGateways []*ec2.NatGateway
	var privateSubnets []*ec2.Subnet

	for subnetName, cidr := range subnetList {
		var subnetType string

		if strings.Contains(subnetName, "private") {
			subnetType = "private"
		} else {
			subnetType = "public"
		}

		var createSubnetErr error
		var currentSubnet *ec2.Subnet

		currentSubnet, createSubnetErr = network.CreateSubnet(ctx, projectName, subnetType, subnetName, cidr, vpcResource)
		if createSubnetErr != nil {
			return createSubnetErr
		}

		indexNum := subnetName[len(subnetName)-1:]

		if subnetType == "public" {
			currentNatGateway, createNatGwErr := network.CreateNatGateway(ctx, projectName, indexNum, currentSubnet, vpcResource)
			if createNatGwErr != nil {
				return createNatGwErr
			}
			natGateways = append(natGateways, currentNatGateway)

			routeTablePublic, createIgwRouteTableErr := network.CreatePublicRouteTable(ctx, projectName, indexNum, vpcResource, "public", "0.0.0.0/0", inetGwResource)
			if createIgwRouteTableErr != nil {
				return createIgwRouteTableErr
			}

			_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, currentSubnet, "public", routeTablePublic)
			if associateRouteTableErr != nil {
				return associateRouteTableErr
			}
		}

		if subnetType == "private" {
			privateSubnets = append(privateSubnets, currentSubnet)
		}
	}

	// Private Subnets - Route Tables and VPC Endpoints
	for i, subnetResource := range privateSubnets {
		indexNum := strconv.Itoa(i + 1)
		routeTablePrivate, createNatRouteTableErr := network.CreatePrivateRouteTable(ctx, projectName, indexNum, vpcResource, "private", "0.0.0.0/0", natGateways[i])
		if createNatRouteTableErr != nil {
			return createNatRouteTableErr
		}

		_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, subnetResource, "private", routeTablePrivate)
		if associateRouteTableErr != nil {
			return associateRouteTableErr
		}
	}

	_, createS3VpcEndpoint := network.CreateS3VpcEndpoint(ctx, projectName, mainRegion, vpcResource)
	if createS3VpcEndpoint != nil {
		return createS3VpcEndpoint
	}

	_, createDynamoDBVpcEndpoint := network.CreateDynamoDBVpcEndpoint(ctx, projectName, mainRegion, vpcResource)
	if createDynamoDBVpcEndpoint != nil {
		return createDynamoDBVpcEndpoint
	}

	// TODO : check what to do with exports and if we need them at all
	//ctx.Export("vpcResource", vpcResource)
	return nil
}
