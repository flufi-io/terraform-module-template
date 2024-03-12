package test

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	log "github.com/sirupsen/logrus"
)

var TimeToDestroy, _ = strconv.Atoi(os.Getenv("TIME_TO_DESTROY"))

func TestMain(m *testing.M) {
	// Decrypt secrets before running tests
	err := runSecretsScript("secrets.sh", "-d", "sandbox")
	if err != nil {
		log.Fatalf("Failed to decrypt secrets: %v", err)
	}
	// Run pre-commit due to the decrypted secrets file.
	err = runInitCommands()
	if err != nil {
		log.Fatalf("Failed to run commands: %v", err)
	}
	_ = runPreCommit()

	// Run tests
	exitVal := m.Run()

	// Encrypt secrets after tests
	err = runSecretsScript("secrets.sh", "-e", "sandbox")
	if err != nil {
		log.Fatalf("Failed to encrypt secrets: %v", err)
	}

	// Exit with the exit code determined by the tests
	os.Exit(exitVal)
}

func TestTerraform(t *testing.T) {
	t.Parallel()
	// Generate a random string
	randHash := random.UniqueId()
	originalName := terraform.GetVariableAsStringFromVarFile(t, "../../examples/complete/terraform.tfvars", "name")
	// Update the name variable with the original value plus the hash
	name := originalName + "-tst-" + randHash
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../examples/complete/",
		VarFiles:     []string{"fixtures.sandbox.us-east-1.tfvars"},
		Vars: map[string]interface{}{
			"name": name,
		},
		Upgrade:              true, // already done in init commands
		Reconfigure:          true, // already done in init commands
		Lock:                 true,
		SetVarsAfterVarFiles: true,
		BackendConfig: map[string]interface{}{
			"bucket":         "flufi-terraform-backend-sandbox-us-east-1",
			"key":            originalName + ".tfstate",
			"region":         "us-east-1",
			"dynamodb_table": "flufi-terraform-backend-sandbox-us-east-1",
			"encrypt":        true,
			"acl":            "bucket-owner-full-control",
		},
	})

	defer terraform.Destroy(t, terraformOptions)
	defer func() {
		timer(TimeToDestroy)
	}()

	terraform.InitAndApply(t, terraformOptions)
}
func TestPreCommit(t *testing.T) {
	err := runPreCommit()
	if err != nil {
		t.Fatalf("Failed to run commands: %s", err)
	}
}

func runSecretsScript(scriptName string, args ...string) error {
	// Prepend the script name to the args slice
	log.Printf("Running script: %s", scriptName)
	commandArgs := append([]string{scriptName}, args...)

	cmd := exec.Command("bash", commandArgs...)
	cmd.Dir = "../../examples/complete/"
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to execute script: %s", err)
		return err
	}
	return nil
}

func runInitCommands() error {
	// Run 'pre-commit install --install-hooks'
	err := runCommand("pre-commit", "install", "--install-hooks")
	if err != nil {
		log.Printf("Failed to run pre-commit install --install-hooks: %s", err)
	}

	// Run 'pre-commit autoupdate'
	err = runCommand("pre-commit", "autoupdate")
	if err != nil {
		log.Printf("Failed to run pre-commit autoupdate: %s", err)
	}

	// Run 'pre-commit run -a'
	err = runCommand("pre-commit", "run", "-a")
	if err != nil {
		log.Printf("Failed to run pre-commit run -a: %s", err)
	}
	return nil
}

func runPreCommit() error {
	// Run 'pre-commit run -a'
	if err := runCommand("pre-commit", "run", "-a"); err != nil {
		return err
	}
	return nil
}

func runCommand(command string, args ...string) error {
	fmt.Printf("Running command: %s %v\n", command, args)
	cmd := exec.Command(command, args...)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		fmt.Printf("Command output: %s\n", string(cmdOutput))
		return err
	}
	fmt.Printf("Command output: %s\n", string(cmdOutput))
	return nil
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
