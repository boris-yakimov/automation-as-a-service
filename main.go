package main

import (
	//"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"automation-as-a-service/network"
)

func init() {
		vpcCidrRange := "10.0.0.0/16"
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		createVpcErr := network.CreateVPC(ctx, vpcCidrRange)
		if createVpcErr != nil {
			return createVpcErr
		}

		// Create an AWS resource (S3 Bucket)
		//bucket, err := s3.NewBucket(ctx, "boris-test", nil)
		//if err != nil {
			//return err
		//}

		//// Export the name of the bucket
		//ctx.Export("bucketName", bucket.ID())
		return nil
	})
}
