package network

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePrivateRouteTable(ctx *pulumi.Context, projectName string, indexNum string, vpcResource *ec2.Vpc, subnetType string, cidrBlock string, natGatewayResource *ec2.NatGateway) (routeTableResourceObject *ec2.RouteTable, createRouteTableErr error) {
	// TODO: add validations to make sure those are not empty
	routeTableName := fmt.Sprintf("%s-%s-route-table-%s", projectName, subnetType, indexNum)

	routeTableResource, createRouteTableErr := ec2.NewRouteTable(ctx, routeTableName, &ec2.RouteTableArgs{
		VpcId: pulumi.StringInput(vpcResource.ID()),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock:    pulumi.String(cidrBlock),
				NatGatewayId: pulumi.StringInput(natGatewayResource.ID()),
			},
			//&ec2.RouteTableRouteArgs{
			// TODO: check if this ID is consistent between all accounts if not figure out how to take it dynamically
			//DestinationPrefixListId: pulumi.String("pl-6ea54007"),
			//VpcEndpointId:           pulumi.StringInput(vpcEndpointResource.ID()),
			//}, // TODO: add prefix list with destinations for VPC gateway endpoints - can it be done with a variadic function to dynamically pick how many endpoints can be passed, since they will likely be with different types
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String(routeTableName),
		},
	},
		pulumi.Parent(natGatewayResource),
		pulumi.DependsOn([]pulumi.Resource{natGatewayResource}),
	)
	if createRouteTableErr != nil {
		return nil, createRouteTableErr
	}

	return routeTableResource, nil
}

func CreatePublicRouteTable(ctx *pulumi.Context, projectName string, indexNum string, vpcResource *ec2.Vpc, subnetType string, cidrBlock string, inetGatewayResource *ec2.InternetGateway) (routeTableResourceObject *ec2.RouteTable, createRouteTableErr error) {
	// TODO: add validations to make sure those are not empty
	routeTableName := fmt.Sprintf("%s-%s-route-table-%s", projectName, subnetType, indexNum)

	routeTableResource, createRouteTableErr := ec2.NewRouteTable(ctx, routeTableName, &ec2.RouteTableArgs{
		VpcId: pulumi.StringInput(vpcResource.ID()),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String(cidrBlock),
				GatewayId: pulumi.StringInput(inetGatewayResource.ID()),
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String(routeTableName),
		},
	},
		pulumi.Parent(inetGatewayResource),
		pulumi.DependsOn([]pulumi.Resource{inetGatewayResource}),
	)
	if createRouteTableErr != nil {
		return nil, createRouteTableErr
	}

	return routeTableResource, nil
}

func AssociateRouteTable(ctx *pulumi.Context, projectName string, indexNum string, subnetResource *ec2.Subnet, subnetType string, routeTable *ec2.RouteTable) (routeTableAssociationObject *ec2.RouteTableAssociation, associateRouteTableErr error) {
	// TODO: add validations to make sure those are not empty
	routeTableAssocName := fmt.Sprintf("%s-%s-route-table-%s", projectName, subnetType, indexNum)

	routeTableAssociationResource, associateRouteTableErr := ec2.NewRouteTableAssociation(ctx, routeTableAssocName, &ec2.RouteTableAssociationArgs{
		SubnetId:     pulumi.StringInput(subnetResource.ID()),
		RouteTableId: pulumi.StringInput(routeTable.ID()),
	},
		pulumi.Parent(routeTable),
		pulumi.DependsOn([]pulumi.Resource{routeTable}),
	)
	if associateRouteTableErr != nil {
		return nil, associateRouteTableErr
	}

	return routeTableAssociationResource, nil
}
