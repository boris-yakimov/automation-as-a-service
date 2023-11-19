package network

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// TODO: check why VPC id is not used here
func CreateNatGateway(ctx *pulumi.Context, vpcId pulumi.StringInput, projectName string, indexNum string, subnetId pulumi.StringInput, vpcResource *ec2.Vpc) (natGwResourceObject *ec2.NatGateway, createNatGwErr error) {
	// TODO: add validations to make sure those are not empty
	natGwName := fmt.Sprintf("%s-natgw-%s", projectName, indexNum)

	eipResource, createEipErr := CreateEIP(ctx, projectName, "natgw", indexNum, vpcResource)
	if createEipErr != nil {
		return nil, createEipErr
	}

	natGwResource, createNatGwErr := ec2.NewNatGateway(ctx, natGwName, &ec2.NatGatewayArgs{
		ConnectivityType: pulumi.String("public"),
		AllocationId:     pulumi.StringInput(eipResource.ID()),
		SubnetId:         pulumi.StringInput(subnetId),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(natGwName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		eipResource,
	}), pulumi.Parent(eipResource),
	)
	if createNatGwErr != nil {
		return nil, createNatGwErr
	}
	return natGwResource, nil
}
