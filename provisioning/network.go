package provisioning

import (
	"automation-as-a-service/modules/network"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Network(ctx *pulumi.Context, projectName string, vpcCidrRange string, subnetList map[string]string) (networkProvisioningError error) {

	// Create AWS VPC
	vpcResource, createVpcErr := network.CreateVPC(ctx, projectName, vpcCidrRange)
	if createVpcErr != nil {
		return createVpcErr
	}

	// Create AWS Internet Gateway
	// TODO: what should I do with this hardcoded index number
	inetGwResource, createIgwErr := network.CreateInternetGateway(ctx, projectName, "1", vpcResource)
	if createIgwErr != nil {
		return createIgwErr
	}

	// TODO: check if I can automate handling of request to increase max number of IPs in account - creating EC2 EIP: AddressLimitExceeded: The maximum number of addresses has been reached.
	// Create VPC Subnets
	for subnetName, cidr := range subnetList {
		var subnetType string
		//var gatewayType string

		if strings.Contains(subnetName, "private") {
			subnetType = "private"
			//gatewayType = "natgw"
		} else {
			subnetType = "public"
			//gatewayType = "igw"
		}

		var createSubnetErr error
		var currentSubnet *ec2.Subnet

		// create subnets
		currentSubnet, createSubnetErr = network.CreateSubnet(ctx, projectName, subnetType, subnetName, cidr, vpcResource)
		if createSubnetErr != nil {
			return createSubnetErr
		}

		indexNum := subnetName[len(subnetName)-1:]

		// TODO: NAT gateways seem to be placed in private subnets, they should be in public
		if subnetType == "private" {
			// TODO: do we really need to create a route table per subnet - maybe create one per public/private type
			// Create a NAT Gateway for each private subnet
			var createNatGwErr error

			currentNatGateway, createNatGwErr := network.CreateNatGateway(ctx, projectName, indexNum, currentSubnet, vpcResource)
			if createNatGwErr != nil {
				return createNatGwErr
			}

			routeTable, createNatRouteTableErr := network.CreateNatRouteTable(ctx, projectName, indexNum, vpcResource, subnetType, "0.0.0.0/0", currentNatGateway)
			if createNatRouteTableErr != nil {
				return createNatRouteTableErr
			}

			_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, currentSubnet, subnetType, routeTable)
			if associateRouteTableErr != nil {
				return associateRouteTableErr
			}
		}

		if subnetType == "public" {
			// TODO: do we really need to create a route table per subnet - maybe create one per public/private type
			routeTable, createIgwRouteTableErr := network.CreateIgwRouteTable(ctx, projectName, indexNum, vpcResource, subnetType, "0.0.0.0/0", inetGwResource)
			if createIgwRouteTableErr != nil {
				return createIgwRouteTableErr
			}

			_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, currentSubnet, subnetType, routeTable)
			if associateRouteTableErr != nil {
				return associateRouteTableErr
			}
		}
	}

	// TODO : check what to do with exports and if we need them at all
	//ctx.Export("vpcResource", vpcResource)
	return nil
}
