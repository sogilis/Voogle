terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  backend "s3" {
    bucket         = "terraform-state-voogle-sogilis-dev"
    dynamodb_table = "tfstate-lock-dynamo-voogle-sogilis-dev"
    key            = "voogle/terraform.tfstate"
    region         = "eu-west-3"
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = "default"
  region  = "eu-west-3"
}

module "voogle-sogilis-dev" {
  source = "./voogle-sogilis-dev"
}

output "voogle-s3-token-id" {
  value = module.voogle-sogilis-dev.voogle-s3-token-id
}

output "voogle-s3-secret" {
  value     = module.voogle-sogilis-dev.voogle-s3-secret
  sensitive = true
}
