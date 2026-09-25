"""Encrypted S3 buckets and KMS key, migrated from Terraform to native Pulumi.

Mirrors the infrastructure previously defined in ``02-terraform-code`` (root
config) and ``01-terraform-module`` (reusable encrypted-bucket module):

- A KMS key (with rotation) used to encrypt the app data bucket
- An app data bucket - SSE-KMS, versioning enabled
- A logs bucket - AES256, versioning currently enabled in the real
  infrastructure (see ``logs_bucket_versioning_enabled`` below)
"""

import pulumi
import pulumi_aws as aws


class EncryptedS3BucketArgs:
    """Arguments for `EncryptedS3Bucket`, mirroring `01-terraform-module`."""

    def __init__(
        self,
        bucket_name: pulumi.Input[str],
        kms_key_arn: pulumi.Input[str] | None = None,
        versioning_enabled: bool = True,
        force_destroy: bool = False,
        tags: dict[str, str] | None = None,
    ) -> None:
        self.bucket_name = bucket_name
        self.kms_key_arn = kms_key_arn
        self.versioning_enabled = versioning_enabled
        self.force_destroy = force_destroy
        self.tags = tags or {}


class EncryptedS3Bucket(pulumi.ComponentResource):
    """An S3 bucket with public access blocked, SSE, and versioning support.

    Mirrors `01-terraform-module`: SSE-KMS if `kms_key_arn` is given,
    otherwise AES256.
    """

    bucket: aws.s3.Bucket
    bucket_id: pulumi.Output[str]
    bucket_arn: pulumi.Output[str]

    def __init__(
        self,
        name: str,
        args: EncryptedS3BucketArgs,
        opts: pulumi.ResourceOptions | None = None,
    ) -> None:
        super().__init__("terraform:encryptedS3:EncryptedBucket", name, {}, opts)

        # These resources were previously top-level (root-stack-parented)
        # when the state was migrated from Terraform; the aliases preserve
        # their identity now that they live under this component.
        child_opts = pulumi.ResourceOptions(
            parent=self,
            aliases=[pulumi.Alias(parent=pulumi.ROOT_STACK_RESOURCE)],
        )

        self.bucket = aws.s3.Bucket(
            name,
            bucket=args.bucket_name,
            force_destroy=args.force_destroy,
            tags=args.tags,
            opts=child_opts,
        )

        if args.kms_key_arn is not None:
            sse_default = aws.s3.BucketServerSideEncryptionConfigurationRuleApplyServerSideEncryptionByDefaultArgs(
                sse_algorithm="aws:kms",
                kms_master_key_id=args.kms_key_arn,
            )
            bucket_key_enabled = True
        else:
            sse_default = aws.s3.BucketServerSideEncryptionConfigurationRuleApplyServerSideEncryptionByDefaultArgs(
                sse_algorithm="AES256",
            )
            bucket_key_enabled = False

        aws.s3.BucketServerSideEncryptionConfiguration(
            name,
            bucket=self.bucket.id,
            rules=[
                aws.s3.BucketServerSideEncryptionConfigurationRuleArgs(
                    apply_server_side_encryption_by_default=sse_default,
                    bucket_key_enabled=bucket_key_enabled,
                )
            ],
            opts=child_opts,
        )

        aws.s3.BucketPublicAccessBlock(
            name,
            bucket=self.bucket.id,
            block_public_acls=True,
            block_public_policy=True,
            ignore_public_acls=True,
            restrict_public_buckets=True,
            opts=child_opts,
        )

        aws.s3.BucketVersioning(
            name,
            bucket=self.bucket.id,
            versioning_configuration=aws.s3.BucketVersioningVersioningConfigurationArgs(
                status="Enabled" if args.versioning_enabled else "Suspended",
            ),
            opts=child_opts,
        )

        self.bucket_id = self.bucket.id
        self.bucket_arn = self.bucket.arn
        self.register_outputs(
            {
                "bucket_id": self.bucket_id,
                "bucket_arn": self.bucket_arn,
            }
        )


config = pulumi.Config()
project = config.get("project") or "myapp"
environment = config.get("environment") or "development"

# The original Terraform code set versioning_enabled=false for the logs
# bucket, but the deployed bucket actually has versioning Enabled. Default
# to the real, deployed state so importing this program produces no diff;
# set this config value to false to suspend versioning on the logs bucket.
logs_bucket_versioning_enabled = config.get_bool("logsBucketVersioningEnabled")
if logs_bucket_versioning_enabled is None:
    logs_bucket_versioning_enabled = True

# ManagedBy is kept as "terraform" to match the tags already applied to the
# real resources, so importing this program produces no diff. Update it (and
# re-apply) once you're ready to record the new ownership.
common_tags = {
    "Project": project,
    "Environment": environment,
    "ManagedBy": "terraform",
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
    name=f"alias/{project}-s3",
    target_key_id=kms_key.key_id,
)

app_bucket = EncryptedS3Bucket(
    "app_bucket_this",
    EncryptedS3BucketArgs(
        bucket_name=f"{project}-{environment}-app-data",
        kms_key_arn=kms_key.arn,
        versioning_enabled=True,
        force_destroy=environment != "production",
        tags={**common_tags, "Purpose": "app-data"},
    ),
)

logs_bucket = EncryptedS3Bucket(
    "logs_bucket_this",
    EncryptedS3BucketArgs(
        bucket_name=f"{project}-{environment}-logs",
        kms_key_arn=None,
        versioning_enabled=logs_bucket_versioning_enabled,
        force_destroy=True,
        tags={**common_tags, "Purpose": "logs"},
    ),
)

pulumi.export("app_bucket_id", app_bucket.bucket_id)
pulumi.export("app_bucket_arn", app_bucket.bucket_arn)
pulumi.export("logs_bucket_id", logs_bucket.bucket_id)
pulumi.export("logs_bucket_arn", logs_bucket.bucket_arn)
pulumi.export("kms_key_arn", kms_key.arn)
