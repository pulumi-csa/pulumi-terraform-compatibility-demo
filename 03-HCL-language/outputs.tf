output "app_bucket_id" {
  description = "Name of the app data S3 bucket"
  value       = module.app_bucket.bucket_id
}

output "app_bucket_arn" {
  description = "ARN of the app data S3 bucket"
  value       = module.app_bucket.bucket_arn
}

output "logs_bucket_id" {
  description = "Name of the logs S3 bucket"
  value       = module.logs_bucket.bucket_id
}

output "logs_bucket_arn" {
  description = "ARN of the logs S3 bucket"
  value       = module.logs_bucket.bucket_arn
}

output "kms_key_arn" {
  description = "ARN of the KMS key used for app bucket encryption"
  value       = aws_kms_key.s3.arn
}
