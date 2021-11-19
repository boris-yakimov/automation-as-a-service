package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"automation-as-a-service/network"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		projectName := "eduspire"
		vpcCidrRange := "10.0.0.0/16"

		// Create AWS VPC
		vpcConfig, createVpcErr := network.CreateVPC(ctx, projectName, vpcCidrRange)
		if createVpcErr != nil {
			return createVpcErr
		}

		// TODO: Not sure if these exports should not be moved on module level
		ctx.Export("vpcArn", vpcConfig.Arn)
		ctx.Export("vpcId", vpcConfig.ID())
		//vpcId := pulumi.String(vpcConfig.ID)
		//fmt.Println(vpcConfig.ID())

		// TODO: figure out output params, callbacks, etc
		fmt.Println(vpcConfig.ID().ToStringOutput())

		//igwConfig, createIgwErr := network.CreateInternetGateway(ctx, projectName, )

		return nil
	})
}
