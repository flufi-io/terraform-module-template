config {
  format = "compact"
  plugin_dir = "~/.tflint.d/plugins"

  module = true
  force = true
  disabled_by_default = false

}

plugin "aws" {
  enabled = true
  version = "0.30.0"
  source  = "github.com/terraform-linters/tflint-ruleset-aws"
}
rule "aws_instance_invalid_type" {
  enabled = true
}

rule "terraform_unused_declarations" {
  enabled = true
}
