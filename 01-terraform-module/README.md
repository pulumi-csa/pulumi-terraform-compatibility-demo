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

## Publishing to Pulumi

Publishing this module to Pulumi's private registry converts it into a multi-language Pulumi package automatically. Requires a Pulumi Cloud Enterprise or Business Critical plan.

### 1. Set your access token

```bash
export PULUMI_ACCESS_TOKEN=<your-access-token>
```

### 2. Publish the module

Run the `upload-module` Go helper from the **repo root** (not from within this directory). The first time, fetch dependencies first:

```bash
cd upload-module && go mod tidy && cd ..
go -C upload-module run . \
    -org HuckStream \
    -name encrypted-s3-bucket \
    -provider aws \
    -version 1.0.0 \
    -path ./01-terraform-module
```

You should see:

```
creating module HuckStream/encrypted-s3-bucket/aws
uploaded HuckStream/encrypted-s3-bucket/aws@1.0.0
```

Navigate to **Platform → Private components** in Pulumi Cloud to confirm the module is published.

### 3. Use the module in a Pulumi program

Once published, the module can be referenced in Terraform configs:

```hcl
module "app-bucket" {
  source  = "tf.pulumi.com/HuckStream/encrypted-s3-bucket/aws"
  version = "1.0.0"
}
```

Or installed as a typed SDK in any Pulumi language project:

```bash
pulumi package add encrypted-s3-bucket-aws@1.0.0
```

See [`04-pulumi-modules`](../04-pulumi-modules) for a Python example.

## Outputs

| Name | Description |
|------|-------------|
| `bucket_id` | Bucket name |
| `bucket_arn` | Bucket ARN |
| `bucket_domain_name` | Global bucket domain |
| `bucket_regional_domain_name` | Region-specific bucket domain |
