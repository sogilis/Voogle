resource "aws_kms_key" "voogle-video-encryption-key" {
  description                        = "This key is used to encrypt bucket that stores voogle videos"
  deletion_window_in_days            = 30
  bypass_policy_lockout_safety_check = false
}

resource "aws_s3_bucket" "voogle-video-s3-bucket" {
  bucket        = "voogle-video"
  acl           = "private"
  force_destroy = false

  versioning {
    enabled = false
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.voogle-video-encryption-key.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }
}

resource "aws_s3_bucket_public_access_block" "voogle-video-s3-bucket-public-access-block" {
  bucket                  = aws_s3_bucket.voogle-video-s3-bucket.id
  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}
