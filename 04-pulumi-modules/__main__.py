import pulumi
import pulumi_aws as aws
import pulumi_encrypted_s3_bucket_aws as s3_module

config = pulumi.Config()
project = config.get("project") or "pyapp"
environment = config.get("environment") or "development"

common_tags = {
    "Project": project,
    "Environment": environment,
    "ManagedBy": "pulumi",
}

kms_key = aws.kms.Key(
    "s3",
    description="KMS key for S3 bucket encryption",
    deletion_window_in_days=7,
    enable_key_rotation=True,
    tags=common_tags,
)

aws.kms.Alias(
    "s3",
    name=pulumi.Output.concat("alias/", project, "-s3"),
    target_key_id=kms_key.key_id,
)

app_bucket = s3_module.Module(
    "app-bucket",
    bucket_name=pulumi.Output.concat(project, "-", environment, "-app-data"),
    kms_key_arn=kms_key.arn,
    versioning_enabled=True,
    force_destroy=environment != "production",
    tags={**common_tags, "Purpose": "app-data"},
)

logs_bucket = s3_module.Module(
    "logs-bucket",
    bucket_name=pulumi.Output.concat(project, "-", environment, "-logs"),
    versioning_enabled=False,
    force_destroy=True,
    tags={**common_tags, "Purpose": "logs"},
)

pulumi.export("app_bucket_id", app_bucket.bucket_id)
pulumi.export("app_bucket_arn", app_bucket.bucket_arn)
pulumi.export("logs_bucket_id", logs_bucket.bucket_id)
pulumi.export("logs_bucket_arn", logs_bucket.bucket_arn)
pulumi.export("kms_key_arn", kms_key.arn)
