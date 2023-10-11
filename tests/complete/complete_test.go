package test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"os"
	"strconv"
	"testing"
	"time"
)

var TimeToDestroy, _ = strconv.Atoi(os.Getenv("TIME_TO_DESTROY"))

func Test(t *testing.T) {
	t.Parallel()

	// Generate a random string
	randHash := uuid.New().String()
	originalName := terraform.GetVariableAsStringFromVarFile(t, "../../examples/complete/terraform.tfvars", "name")

	// Update the name variable with the original value plus the hash
	name := originalName + "-terratest-" + randHash

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/complete/",
		VarFiles:     []string{"terraform.tfvars"},
		Vars: map[string]interface{}{
			"name": name,
		},
		Upgrade:     true,
		Reconfigure: true,
		Lock:        true,
	}

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
