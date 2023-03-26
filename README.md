# terraform-module-template

<!-- BEGIN_TF_DOCS -->
# Examples
```hcl
module "random" {
  source = "../../"
  length = var.length
}

```

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.2.0 |
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_length"></a> [length](#input\_length) | Length of the password | `number` | `14` | no |
## Outputs

| Name | Description |
|------|-------------|
| <a name="output_random_password"></a> [random\_password](#output\_random\_password) | this is a dummy password |
## Resources

- resource.random_password.password (main.tf#1)
<!-- END_TF_DOCS -->
