package main

import (
	"automation-as-a-service/provisioning"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// TODO: configure a list of env vars or a config file that passes a list of things that should be created and if enabled add required variables for them
		// TODO: set limit on how long a project name can be to not exhaust char limits when creating resources with longer names
		// TODO: add unit tests
		projectName := "temp-test"
		vpcCidrRange := "10.0.0.0/16"
		subnetList := map[string]string{
			"private-subnet1": "10.0.0.0/20", // 4k IPs per subnet
			"private-subnet2": "10.0.32.0/20",
			"private-subnet3": "10.0.64.0/20",

			"public-subnet1": "10.0.128.0/20",
			"public-subnet2": "10.0.160.0/20",
			"public-subnet3": "10.0.192.0/20",
		}

		listOfEcrRepos := map[string]string{
			"test-ecr-repo1": "test-app-docker",
			"test-ecr-repo2": "test-app-helm",
			"test-ecr-repo3": "test-app-base-image",
		}

		networkProvisioningErr := provisioning.Network(ctx, projectName, vpcCidrRange, subnetList)
		if networkProvisioningErr != nil {
			return networkProvisioningErr
		}

		ecrProvisioningErr := provisioning.Ecr(ctx, projectName, listOfEcrRepos)
		if ecrProvisioningErr != nil {
			return ecrProvisioningErr
		}

		return nil
	})
}
