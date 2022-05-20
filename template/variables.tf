variable "name" {
  type        = string
  description = "Name of the repository"
}
variable "description" {
  type        = string
  description = "Description of the repository"
}

variable "visibility" {
  type        = string
  description = "Visibility of the repository"
  default     = "public"
}