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

## Using Pulumi as the State Backend

Instead of storing state locally or in S3, you can use [Pulumi Cloud](https://app.pulumi.com) as a Terraform state backend. This gives you state history, drift detection, and a UI to browse your resources — with no extra infrastructure to manage.

> **If you have existing local state, back it up before making any changes:**
>
> ```bash
> terraform state pull > terraform_state_backup.tfstate
> ```

### 1. Authenticate

```bash
terraform login tf.pulumi.com
```

This opens your browser to the Pulumi Cloud access tokens page. From there:

1. Click **Create token**
2. Enter a description (e.g. `used for terraform backend`) and choose an expiration
3. Click **Create token** — Pulumi will display the token once; copy it
4. Paste the token into the terminal prompt and press Enter

The token is saved to `~/.terraform.d/credentials.tfrc.json` for all future operations. For CI/CD, set the environment variable instead:

```bash
export TF_TOKEN_tf_pulumi_com=<your-pulumi-access-token>
```

### 2. Configure the backend

Uncomment the `backend "remote"` block in [`versions.tf`](./versions.tf) and fill in your organization and workspace name:

```hcl
backend "remote" {
  hostname     = "tf.pulumi.com"
  organization = "<your-pulumi-org>"

  workspaces {
    name = "<project>_<stack>"
  }
}
```

### 3. Migrate existing state

Re-initialize Terraform to migrate your state to Pulumi Cloud:

```bash
terraform init -migrate-state
```

Verify with `terraform plan` — it should report no changes.

> **Local modules require local execution mode.** The `remote` backend defaults to running operations in Pulumi Cloud, where local paths like `../01-terraform-module` are not accessible. Set execution mode to `local` so that `terraform`/`tofu` runs on your machine and only state is stored in Pulumi Cloud:
>
> ```bash
> pulumi stack tag set terraform:execution-mode local -s HuckStream/terraform/encryptedBucket
> ```
>
> Alternatively, set the `terraform:execution-mode` tag to `local` in Pulumi Cloud.

Once migrated, view the stack and its state history in Pulumi Cloud:
https://app.pulumi.com/<your-organization-name>/<selected-project-name>/<selected-stack-name>

In this example, that's https://app.pulumi.com/HuckStream/terraform/encryptedBucket

## Variables

| Name          | Default       | Description               |
| ------------- | ------------- | ------------------------- |
| `project`     | `myapp`       | Prefix for resource names |
| `environment` | `development` | Deployment environment    |
| `aws_region`  | `us-east-1`   | AWS region                |
