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
	//var subnetToNatMapping map[*ec2.Subnet]*ec2.NatGateway
	var natGateways []*ec2.NatGateway
	for subnetName, cidr := range subnetList {
		var subnetType string

		if strings.Contains(subnetName, "private") {
			subnetType = "private"
		} else {
			subnetType = "public"
		}

		var createSubnetErr error
		var currentSubnet *ec2.Subnet

		// create subnets
		currentSubnet, createSubnetErr = network.CreateSubnet(ctx, projectName, subnetType, subnetName, cidr, vpcResource)
		if createSubnetErr != nil {
			return createSubnetErr
		}

		indexNum := subnetName[len(subnetName)-1:]

		//var routeTablePrivate *ec2.RouteTable
		if subnetType == "public" {
			currentNatGateway, createNatGwErr := network.CreateNatGateway(ctx, projectName, indexNum, currentSubnet, vpcResource)
			if createNatGwErr != nil {
				return createNatGwErr
			}
			natGateways = append(natGateways, currentNatGateway)
			//subnetToNatMapping[currentSubnet] = currentNatGateway

			routeTablePublic, createIgwRouteTableErr := network.CreateIgwRouteTable(ctx, projectName, indexNum, vpcResource, "public", "0.0.0.0/0", inetGwResource)
			if createIgwRouteTableErr != nil {
				return createIgwRouteTableErr
			}

			_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, currentSubnet, "public", routeTablePublic)
			if associateRouteTableErr != nil {
				return associateRouteTableErr
			}
		}
	}

	//for currentSubnet, currentNatGateway := range subnetToNatMapping {
	for subnetName, _ := range subnetList {
		if strings.Contains(subnetName, "private") {
			indexNum := subnetName[len(subnetName)-1:]
			subnetId, getSubnetIdErr := network.GetSubnetIdByName(ctx, subnetName)
			if getSubnetIdErr != nil {
				return getSubnetIdErr
			}
			currentSubnet, getCurrentSubnetErr := network.GetSubnetResource(ctx, subnetName, subnetId)
			if getCurrentSubnetErr != nil {
				return getCurrentSubnetErr
			}
			for _, currentNatGateway := range natGateways {
				routeTablePrivate, createNatRouteTableErr := network.CreateNatRouteTable(ctx, projectName, indexNum, vpcResource, "private", "0.0.0.0/0", currentNatGateway)
				if createNatRouteTableErr != nil {
					return createNatRouteTableErr
				}

				// Do nothing for private subnets
				_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, currentSubnet, "private", routeTablePrivate)
				if associateRouteTableErr != nil {
					return associateRouteTableErr
				}
			}
		}
	}

	// TODO : check what to do with exports and if we need them at all
	//ctx.Export("vpcResource", vpcResource)
	return nil
}
