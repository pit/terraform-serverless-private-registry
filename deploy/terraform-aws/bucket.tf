resource "aws_s3_bucket" "this" {
  bucket = var.storage_bucket_name
  acl    = "private"

  versioning {
    enabled = true
  }

  object_lock_configuration {
    object_lock_enabled = "Enabled"
    # rule {
    #   default_retention {
    #     mode = "GOVERNANCE"
    #   }
    # }
  }

  tags = merge({
    Name = var.storage_bucket_name
  }, var.tags)
}

resource "aws_s3_bucket_policy" "this" {
  bucket = aws_s3_bucket.this.id
  policy = data.aws_iam_policy_document.this.json
}

data "aws_iam_policy_document" "this" {
  statement {
    sid = "AllowReadOnlyOperationsToReadOnlyArns"

    actions = [
      "s3:ListBucket",
      "s3:GetObject",
      "s3:GetObjectAcl",
      "s3:GetObjectTagging",
      "s3:GetObjectVersion",
    ]

    principals {
      type = "AWS"
      identifiers = sort(distinct(concat(
        var.storage_bucket_readonly_arns,
        var.storage_bucket_readwrite_arns,
        [module.lambdas_role.iam_role_arn],
      )))
    }

    resources = [
      aws_s3_bucket.this.arn,
      "${aws_s3_bucket.this.arn}/*",
    ]
  }

  statement {
    sid = "AllowReadWriteOperationsToReadWriteArns"

    actions = [
      "s3:PutObject",
      "s3:PutObjectAcl",
    ]

    principals {
      type = "AWS"
      identifiers = sort(distinct(concat(
        var.storage_bucket_readwrite_arns,
        [module.lambdas_role.iam_role_arn],
      )))
    }

    resources = [
      "${aws_s3_bucket.this.arn}/*",
    ]
  }

  statement {
    sid = "AllowAllOperationsToAdminArns"

    actions = [
      "s3:*",
    ]

    principals {
      type        = "AWS"
      identifiers = sort(distinct(var.storage_bucket_admin_arns))
    }

    resources = [
      aws_s3_bucket.this.arn,
      "${aws_s3_bucket.this.arn}/*",
    ]
  }
}

resource "aws_s3_bucket_public_access_block" "this" {
  bucket = aws_s3_bucket.this.id

  depends_on = [aws_s3_bucket_policy.this]

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# resource "aws_lambda_permission" "lambda_custom_signature_permission" {
#   function_name = module.lambda_custom_s3_signature.lambda_function_name

#   statement_id = "AllowInvokeFromS3Bucket"
#   action       = "lambda:InvokeFunction"
#   principal    = "s3.amazonaws.com"

#   source_arn = aws_s3_bucket.this.arn
# }
