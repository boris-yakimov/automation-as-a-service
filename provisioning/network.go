package provisioning

import (
	"automation-as-a-service/modules/network"
	"fmt"
	"sort"
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
	var listOfNatGateways []*ec2.NatGateway
	privateSubnets := make(map[string]*ec2.Subnet)

	for subnetName, cidrRange := range subnetList {
		var subnetType string

		if strings.Contains(subnetName, "private") {
			subnetType = "private"
		} else {
			subnetType = "public"
		}

		var createSubnetErr error
		var currentSubnet *ec2.Subnet

		currentSubnet, createSubnetErr = network.CreateSubnet(ctx, projectName, subnetType, subnetName, cidrRange, vpcResource)
		fmt.Printf("Created subnet: %s (%s) with CIDR: %s\n", subnetName, subnetType, cidrRange)
		if createSubnetErr != nil {
			return createSubnetErr
		}

		indexNum := subnetName[len(subnetName)-1:]

		if subnetType == "public" {
			currentNatGateway, createNatGwErr := network.CreateNatGateway(ctx, projectName, indexNum, currentSubnet, vpcResource)
			fmt.Printf("Created NAT Gateway for subnet %s\n", subnetName)
			if createNatGwErr != nil {
				return createNatGwErr
			}
			listOfNatGateways = append(listOfNatGateways, currentNatGateway)

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
			privateSubnets[cidrRange] = currentSubnet
		}
	}

	// Sort list of CIDR ranges and their associated to NAT gateways to prevent random assignment of route table to nat gateway that constantly detects route table changes
	sortedCidrRanges := make([]string, 0, len(privateSubnets))
	for k := range privateSubnets {
		sortedCidrRanges = append(sortedCidrRanges, k)
	}
	// TODO: make sure to understand how stable slice works here and what is the different between the regular sort
	//sort.Strings(sortedCidrRanges)
	sort.SliceStable(sortedCidrRanges, func(i, j int) bool {
		return sortedCidrRanges[i] < sortedCidrRanges[j]
	})

	natGateways := make(map[string]*ec2.NatGateway)
	for i, natResource := range listOfNatGateways {
		cidr := sortedCidrRanges[i]
		natGateways[cidr] = natResource
	}
	fmt.Printf("NAT Gateway Assignments: %v\n", natGateways)
	//fmt.Println(natGateways)

	// VPC endpoints - have to be prepared before private route tables are created
	_, createS3VpcEndpoint := network.CreateS3VpcEndpoint(ctx, projectName, mainRegion, vpcResource)
	if createS3VpcEndpoint != nil {
		return createS3VpcEndpoint
	}

	_, createDynamoDBVpcEndpoint := network.CreateDynamoDBVpcEndpoint(ctx, projectName, mainRegion, vpcResource)
	if createDynamoDBVpcEndpoint != nil {
		return createDynamoDBVpcEndpoint
	}

	// Private Subnets - Route Tables and VPC Endpoints
	for i, cidrRange := range sortedCidrRanges {
		indexNum := strconv.Itoa(i + 1)
		routeTablePrivate, createNatRouteTableErr := network.CreatePrivateRouteTable(ctx, projectName, indexNum, vpcResource, "private", "0.0.0.0/0", natGateways[cidrRange])
		fmt.Printf("Created Route Table for CIDR %s\n", cidrRange)
		if createNatRouteTableErr != nil {
			return createNatRouteTableErr
		}

		_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNum, privateSubnets[cidrRange], "private", routeTablePrivate)
		fmt.Println("Associated Route Table with subnet: ", privateSubnets[cidrRange])
		if associateRouteTableErr != nil {
			return associateRouteTableErr
		}
	}

	// TODO : check what to do with exports and if we need them at all
	ctx.Export("vpcResource", vpcResource)
	return nil
}
