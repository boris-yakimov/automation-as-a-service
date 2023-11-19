package provisioning

import (
	"automation-as-a-service/modules/ecr"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Ecr(ctx *pulumi.Context, projectName string, listOfEcrRepos map[string]string) (createListOfEcrReposErr error) {
	// TODO: ecrResource should be in a map of all repos that is returned at the end
	// Create ECR repos
	for _, ecrRepoName := range listOfEcrRepos {
		ecrRepoResource, createListOfEcrReposErr := ecr.CreateECR(ctx, projectName, ecrRepoName)
		if createListOfEcrReposErr != nil {
			return createListOfEcrReposErr
		}

		// TODO: figure out where we can get those from - maybe an env var or lists with config options (or config file with config options) - for now it is hardcoded as true
		var enableEcrLifecyclePolicy bool = true
		var imageRetentionPeriodInDays = "90"
		var ecrLifecyclePolicyName = fmt.Sprintf("%s-lifecycle-policy", ecrRepoName)
		if enableEcrLifecyclePolicy {
			// Attach lifecycle policy to each ECR repo
			configureEcrLifecycleErr := ecr.ConfigureEcrLifecyclePolicy(ctx, ecrRepoName, ecrLifecyclePolicyName, imageRetentionPeriodInDays, ecrRepoResource)
			if configureEcrLifecycleErr != nil {
				return configureEcrLifecycleErr
			}
		}
	}

	return nil
}
