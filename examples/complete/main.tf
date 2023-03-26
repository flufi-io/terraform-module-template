module "random" {
  source  = "../../"
  length  = var.length
  context = module.this.context
}
