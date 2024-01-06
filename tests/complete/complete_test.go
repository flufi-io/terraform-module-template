package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

var TimeToDestroy, _ = strconv.Atoi(os.Getenv("TIME_TO_DESTROY"))

func Test(t *testing.T) {
	t.Parallel()
	// Run secrets.sh script
	runSecretsScript(t, "secrets.sh", "-d", "sandbox")
	// Generate a random string
	randHash := random.UniqueId()
	originalName := terraform.GetVariableAsStringFromVarFile(t, "../../examples/complete/fixtures.sandbox.tfvars.json", "name")
	// Update the name variable with the original value plus the hash
	name := originalName + "-terratest-" + randHash
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../examples/complete/",
		VarFiles:     []string{"fixtures.sandbox.tfvars.json"},
		Vars: map[string]interface{}{
			"name": name,
		},
		Upgrade:              true,
		Reconfigure:          true,
		Lock:                 true,
		SetVarsAfterVarFiles: true,
	})

	defer runSecretsScript(t, "secrets.sh", "-e", "sandbox")
	defer terraform.Destroy(t, terraformOptions)
	defer func() {
		timer(TimeToDestroy)
	}()

	terraform.InitAndApply(t, terraformOptions)
}

func runSecretsScript(t *testing.T, scriptName string, args ...string) {
	// Prepend the script name to the args slice
	t.Log("Running script: " + scriptName)
	commandArgs := append([]string{scriptName}, args...)

	cmd := exec.Command("bash", commandArgs...)
	cmd.Dir = "../../examples/complete/"
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to execute script: %s", err)
	}
}

func timer(s int) {
	for {
		if s <= 0 {
			break
		} else {
			fmt.Println(s)
			time.Sleep(1 * time.Second) // wait 1 sec
			s--                         // reduce time
		}
	}
}
