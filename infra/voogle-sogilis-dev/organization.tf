variable "root_organization" {}

resource "aws_organizations_organizational_unit" "voogle-sogilis-dev" {
  name      = "voogle-sogilis-dev"
  parent_id = var.root_organization
}
