package network

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateS3VpcEndpoint(ctx *pulumi.Context, projectName string, mainRegion string, vpcResource *ec2.Vpc) (vpcEndpointResource *ec2.VpcEndpoint, createS3VpcEndpointErr error) {
	vpcEndpointName := fmt.Sprintf("%s-s3-vpc-gateway-endpoint", projectName)
	s3ServiceName := fmt.Sprintf("com.amazonaws.%s.s3", mainRegion)
	_, createS3VpcEndpointErr = ec2.NewVpcEndpoint(ctx, vpcEndpointName, &ec2.VpcEndpointArgs{
		VpcId:       pulumi.StringInput(vpcResource.ID()),
		ServiceName: pulumi.String(s3ServiceName),
		// TODO: switch this to non hardcoded but dynamic values for route table id
		RouteTableIds: pulumi.StringArray{
			pulumi.String("rtb-0b63b6fa519b6c77e"),
			pulumi.String("rtb-038a88c6f1f0d4258"),
			pulumi.String("rtb-089010a8ca4c5b45e"),
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

func CreateDynamoDBVpcEndpoint(ctx *pulumi.Context, projectName string, mainRegion string, vpcResource *ec2.Vpc) (vpcEndpointResource *ec2.VpcEndpoint, createDynamoVpcEndpointErr error) {
	vpcEndpointName := fmt.Sprintf("%s-dynamo-vpc-gateway-endpoint", projectName)
	dynamodbServiceName := fmt.Sprintf("com.amazonaws.%s.dynamodb", mainRegion)
	_, createDynamoVpcEndpointErr = ec2.NewVpcEndpoint(ctx, vpcEndpointName, &ec2.VpcEndpointArgs{
		VpcId:       pulumi.StringInput(vpcResource.ID()),
		ServiceName: pulumi.String(dynamodbServiceName),
		// TODO: switch this to non hardcoded but dynamic values for route table id
		RouteTableIds: pulumi.StringArray{
			pulumi.String("rtb-0b63b6fa519b6c77e"),
			pulumi.String("rtb-038a88c6f1f0d4258"),
			pulumi.String("rtb-089010a8ca4c5b45e"),
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
