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

	return nil
}
