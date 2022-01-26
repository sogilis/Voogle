resource "aws_iam_access_key" "voogle-sqsc" {
  user = aws_iam_user.sqsc.name
}

resource "aws_iam_user" "sqsc" {
  name = "squarescale"
  path = "/voogle-sogilis-dev/"
}

resource "aws_iam_user_policy_attachment" "sqsc-attach" {
  user       = aws_iam_user.sqsc.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}

data "aws_caller_identity" "current" {}

output "account-id" {
  value = data.aws_caller_identity.current.account_id
}

output "token-id" {
  value = aws_iam_access_key.voogle-sqsc.id
}
output "secret" {
  value     = aws_iam_access_key.voogle-sqsc.secret
  sensitive = true
}


# IAM voogle S3 bucket
resource "aws_iam_user" "voogle-s3" {
  name = "voogle-s3"
  path = "/voogle-sogilis-dev/"
}

resource "aws_iam_access_key" "voogle-s3-access-key" {
  user = aws_iam_user.voogle-s3.name
}

resource "aws_iam_user_group_membership" "voogle-s3-group-membership" {
  user = aws_iam_user.voogle-s3.name

  groups = [
    aws_iam_group.voogle-s3-group.name,
  ]
}

output "voogle-s3-token-id" {
  value = aws_iam_access_key.voogle-s3-access-key.id
}

output "voogle-s3-secret" {
  value     = aws_iam_access_key.voogle-s3-access-key.secret
  sensitive = true
}
