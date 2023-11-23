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

	// TODO: check if it makes sense for this to be refactored to happen in each function instead of on provisioning module
	vpcId := vpcResource.ID()

	// Create AWS Internet Gateway
	// TODO: what should I do with this hardcoded index number
	inetGwResource, createIgwErr := network.CreateInternetGateway(ctx, vpcId, projectName, "1", vpcResource)
	if createIgwErr != nil {
		return createIgwErr
	}

	//inetGwId := inetGwResource.ID()

	// TODO: check if I can automate handling of request to increase max number of IPs in account - creating EC2 EIP: AddressLimitExceeded: The maximum number of addresses has been reached.

	//var subnetResource *ec2.Subnet
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
		currentSubnet, createSubnetErr = network.CreateSubnet(ctx, vpcId, projectName, subnetType, subnetName, cidr, vpcResource)
		if createSubnetErr != nil {
			return createSubnetErr
		}

		indexNum := subnetName[len(subnetName)-1:]

		// TODO: make sure that 3 NATs are actually placed in 3 separate AZs
		if subnetType == "private" {
			// TODO: do we really need to create a route table per subnet - maybe create one per public/private type
			// Create a NAT Gateway for each private subnet
			//var currentNatGateway *ec2.NatGateway
			var createNatGwErr error

			currentNatGateway, createNatGwErr := network.CreateNatGateway(ctx, vpcId, projectName, indexNum, currentSubnet, vpcResource)
			if createNatGwErr != nil {
				return createNatGwErr
			}

			routeTable, createNatRouteTableErr := network.CreateNatRouteTable(ctx, projectName, indexNum, vpcId, subnetType, "0.0.0.0/0", currentNatGateway)
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
			routeTable, createIgwRouteTableErr := network.CreateIgwRouteTable(ctx, projectName, indexNum, vpcId, subnetType, "0.0.0.0/0", inetGwResource)
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
	//ctx.Export("vpcId", vpcId)
	return nil
}
