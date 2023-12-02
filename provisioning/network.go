package provisioning

import (
	"automation-as-a-service/modules/network"
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
	privateSubnets := make(map[string]*ec2.Subnet)
	natGatewayIdToResourceMap := make(map[string]*ec2.NatGateway)
	var sortedListOfNatIds []string

	// channel for pulumi ApplyT
	applytFuncDone := make(chan bool)

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
		// debug
		//fmt.Printf("Created subnet: %s (%s) with CIDR: %s\n", subnetName, subnetType, cidrRange)
		if createSubnetErr != nil {
			return createSubnetErr
		}

		indexNum := subnetName[len(subnetName)-1:]

		if subnetType == "public" {
			currentNatGateway, createNatGwErr := network.CreateNatGateway(ctx, projectName, indexNum, currentSubnet, vpcResource)
			// debug
			//fmt.Printf("Created NAT Gateway for subnet %s\n", subnetName)
			if createNatGwErr != nil {
				return createNatGwErr
			}
			// get ID of NAT as string from output param of NAT resource
			currentNatGateway.ID().ApplyT(func(id string) error {
				// sort NAT IDs to bypass bug where order of NAT gateways keep changing at random and route tables keep trying to get repointed at every run
				sortedListOfNatIds = append(sortedListOfNatIds, id)
				natGatewayIdToResourceMap[id] = currentNatGateway
				// block channel until ApplyT function finishes as we are trying get NAT id from the function which is by default asynchronous, if we don't wait we hit an issue where list of NAT Ids is empty
				applytFuncDone <- true
				return nil
			})
			// TODO: make sure I understand how goroutines and this specifically works
			<-applytFuncDone

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

	// make sure list of NAT Ids is sorted otherwise route tables keep trying to assign different NATs on every run
	// TODO: make sure to understand how stable slice works here and what is the difference between the regular sort
	sort.SliceStable(sortedListOfNatIds, func(i, j int) bool {
		return sortedListOfNatIds[i] < sortedListOfNatIds[j]
	})

	var listOfPrivateRouteTables []*ec2.RouteTable

	indexNumCreateRoute := "0"
	// Private Subnets - Route Tables and VPC Endpoints
	for _, natId := range sortedListOfNatIds {
		routeTablePrivate, createNatRouteTableErr := network.CreatePrivateRouteTable(ctx, projectName, indexNumCreateRoute, vpcResource, "private", "0.0.0.0/0", natGatewayIdToResourceMap[natId])
		if createNatRouteTableErr != nil {
			return createNatRouteTableErr
		}
		counter, _ := strconv.Atoi(indexNumCreateRoute)
		indexNumCreateRoute = strconv.Itoa(counter + 1)
		// debug
		//fmt.Printf("Created Route Table for CIDR %s\n", cidrRange)
		//routeTablePrivate.ID().ApplyT(func(id string) error {
		//fmt.Printf("route table id: %s\n", id)
		//return nil
		//})

		// Used for VPC endpoint attachments later
		listOfPrivateRouteTables = append(listOfPrivateRouteTables, routeTablePrivate)
	}

	// TODO: how the hell does this work, this is just initialized and than used in for loop bellow but no actual values are assined to it anywhere but actual assignment on route table level in pulumi and aws console seems to work as expected
	sortedRouteTables := make(map[string]*ec2.RouteTable)

	indexNumAssocRoute := "0"
	for cidrRange, routeTable := range sortedRouteTables {
		_, associateRouteTableErr := network.AssociateRouteTable(ctx, projectName, indexNumAssocRoute, privateSubnets[cidrRange], "private", routeTable)
		if associateRouteTableErr != nil {
			return associateRouteTableErr
		}
		counter, _ := strconv.Atoi(indexNumAssocRoute)
		indexNumAssocRoute = strconv.Itoa(counter + 1)
		// debug
		//fmt.Println("Associated Route Table with subnet: ", privateSubnets[cidrRange])
		//fmt.Println(cidrRange + "\n")
		//routeTable.ID().ApplyT(func(id string) error {
		//fmt.Printf("%s -> %s\n", cidrRange, id)
		//return nil
		//})
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
