package network

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateSubnet(ctx *pulumi.Context, projectName string, subnetType string, subnetName string, subnetRange string, vpcResource *ec2.Vpc) (subnetResourceOjbect *ec2.Subnet, createSubnetErr error) {
	if subnetType != "public" && subnetType != "private" {
		return nil, fmt.Errorf("Incorrect subnet type, supported types are \"public\" and \"private\"")
	}

	fullSubnetName := fmt.Sprintf("%s-%s", projectName, subnetName)

	//TODO: add check if subnetRange not in CIDR format

	subnetResource, createSubnetErr := ec2.NewSubnet(ctx, fullSubnetName, &ec2.SubnetArgs{
		VpcId:     pulumi.StringInput(vpcResource.ID()),
		CidrBlock: pulumi.String(subnetRange),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(subnetName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	}, pulumi.Parent(vpcResource),
	)
	if createSubnetErr != nil {
		return nil, createSubnetErr
	}
	return subnetResource, nil
}
