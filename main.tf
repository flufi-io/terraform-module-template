resource "random_password" "password" {
  length           = var.length
  special          = var.special
  override_special = var.override_special
}
