package test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestCompleteExample(t *testing.T) {
	godotenv.Load("../../.env")

	// DELAY is the time in seconds to run terraform destroy after terraform apply
	DELAY, _ := strconv.Atoi(os.Getenv("DELAY"))
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../examples/complete",
		Vars:         map[string]interface{}{},
		Upgrade:      true,
		Reconfigure:  true,
	})

	defer terraform.Destroy(t, terraformOptions)
	// Delay the execution of the terraform destroy
	defer func() {
		delay(DELAY)
	}()

	terraform.InitAndApply(t, terraformOptions)

}

func delay(seconds int) {
	for {
		if seconds <= 0 {
			break
		} else {
			log.Println(seconds)
			time.Sleep(1 * time.Second)
			seconds--
		}
	}
}
