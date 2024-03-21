package test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

var TimeToDestroy, _ = strconv.Atoi(os.Getenv("TIME_TO_DESTROY"))

func TestTerraform(t *testing.T) {
	t.Parallel()
	// Generate a random string
	randHash := random.UniqueId()
	originalName := terraform.GetVariableAsStringFromVarFile(t, "../../examples/basic/terraform.tfvars", "name")
	// Update the name variable with the original value plus the hash
	name := originalName + "-tst-" + randHash
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../examples/basic/",
		VarFiles:     []string{"terraform.tfvars"},
		Vars: map[string]interface{}{
			"name": name,
		},
		Upgrade:              true,
		Reconfigure:          true,
		Lock:                 true,
		SetVarsAfterVarFiles: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	defer func() {
		timer(TimeToDestroy)
	}()

	terraform.InitAndApply(t, terraformOptions)
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
