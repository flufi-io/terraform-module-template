locals {
  template_dirs = {
    config = {
      path = "${path.module}/template/.config"
    }
    github = {
      path = "${path.module}/template/.github"
    }
    workflows = {
      path = "${local.template_dirs.github.path}/workflows"
    }
  }
#  template_files = {
#    terraform-docs-config = {
#      path = "${local.template_dirs.config.path}/.terraform-docs.yml"
#    }
#    atlantis = {
#      path = "${local.template_dirs.workflows.path}/atlantis.yml"
#    }
#    terraform-docs = {
#      path = "${local.template_dirs.workflows.path}/terraform-docs.yml"
#    }
#    terratest = {
#      path = "${local.template_dirs.workflows.path}/terratest.yml"
#    }
#    tfsec = {
#      path = "${local.template_dirs.workflows.path}/tfsec.yml"
#    }
#  }
  template_files = [
    "${path.module}/template/.config.terraform-docs.yml",
    "${path.module}/template/.github/atlantis.yml",
    "${path.module}/template/.github/terraform-docs.yml",
    "${path.module}/template/.github/terratest.yml",
    "${path.module}/template/.github/tfsec.yml",
    "${path.module}/template/.github/atlantis.yml",
  ]
}
resource "github_repository" "template" {
  name        = var.name
  description = var.description
  visibility  = var.visibility
  auto_init   = true
}

resource "github_branch" "main" {
  repository = github_repository.template.name
  branch     = "main"
}

resource "github_branch_default" "default"{
  repository = github_repository.template.name
  branch     = github_branch.main.branch
}

resource "github_repository_file" "template_files" {
  for_each            = tolist(toset(local.template_files))
  repository          = github_repository.template.name
  branch              = github_branch_default.default.branch
  file                = each.value["path"]
  content             = file(each.value["path"])
  commit_message      = "Managed by Terraform"
  overwrite_on_create = true
}

