package ecr

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateECR(ctx *pulumi.Context, projectName string, ecrRepoName string) (ecrResourceObject *ecr.Repository, createEcrError error) {
	var ecrResource *ecr.Repository

	// TODO: add optional flag if ImageScanning should be enabled - low prio
	ecrResource, createEcrError = ecr.NewRepository(ctx, ecrRepoName, &ecr.RepositoryArgs{
		// passing name as input is not required but otherwise generates random strings at the end of name causing issues attaching policy after since policy supports only attachment by name (which is also not an output param of create ECR repo, therefore cannot be taken dynamically by policy creation...(this is annoying in pulumi)
		Name: pulumi.String(ecrRepoName),
		ImageScanningConfiguration: &ecr.RepositoryImageScanningConfigurationArgs{
			ScanOnPush: pulumi.Bool(true),
		},
		ImageTagMutability: pulumi.String("MUTABLE"),
		Tags: pulumi.StringMap{
			"Name":      pulumi.String(ecrRepoName),
			"ManagedBy": pulumi.String("pulumi"),
			"Project":   pulumi.String(projectName),
		},
	})
	if createEcrError != nil {
		return nil, createEcrError
	}

	return ecrResource, nil
}

func ConfigureEcrLifecyclePolicy(ctx *pulumi.Context, ecrRepoName string, ecrLifecyclePolicyName string, imageRetentionPeriod string, ecrRepoResource *ecr.Repository) (configureEcrLifecyclePolicyErr error) {
	_, configureEcrLifecyclePolicyErr = ecr.NewLifecyclePolicy(ctx, ecrLifecyclePolicyName, &ecr.LifecyclePolicyArgs{
		Repository: pulumi.String(ecrRepoName),
		Policy: pulumi.Any(`{
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
						"countType": "imageCountMoreThan",
						"countNumber": 50 
					},
					"action": {
						"type": "expire"
					}
				}
			]
		}`)},
		pulumi.Parent(ecrRepoResource),
		pulumi.DependsOn([]pulumi.Resource{ecrRepoResource}),
	)
	if configureEcrLifecyclePolicyErr != nil {
		return configureEcrLifecyclePolicyErr
	}

	return nil
}
