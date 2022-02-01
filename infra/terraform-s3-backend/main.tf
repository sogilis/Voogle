terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = "default"
  region  = "eu-west-3"
}

resource "aws_kms_key" "tfstate-encryption-key-voogle-sogilis-dev" {
  description             = "This key is used to encrypt bucket that stores the terraform state voogle sogilis dev"
  deletion_window_in_days = 10
}

resource "aws_dynamodb_table" "dynamodb-tfstate-lock-voogle-sogilis-dev" {
  name           = "tfstate-lock-dynamo-voogle-sogilis-dev"
  hash_key       = "LockID"
  read_capacity  = 20
  write_capacity = 20

  attribute {
    name = "LockID"
    type = "S"
  }
}

resource "aws_s3_bucket" "tf-state-bucket-voogle-sogilis-dev" {
  bucket = "terraform-state-voogle-sogilis-dev"
  acl    = "private"

  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.tfstate-encryption-key-voogle-sogilis-dev.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }
}
