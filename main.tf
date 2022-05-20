locals {
  template_files = [
    "../../template/.config/.terraform-docs.yml",

    "../../template/.github/workflows/atlantis.yml",
    "../../template/.github/workflows/terraform-docs.yml",
    "../../template/.github/workflows/terratest.yml",
    "../../template/.github/workflows/tfsec.yml",

    "../../template/.gitignore",
    "../../template/main.tf",
    "../../template/variables.tf",
    "../../template/versions.tf",
    "../../template/outputs.tf",
    "../../template/README.md",

    "../../template/examples/complete/main.tf",
    "../../template/examples/complete/variables.tf",
    "../../template/examples/complete/versions.tf",
    "../../template/examples/complete/outputs.tf",
    "../../template/examples/complete/terraform.tfvars",
    "../../template/examples/complete/providers.tf",

    "../../template/tests/complete/complete_test.go",
  ]
}
resource "github_repository" "template" {
  name        = var.name
  description = var.description
  visibility  = "public"
  auto_init   = true
}

#resource "github_branch" "main" {
#  repository = github_repository.template.name
#  branch     = "main"
#}

#resource "github_branch_default" "default"{
#  repository = github_repository.template.name
#  branch     = "main"
#}
#resource "github_branch_protection_v3" "main" {
#  repository     = github_repository.template.name
#  branch         = github_branch_default.default.branch
#  enforce_admins = true
#
#  restrictions {
#    users = ["pipo-flufi"]
#  }
#}


resource "github_repository_file" "template_files" {
  depends_on = [github_repository.template]
  for_each            = toset(local.template_files)
  repository          = github_repository.template.name
  branch              = "main"
  file                = replace(each.value, "../../template/", "")
  content             = file(each.value)
  commit_message      = "Managed by Terraform"
  overwrite_on_create = true
  commit_author       = "Terraform User"
  commit_email        = "terraform@example.com"
}
