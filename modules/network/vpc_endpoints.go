package network

import (
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateS3VpcEndpoint(ctx *pulumi.Context, projectName string, mainRegion string, vpcResource *ec2.Vpc, listOfPrivateRouteTables []*ec2.RouteTable) (vpcEndpointResource *ec2.VpcEndpoint, createS3VpcEndpointErr error) {

	if len(listOfPrivateRouteTables) < 2 {
		currentLen := strconv.Itoa(len(listOfPrivateRouteTables))
		return nil, fmt.Errorf("listOfPrivateRouteTables should contain exactly 3 route tables, it currently contains %s", currentLen)
	}

	vpcEndpointName := fmt.Sprintf("%s-s3-vpc-gateway-endpoint", projectName)
	s3ServiceName := fmt.Sprintf("com.amazonaws.%s.s3", mainRegion)

	_, createS3VpcEndpointErr = ec2.NewVpcEndpoint(ctx, vpcEndpointName, &ec2.VpcEndpointArgs{
		VpcId:       pulumi.StringInput(vpcResource.ID()),
		ServiceName: pulumi.String(s3ServiceName),
		RouteTableIds: pulumi.StringArray{
			pulumi.StringInput(listOfPrivateRouteTables[0].ID()),
			pulumi.StringInput(listOfPrivateRouteTables[1].ID()),
			pulumi.StringInput(listOfPrivateRouteTables[2].ID()),
		},
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(vpcEndpointName),
			"ManagedBy": pulumi.String("pulumi"),
		},
	})

	if createS3VpcEndpointErr != nil {
		return nil, createS3VpcEndpointErr
	}

	return vpcEndpointResource, nil
}

func CreateDynamoDBVpcEndpoint(ctx *pulumi.Context, projectName string, mainRegion string, vpcResource *ec2.Vpc, listOfPrivateRouteTables []*ec2.RouteTable) (vpcEndpointResource *ec2.VpcEndpoint, createDynamoVpcEndpointErr error) {

	if len(listOfPrivateRouteTables) < 2 {
		currentLen := strconv.Itoa(len(listOfPrivateRouteTables))
		return nil, fmt.Errorf("listOfPrivateRouteTables should contain exactly 3 route tables, it currently contains %s", currentLen)
	}

	vpcEndpointName := fmt.Sprintf("%s-dynamo-vpc-gateway-endpoint", projectName)
	dynamodbServiceName := fmt.Sprintf("com.amazonaws.%s.dynamodb", mainRegion)

	_, createDynamoVpcEndpointErr = ec2.NewVpcEndpoint(ctx, vpcEndpointName, &ec2.VpcEndpointArgs{
		VpcId:       pulumi.StringInput(vpcResource.ID()),
		ServiceName: pulumi.String(dynamodbServiceName),
		RouteTableIds: pulumi.StringArray{
			pulumi.StringInput(listOfPrivateRouteTables[0].ID()),
			pulumi.StringInput(listOfPrivateRouteTables[1].ID()),
			pulumi.StringInput(listOfPrivateRouteTables[2].ID()),
		},
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(vpcEndpointName),
			"ManagedBy": pulumi.String("pulumi"),
		},
	})

	if createDynamoVpcEndpointErr != nil {
		return nil, createDynamoVpcEndpointErr
	}

	return vpcEndpointResource, nil
}
