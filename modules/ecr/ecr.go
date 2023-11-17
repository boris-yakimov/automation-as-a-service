package ecr

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateECR(ctx *pulumi.Context, projectName string, listOfEcrRepos map[string]string) (ecrResourceObject *ecr.Repository, createEcrError error) {
	var ecrResource *ecr.Repository

	for _, ecrRepoName := range listOfEcrRepos {
		// TODO: ecrResource should be in a map of all repos that is returned at the end
		ecrResource, createEcrError = ecr.NewRepository(ctx, ecrRepoName, &ecr.RepositoryArgs{
			// TODO: change to name from map
			//Name: pulumi.String(ecrRepo),
			ImageScanningConfiguration: &ecr.RepositoryImageScanningConfigurationArgs{
				ScanOnPush: pulumi.Bool(true),
			},
			ImageTagMutability: pulumi.String("MUTABLE"),
			Tags: pulumi.StringMap{
				"Name":      pulumi.String("testEcrRepo"),
				"ManagedBy": pulumi.String("pulumi"),
				"Project":   pulumi.String(projectName),
			},
		})
		if createEcrError != nil {
			return nil, createEcrError
		}
	}

	// TODO: this has to be a map of objects, even if it contains a single one
	return ecrResource, nil
}
