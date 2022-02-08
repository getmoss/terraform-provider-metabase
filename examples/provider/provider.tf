variable "host" {}
variable "username" {}
variable "password" {}

provider "metabase" {
  username = var.username
  password = var.password
  host     = var.host
}