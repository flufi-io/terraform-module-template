module "random" {
  source           = "../../"
  length           = var.length
  context          = module.this.context
  override_special = var.override_special
  special          = var.special
}
