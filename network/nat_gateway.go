package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateNatGateway(ctx *pulumi.Context, vpcId pulumi.StringInput, projectName string, subnetId string) (natGwConfigObject *ec2.NatGateway, createNatGwErr error) {
	natGwName := projectName + "-igw"

	natGwConfig, createNatGwErr := ec2.NewNatGateway(ctx, natGwName, &ec2.NatGatewayArgs{
		ConnectivityType: pulumi.String("public"),
		AllocationId:     pulumi.Any(aws_eip.Example.Id),
		SubnetId:         pulumi.Any(aws_subnet.Example.Id),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(natGwName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{
		aws_internet_gateway.Example,
	}))
	if createNatGwErr != nil {
		return nil, createNatGwErr
	}
	return natGwConfig, nil
}
