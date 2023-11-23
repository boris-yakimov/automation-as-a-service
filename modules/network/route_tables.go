package network

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// TODO seems route tables are not waiting for the natgateway and therefore failing to create, dependency on nat resource doesn't seem to wait for it
func CreateNatRouteTable(ctx *pulumi.Context, projectName string, indexNum string, vpcId pulumi.StringInput, subnetType string, cidrBlock string, natGatewayResource *ec2.NatGateway) (routeTableResourceObject *ec2.RouteTable, createRouteTableErr error) {
	routeTableName := fmt.Sprintf("%s-%s-route-table-%s", projectName, subnetType, indexNum)

	natGatewayId := natGatewayResource.ID()

	routeTableResource, createRouteTableErr := ec2.NewRouteTable(ctx, routeTableName, &ec2.RouteTableArgs{
		VpcId: pulumi.StringInput(vpcId),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock:    pulumi.String(cidrBlock),
				NatGatewayId: pulumi.StringInput(natGatewayId),
			},
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

func CreateIgwRouteTable(ctx *pulumi.Context, projectName string, indexNum string, vpcId pulumi.StringInput, subnetType string, cidrBlock string, inetGatewayResource *ec2.InternetGateway) (routeTableResourceObject *ec2.RouteTable, createRouteTableErr error) {
	routeTableName := fmt.Sprintf("%s-%s-route-table-%s", projectName, subnetType, indexNum)

	inetGatewayId := inetGatewayResource.ID()

	routeTableResource, createRouteTableErr := ec2.NewRouteTable(ctx, routeTableName, &ec2.RouteTableArgs{
		VpcId: pulumi.StringInput(vpcId),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String(cidrBlock),
				GatewayId: pulumi.StringInput(inetGatewayId),
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
