package ecr

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateECR(ctx *pulumi.Context, projectName string, listOfEcrRepos map[string]string) (ecrResourceObject *ecr.Repository, createEcrError error) {
	var ecrResource *ecr.Repository

	// TODO: add optional flag if ImageScanning should be enabled - low prio

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

func ConfigureEcrLifecyclePolicy(ctx *pulumi.Context, ecrLifecyclePolicyName string, imageRetentionPeriod string) (configureEcrLifecyclePolicyErr error) {
	lifeCyclePolicyJson := fmt.Sprintf(`{
   "rules": [
				{
					"rulePriority": 1,
					"description": "Expire images older than 14 days",
					"selection": {
						"tagStatus": "untagged",
						"countType": "sinceImagePushed",
						"countUnit": "days",
						"countNumber": 14
					},
					"action": {
						"type": "expire"
					}
				},
				{
					"rulePriority": 2,
					"description": "Expire all images older than <imateRetentionPeriod> days",
					"selection": {
						"tagStatus": "any",
						"countType": "sinceImagePushed",
						"countUnit": "days",
						"countNumber": %s
					},
					"action": {
						"type": "expire"
					}
				},
				{
					"rulePriority": 3,
					"description": "Keep at most 100 images, expire all others",
					"selection": {
						"tagStatus": "any",
						"countType": "imageCountMoreThan",
						"countNumber": 100
					},
					"action": {
						"type": "expire"
					}
				}
			]
		}`, imageRetentionPeriod)

	_, configureEcrLifecyclePolicyErr = ecr.NewLifecyclePolicy(ctx, ecrLifecyclePolicyName, &ecr.LifecyclePolicyArgs{
		Repository: pulumi.String(ecrLifecyclePolicyName),
		Policy:     pulumi.Any(lifeCyclePolicyJson),
	})
	if configureEcrLifecyclePolicyErr != nil {
		return configureEcrLifecyclePolicyErr
	}

	return nil
}
