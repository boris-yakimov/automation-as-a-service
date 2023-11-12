package network

import (
	"fmt"
	"reflect"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateRouteTable(ctx *pulumi.Context, projectName string, vpcId pulumi.StringInput, gatewayType string, gatewayId pulumi.StringInput, cidrBlock pulumi.StringInput) (routeTableResourceObject *ec2.RouteTable, createRouteTableErr error) {
	routeTableName := projectName + "route_table"

	var gatewayTypeArg string
	if gatewayType == "NAT" {
		gatewayTypeArg = "NatGatewayId"
	} else if gatewayType == "IGW" {
		gatewayTypeArg = "GatewayId"
	} else {
		return nil, fmt.Errorf("Invalid value for gatewayType - supported values are NAT and IGW")
	}

	routeTableRouteArgs := &ec2.RouteTableRouteArgs{
		CidrBlock: pulumi.StringInput(cidrBlock),
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
			//&ec2.RouteTableRouteArgs{
			//CidrBlock:      pulumi.StringInput(cidrBlock),
			//gatewayTypeArg: pulumi.StringInput(gatewayId),
			//},
		},
	})
	if createRouteTableErr != nil {
		return nil, createRouteTableErr
	}
	return routeTableResource, nil
}
