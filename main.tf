locals {
  template_files = [
    "${path.module}/template/.config/.terraform-docs.yml",

    "${path.module}/template/.github/workflows/atlantis.yml",
    "${path.module}/template/.github/workflows/terraform-docs.yml",
    "${path.module}/template/.github/workflows/terratest.yml",
    "${path.module}/template/.github/workflows/tfsec.yml",

    "${path.module}/template/.gitignore",
    "${path.module}/template/main.tf",
    "${path.module}/template/variables.tf",
    "${path.module}/template/versions.tf",
    "${path.module}/template/outputs.tf",
    "${path.module}/template/README.md",

    "${path.module}/template/examples/complete/main.tf",
    "${path.module}/template/examples/complete/variables.tf",
    "${path.module}/template/examples/complete/versions.tf",
    "${path.module}/template/examples/complete/outputs.tf",
    "${path.module}/template/examples/complete/terraform.tfvars",
    "${path.module}/template/examples/complete/providers.tf",

    "${path.module}/template/tests/complete/complete_test.go",
  ]
}
resource "github_repository" "template" {
  name        = var.name
  description = var.description
  visibility  = var.visibility
  auto_init   = true

}

#resource "github_branch" "main" {
#  repository = github_repository.template.name
#  branch     = "main"
#}

resource "github_branch_default" "default"{
  repository = github_repository.template.name
  branch     = "main"
}

resource "github_repository_file" "template_files" {
  depends_on = [github_branch_default.default]
  for_each            = toset(local.template_files)
  repository          = github_repository.template.name
  branch              = github_branch_default.default.branch
  file                = each.value
  content             = file(each.value)
  commit_message      = "Managed by Terraform"
  overwrite_on_create = true
  commit_author       = "Terraform User"
  commit_email        = "terraform@example.com"
}

