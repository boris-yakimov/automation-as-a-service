package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateNatGateway(ctx *pulumi.Context, vpcId pulumi.StringInput, projectName string, subnetId pulumi.StringInput, igwResource *ec2.InternetGateway) (natGwResourceObject *ec2.NatGateway, createNatGwErr error) {
	// TODO: make this take an ID or count or something to not conflict when more than 1 nat has to be created
	natGwName := projectName + "-natgw"

	eipResource, createEipErr := CreateEIP(ctx, projectName, "natgw")
	if createEipErr != nil {
		return nil, createEipErr
	}

	natGwResource, createNatGwErr := ec2.NewNatGateway(ctx, natGwName, &ec2.NatGatewayArgs{
		ConnectivityType: pulumi.String("public"),
		//AllocationId:     pulumi.Any(aws_eip.Example.Id),
		AllocationId: pulumi.StringInput(eipResource.ID()),
		//SubnetId:     pulumi.Any(aws_subnet.Example.Id),
		SubnetId: pulumi.StringInput(subnetId),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(natGwName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		igwResource,
	}))
	if createNatGwErr != nil {
		return nil, createNatGwErr
	}
	return natGwResource, nil
}
