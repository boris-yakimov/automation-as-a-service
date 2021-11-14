package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateVPC(ctx *pulumi.Context, vpcCidrRange string) error {
	_, createVpcErr := ec2.NewVpc(ctx, "main", &ec2.VpcArgs{
		CidrBlock: pulumi.String(vpcCidrRange),
	})
	if createVpcErr != nil {
		return createVpcErr
	}
	return nil
}
