# Terraform Example — Using the Encrypted S3 Module

Example Terraform root module that calls the [`encrypted-s3-bucket` module](../01-terraform-module) twice to provision two buckets with different configurations.

## What it creates

- A KMS key (with rotation) for app data encryption
- An **app data bucket** — SSE-KMS, versioning enabled
- A **logs bucket** — AES256, versioning disabled

## Usage

```bash
terraform init
terraform apply -var="project=myapp" -var="environment=development"
```

## Variables

| Name | Default | Description |
|------|---------|-------------|
| `project` | `myapp` | Prefix for resource names |
| `environment` | `development` | Deployment environment |
| `aws_region` | `us-east-1` | AWS region |
