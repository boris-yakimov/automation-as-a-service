package network

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateRouteTable(ctx *pulumi.Context, projectName string, vpcId pulumi.StringInput, gatewayType string, gatewayId pulumi.StringInput, cidrBlock string) (routeTableResourceObject *ec2.RouteTable, createRouteTableErr error) {
	routeTableName := projectName + "route-table"

	var gatewayTypeArg string
	if gatewayType == "NATGW" {
		gatewayTypeArg = "NatGatewayId"
	} else if gatewayType == "IGW" {
		gatewayTypeArg = "GatewayId"
	} else {
		return nil, fmt.Errorf("Invalid value for gatewayType - supported values are NAT and IGW")
	}

	routeTableRouteArgs := &ec2.RouteTableRouteArgs{
		CidrBlock: pulumi.String(cidrBlock),
	}

	routeTableRouteArgsValue := reflect.ValueOf(routeTableRouteArgs).Elem()
	field := routeTableRouteArgsValue.FieldByName(gatewayTypeArg)

	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(pulumi.StringInput(gatewayId)))
	} else {
		return nil, fmt.Errorf("Invalid gatewayTypeArg when creating route table")
	}

	routeTableResource, createRouteTableErr := ec2.NewRouteTable(ctx, routeTableName, &ec2.RouteTableArgs{
		VpcId: pulumi.StringInput(vpcId),
		Routes: ec2.RouteTableRouteArray{
			routeTableRouteArgs,
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String(routeTableName),
		},
	})
	if createRouteTableErr != nil {
		return nil, createRouteTableErr
	}
	return routeTableResource, nil
}

func AssociateRouteTable(ctx context.Context, routeTableId pulumi.StringInput, subnetId pulumi.StringInput) (routeTableAssociationObject *ec2.RouteTableAssociation, associateRouteTableErr error) {

	routeTableAssociationResource, associateRouteTableErr := ec2.NewRouteTableAssociationResource(ctx, "RouteTableAssociation", &ec2.RouteTableAssociationArgs{
		SubnetId:     pulumi.Any(subnetId),
		routeTableId: pulumi.Any(routeTableId),
	})
	if associateRouteTableErr != nil {
		return nil, associateRouteTableErr
	}

	return routeTableAssociationResource, nil
}
