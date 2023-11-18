package network

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateInternetGateway(ctx *pulumi.Context, vpcId pulumi.StringInput, projectName string, indexNum string) (igwResourceObject *ec2.InternetGateway, createIgwErr error) {
	// TODO: make this take an ID or count or something to not conflict when more than 1 nat has to be created
	igwName := fmt.Sprintf("%s-igw-%s", projectName, indexNum)

	igwResource, createIgwErr := ec2.NewInternetGateway(ctx, igwName, &ec2.InternetGatewayArgs{
		VpcId: pulumi.StringInput(vpcId),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(igwName),
			"ManagedBy": pulumi.String("pulumi"),
		},
	})
	if createIgwErr != nil {
		return nil, createIgwErr
	}
	return igwResource, nil
}
