package provisioning

import (
	"automation-as-a-service/modules/ecr"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Ecr(ctx *pulumi.Context, projectName string, listOfEcrRepos map[string]string) (createListOfEcrReposErr error) {
	_, createListOfEcrReposErr = ecr.CreateECR(ctx, projectName, listOfEcrRepos)
	if createListOfEcrReposErr != nil {
		return createListOfEcrReposErr
	}

	// TODO: figure out where we can get those from - maybe an env var or lists with config options (or config file with config options) - for now it is hardcoded as true
	var enableEcrLifecyclePolicy bool = true
	var imageRetentionPeriodInDays = "90"
	var ecrLifecyclePolicyName = "image-cleanup-policy"
	if enableEcrLifecyclePolicy {
		configureEcrLifecycleErr := ecr.ConfigureEcrLifecyclePolicy(ctx, ecrLifecyclePolicyName, imageRetentionPeriodInDays)
		if configureEcrLifecycleErr != nil {
			return configureEcrLifecycleErr
		}
	}

	return nil
}
