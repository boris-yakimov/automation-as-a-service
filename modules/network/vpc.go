package network

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateVPC(ctx *pulumi.Context, projectName string, vpcCidrRange string) (vpcResourceObject *ec2.Vpc, createVpcErr error) {
	vpcName := fmt.Sprintf("%s-vpc", projectName)

	vpcResource, createVpcErr := ec2.NewVpc(ctx, vpcName, &ec2.VpcArgs{
		CidrBlock:          pulumi.String(vpcCidrRange),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport:   pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(vpcName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	})
	if createVpcErr != nil {
		return nil, createVpcErr
	}
	return vpcResource, nil
}
