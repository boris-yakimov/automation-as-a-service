package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateInternetGateway(ctx *pulumi.Context, vpcId pulumi.StringInput, projectName string) (igwResourceObject *ec2.InternetGateway, createIgwErr error) {
	igwName := projectName + "-igw"

	igwResource, createIgwErr := ec2.NewInternetGateway(ctx, igwName, &ec2.InternetGatewayArgs{
		VpcId: pulumi.StringInput(vpcId),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(igwName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	})
	if createIgwErr != nil {
		return nil, createIgwErr
	}
	return igwResource, nil
}
