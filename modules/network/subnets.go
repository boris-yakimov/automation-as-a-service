package network

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateSubnet(ctx *pulumi.Context, vpcId pulumi.StringInput, projectName string, subnetType string, subnetName string, subnetRange string) (subnetResourceOjbect *ec2.Subnet, createSubnetErr error) {
	if subnetType != "public" && subnetType != "private" {
		return nil, fmt.Errorf("Incorrect subnet type, supported types are \"public\" and \"private\"")
	}

	fullSubnetName := fmt.Sprintf("%s-%s", projectName, subnetName)

	//TODO: add check if subnetRange not in CIDR format

	subnetResource, createSubnetErr := ec2.NewSubnet(ctx, fullSubnetName, &ec2.SubnetArgs{
		VpcId:     pulumi.StringInput(vpcId),
		CidrBlock: pulumi.String(subnetRange),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(subnetName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	})
	if createSubnetErr != nil {
		return nil, createSubnetErr
	}
	return subnetResource, nil
}
