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
	natGateways := make(map[string]*ec2.NatGateway)

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
		//fmt.Printf("Created subnet: %s (%s) with CIDR: %s\n", subnetName, subnetType, cidrRange)
		if createSubnetErr != nil {
			return createSubnetErr
		}

		indexNum := subnetName[len(subnetName)-1:]

		if subnetType == "public" {
			currentNatGateway, createNatGwErr := network.CreateNatGateway(ctx, projectName, indexNum, currentSubnet, vpcResource)
			//fmt.Printf("Created NAT Gateway for subnet %s\n", subnetName)
			if createNatGwErr != nil {
				return createNatGwErr
			}
			// TODO: convert this to a map of cidr to nat to force an order !!!
			listOfNatGateways = append(listOfNatGateways, currentNatGateway)
			natGateways[cidrRange] = currentNatGateway

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

	// Sort list of CIDR ranges and their associated to NAT gateways to prevent random assignment of route table to NAT gateway that constantly detects route table changes
	sortedCidrRanges := make([]string, 0, len(privateSubnets))
	for k := range privateSubnets {
		sortedCidrRanges = append(sortedCidrRanges, k)
	}
	// TODO: make sure to understand how stable slice works here and what is the difference between the regular sort
	sort.SliceStable(sortedCidrRanges, func(i, j int) bool {
		return sortedCidrRanges[i] < sortedCidrRanges[j]
	})

	//natGateways := make(map[string]*ec2.NatGateway)
	//for i, natResource := range listOfNatGateways {
	//cidr := sortedCidrRanges[i]
	//natGateways[cidr] = natResource
	//}
	// debug
	//fmt.Printf("NAT Gateway Assignments: %v\n", natGateways)
	// TODO: fix name
	sortedRouteTables := make(map[string]*ec2.RouteTable)

	var listOfPrivateRouteTables []*ec2.RouteTable

	// Private Subnets - Route Tables and VPC Endpoints
	//for i, cidrRange := range sortedCidrRanges {
	indexNumTemp := "0"
	for cidrRange, natGateway := range natGateways {
		fmt.Println(pulumi.StringInput(natGateway.ID()))
		//indexNum := strconv.Itoa(i + 1)

		// TODO: HERE IS THE PROBLEM BECAUSE NAT KEEPS GETTING ASSIGNED TO A DIFFERENT ROUTE TABLE EVERY TIME !!!
		routeTablePrivate, createNatRouteTableErr := network.CreatePrivateRouteTable(ctx, projectName, indexNumTemp, vpcResource, "private", "0.0.0.0/0", natGateways[cidrRange])
		if createNatRouteTableErr != nil {
			return createNatRouteTableErr
		}
		counter, _ := strconv.Atoi(indexNumTemp)
		indexNumTemp = strconv.Itoa(counter + 1)
		// debug
		//fmt.Printf("Created Route Table for CIDR %s\n", cidrRange)
		//routeTablePrivate.ID().ApplyT(func(id string) error {
		//fmt.Printf("route table id: %s\n", id)
		//return nil
		//})

		listOfPrivateRouteTables = append(listOfPrivateRouteTables, routeTablePrivate)
		sortedRouteTables[cidrRange] = routeTablePrivate
	}

	indexNumTemp1 := "0"
	for cidrRange, routeTable := range sortedRouteTables {
		//indexNum := strconv.Itoa(i + 1)
		routeTableAssocResource, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNumTemp1, privateSubnets[cidrRange], "private", routeTable)
		if associateRouteTableErr != nil {
			return associateRouteTableErr
		}
		counter, _ := strconv.Atoi(indexNumTemp1)
		indexNumTemp1 = strconv.Itoa(counter + 1)
		// debug
		//fmt.Println("Associated Route Table with subnet: ", privateSubnets[cidrRange])
		fmt.Println(cidrRange + "\n")
		routeTable.ID().ApplyT(func(id string) error {
			fmt.Printf("%s -> %s\n", cidrRange, id)
			return nil
		})
		routeTableAssocResource.ID().ApplyT(func(id string) error {
			fmt.Printf("route table assoc id: %s\n", id)
			return nil
		})
	}

	_, createS3VpcEndpoint := network.CreateS3VpcEndpoint(ctx, projectName, mainRegion, vpcResource, listOfPrivateRouteTables)
	if createS3VpcEndpoint != nil {
		return createS3VpcEndpoint
	}

	_, createDynamoDBVpcEndpoint := network.CreateDynamoDBVpcEndpoint(ctx, projectName, mainRegion, vpcResource, listOfPrivateRouteTables)
	if createDynamoDBVpcEndpoint != nil {
		return createDynamoDBVpcEndpoint
	}

	return nil
}
