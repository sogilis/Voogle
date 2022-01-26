resource "aws_kms_key" "voogle-video-encryption-key" {
  description             = "This key is used to encrypt bucket that stores voogle videos"
  deletion_window_in_days = 10
}

resource "aws_s3_bucket" "voogle-video-s3-bucket" {
  bucket = "voogle-video"
  acl    = "private"

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
