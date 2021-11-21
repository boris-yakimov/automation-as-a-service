package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateEIP(ctx *pulumi.Context, projectName string, eipPurpose string) (eipResourceObject *ec2.Eip, createEipErr error) {
	// TODO: add validations to make sure those are not empty
	eipName := projectName + "-eip-" + eipPurpose

	// TODO: make it depend on the VPC
	// TODO: check issues with err         * Error creating EIP: InvalidPublicIpv4Pool.NotFound: The pool ID 'ipv4pool-ec2-012345' does not exist. status code: 400, request id: 98edaf82-d7a7-448f-be6c-a24ab834d477
	eipResource, createEipErr := ec2.NewEip(ctx, eipName, &ec2.EipArgs{
		PublicIpv4Pool: pulumi.String("ipv4pool-ec2-012345"),
		Vpc:            pulumi.Bool(true),
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
