package test

import (
	"github.com/gruntwork-io/terratest/modules/random"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCompleteExample(t *testing.T) {
	// retryable errors in terraform testing.
	uniqueID := random.UniqueId()
	repoName := "template-repository-" + uniqueID
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../examples/complete",
		Vars: map[string]interface{}{
			"name":        repoName,
			"description": "testing template repository",
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "repository_name")
	assert.Equal(t, repoName, output)
}
