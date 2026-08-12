# Encrypted S3 Bucket — Terraform Module

A reusable Terraform module that provisions an encrypted S3 bucket with public access blocked and versioning support.

## Resources

- `aws_s3_bucket` — the bucket
- `aws_s3_bucket_server_side_encryption_configuration` — SSE-KMS (if a KMS key ARN is provided) or AES256
- `aws_s3_bucket_public_access_block` — all public access blocked
- `aws_s3_bucket_versioning` — enabled/suspended via variable

## Inputs

| Name | Type | Default | Description |
|------|------|---------|-------------|
| `bucket_name` | `string` | — | Name of the S3 bucket |
| `kms_key_arn` | `string` | `null` | KMS key ARN for SSE-KMS; falls back to AES256 if null |
| `versioning_enabled` | `bool` | `true` | Enable object versioning |
| `force_destroy` | `bool` | `false` | Allow bucket deletion even when non-empty |
| `tags` | `map(string)` | `{}` | Tags applied to all resources |

## Outputs

| Name | Description |
|------|-------------|
| `bucket_id` | Bucket name |
| `bucket_arn` | Bucket ARN |
| `bucket_domain_name` | Global bucket domain |
| `bucket_regional_domain_name` | Region-specific bucket domain |
