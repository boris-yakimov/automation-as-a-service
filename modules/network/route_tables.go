package network

import (
	"fmt"
	"reflect"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateRouteTable(ctx *pulumi.Context, projectName string, indexNum string, vpcId pulumi.StringInput, gatewayType string, subnetType string, gatewayId pulumi.StringInput, cidrBlock string) (routeTableResourceObject *ec2.RouteTable, createRouteTableErr error) {
	var gatewayTypeArg string

	if gatewayType == "natgw" {
		gatewayTypeArg = "NatGatewayId"
	} else if gatewayType == "igw" {
		gatewayTypeArg = "GatewayId"
	} else {
		return nil, fmt.Errorf("Invalid value for gatewayType - supported values are NAT and IGW")
	}

	routeTableName := fmt.Sprintf("%s-%s-route-table-%s", projectName, subnetType, indexNum)

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

func AssociateRouteTable(ctx *pulumi.Context, projectName string, indexNum string, routeTableId pulumi.StringInput, subnetId pulumi.StringInput, subnetType string) (routeTableAssociationObject *ec2.RouteTableAssociation, associateRouteTableErr error) {
	routeTableAssocName := fmt.Sprintf("%s-%s-route-table-%s", projectName, subnetType, indexNum)

	routeTableAssociationResource, associateRouteTableErr := ec2.NewRouteTableAssociation(ctx, routeTableAssocName, &ec2.RouteTableAssociationArgs{
		SubnetId:     pulumi.StringInput(subnetId),
		RouteTableId: pulumi.StringInput(routeTableId),
	})
	if associateRouteTableErr != nil {
		return nil, associateRouteTableErr
	}

	return routeTableAssociationResource, nil
}
