resource "aws_kms_key" "s3" {
  description             = "KMS key for S3 bucket encryption"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = local.common_tags
}

resource "aws_kms_alias" "s3" {
  name          = "alias/${var.project}-s3"
  target_key_id = aws_kms_key.s3.key_id
}

module "app_bucket" {
  source = "../01-terraform-module"

  bucket_name        = "${var.project}-${var.environment}-app-data"
  kms_key_arn        = aws_kms_key.s3.arn
  versioning_enabled = true
  force_destroy      = var.environment != "production"

  tags = merge(local.common_tags, {
    Purpose = "app-data"
  })
}

module "logs_bucket" {
  source = "../01-terraform-module"

  bucket_name        = "${var.project}-${var.environment}-logs"
  kms_key_arn        = null
  versioning_enabled = false
  force_destroy      = true

  tags = merge(local.common_tags, {
    Purpose = "logs"
  })
}

locals {
  common_tags = {
    Project     = var.project
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}
