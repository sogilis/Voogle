terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  backend "s3" {
    bucket         = "terraform-state-voogle-sogilis-dev"
    dynamodb_table = "tfstate-lock-dynamo"
    key            = "voogle/terraform.tfstate"
    region         = "eu-west-3"
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = "default"
  region  = "eu-west-3"
}

resource "aws_organizations_organization" "sogilis" {
  # (resource arguments)
}

module "voogle-sogilis-dev" {
  root_organization = aws_organizations_organization.sogilis.roots[0].id
  source            = "./voogle-sogilis-dev"
}

output "account-id" {
  value = module.voogle-sogilis-dev.account-id
}

output "token-id" {
  value = module.voogle-sogilis-dev.token-id
}

output "secret" {
  value     = module.voogle-sogilis-dev.secret
  sensitive = true
}

output "voogle-s3-token-id" {
  value = module.voogle-sogilis-dev.voogle-s3-token-id
}

output "voogle-s3-secret" {
  value     = module.voogle-sogilis-dev.voogle-s3-secret
  sensitive = true
}
