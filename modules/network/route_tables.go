package network

import (
	"fmt"
	"reflect"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// TODO seems route tables are not waiting for the natgateway and therefore failing to create, dependency on nat resource doesn't seem to wait for it
func CreateRouteTable(ctx *pulumi.Context, projectName string, indexNum string, vpcId pulumi.StringInput, gatewayType string, subnetType string, gatewayId pulumi.StringInput, cidrBlock string, gatewayResource interface{}) (routeTableResourceObject *ec2.RouteTable, createRouteTableErr error) {
	var gatewayTypeArg string
	var gatewayObjectType pulumi.Resource

	switch currentGatewayResourceType := gatewayResource.(type) {
	case *ec2.NatGateway:
		gatewayTypeArg = "NatGatewayId"
		gatewayObjectType = currentGatewayResourceType
	case *ec2.InternetGateway:
		gatewayTypeArg = "GatewayId"
		gatewayObjectType = currentGatewayResourceType
	default:
		return nil, fmt.Errorf("Invalid value for gatewayResource - supported types are *ec2.NatGateway and *ec2.InternetGateway ")
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
	},
		pulumi.Parent(gatewayObjectType),
		pulumi.DependsOn([]pulumi.Resource{gatewayObjectType}),
	)
	if createRouteTableErr != nil {
		return nil, createRouteTableErr
	}

	return routeTableResource, nil
}

func AssociateRouteTable(ctx *pulumi.Context, projectName string, indexNum string, subnetId pulumi.StringInput, subnetType string, routeTable *ec2.RouteTable) (routeTableAssociationObject *ec2.RouteTableAssociation, associateRouteTableErr error) {
	routeTableAssocName := fmt.Sprintf("%s-%s-route-table-%s", projectName, subnetType, indexNum)
	//routeTableId := routeTable.ID()

	routeTableAssociationResource, associateRouteTableErr := ec2.NewRouteTableAssociation(ctx, routeTableAssocName, &ec2.RouteTableAssociationArgs{
		SubnetId:     pulumi.StringInput(subnetId),
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
