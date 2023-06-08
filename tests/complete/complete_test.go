package test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"os"
	"strconv"
	"testing"
	"time"
)

var TimeToDestroy, _ = strconv.Atoi(os.Getenv("TIME_TO_DESTROY"))

func TestWorkspace(t *testing.T) {
	t.Parallel()

	// Generate a random string
	randHash, _ := generateRandomString(10)
	tempTfvarsFile := fmt.Sprintf("terraform_%s.tfvars", randHash)

	// Prepare the content for the tfvars file
	content := []byte(fmt.Sprintf("name = %s", randHash))

	// Write the content to the tfvars file
	err := os.WriteFile(tempTfvarsFile, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	originalName := terraform.GetVariableAsStringFromVarFile(t, tempTfvarsFile, "name")

	// Update the name variable with the original value plus the hash
	name := originalName + randHash

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/complete/",
		VarFiles:     []string{tempTfvarsFile},
		Vars: map[string]interface{}{
			"name": name,
		},
		Upgrade:     true,
		Reconfigure: true,
	}

	defer terraform.Destroy(t, terraformOptions)
	defer func() {
		timer(TimeToDestroy)
	}()

	defer func() {
		err := os.Remove(tempTfvarsFile) // Clean up the temporary tfvars file
		if err != nil {
			t.Logf("Failed to remove temp file %s: %v", tempTfvarsFile, err)
		}
	}()

	terraform.InitAndApply(t, terraformOptions)
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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
