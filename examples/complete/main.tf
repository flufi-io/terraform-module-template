module "random" {
  source     = "../../"
  context    = module.this.context
  attributes = [var.secret]
}
