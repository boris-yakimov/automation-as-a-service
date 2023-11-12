package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateEIP(ctx *pulumi.Context, projectName string, eipPurpose string) (eipResourceObject *ec2.Eip, createEipErr error) {
	// TODO: add validations to make sure those are not empty
	eipName := projectName + "-eip-" + eipPurpose

	eipResource, createEipErr := ec2.NewEip(ctx, eipName, &ec2.EipArgs{
		Domain:         pulumi.String("vpc"),
		PublicIpv4Pool: pulumi.String("amazon"),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(eipName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	})
	if createEipErr != nil {
		return nil, createEipErr
	}
	return eipResource, nil
}
