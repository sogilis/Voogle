resource "aws_iam_group" "voogle-s3-group" {
  name = "voogle-s3-group"
  path = "/voogle-sogilis-dev/"
}

resource "aws_iam_group_policy" "voogle-s3-group-policy" {
  name  = "voogle-s3-group-policy"
  group = aws_iam_group.voogle-s3-group.name

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect":"Allow",
        "Action":[
          "s3:CreateBucket",
          "s3:ListAllMyBuckets",
        ],
        "Resource":[
          "arn:aws:s3:::*"
        ]
      },
      {
        "Effect" : "Allow",
        "Action" : [
          "s3:ListBucket",
        ],
        "Resource" : [
          "arn:aws:s3:::${aws_s3_bucket.voogle-video-s3-bucket.bucket}",
        ]
      },
      {
        "Effect" : "Allow",
        "Action" : [
          "s3:GetObject",
          "s3:PutObject",
        ],
        "Resource" : [
          "arn:aws:s3:::${aws_s3_bucket.voogle-video-s3-bucket.bucket}/*"
        ]
      },
      {
        "Effect" : "Allow",
        "Action" : [
          "kms:Decrypt",
          "kms:GenerateDataKey",
        ],
        "Resource" : [
          "arn:aws:kms:eu-west-3:341514320952:key/f889ca9c-80cc-4cf5-b593-736071618385"
        ]
      },
    ]
  })
}
