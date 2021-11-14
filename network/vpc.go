package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)


//TODO: configure output params as a struct ? 
func CreateVPC(ctx *pulumi.Context, projectName string, vpcCidrRange string) (vpcConfigObject *ec2.Vpc, createVpcErr error ) {
	vpcName := projectName + "-vpc"

	vpcConfig, createVpcErr := ec2.NewVpc(ctx, vpcName, &ec2.VpcArgs{
		CidrBlock: pulumi.String(vpcCidrRange),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport: pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name": pulumi.String(vpcName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	})
	if createVpcErr != nil {
		return nil, createVpcErr
	}
	return vpcConfig, nil
}
