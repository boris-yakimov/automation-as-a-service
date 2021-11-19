package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateInternetGateway(ctx *pulumi.Context, vpcId pulumi.StringInput, projectName string) (igwConfigObject *ec2.InternetGateway, createIgwErr error) {
	igwName := projectName + "-igw"

	igwConfig, createIgwErr := ec2.NewInternetGateway(ctx, igwName, &ec2.InternetGatewayArgs{
		VpcId: pulumi.StringInput(vpcId),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(igwName),
			"ManagedBy": pulumi.String("Pulumi"),
		},
	})
	if createIgwErr != nil {
		return nil, createIgwErr
	}
	return igwConfig, nil
}
