# Terraform + Pulumi Demo

This demo shows how Pulumi can enhance an existing Terraform workflow — starting with a plain Terraform project.

## What's in this repo

| Directory                                      | What it is                                                                       |
| ---------------------------------------------- | -------------------------------------------------------------------------------- |
| [`01-terraform-module`](./01-terraform-module) | A reusable Terraform module for an encrypted S3 bucket                           |
| [`02-terraform-code`](./02-terraform-code)     | A root module that calls `01-terraform-module` and provisions real AWS resources |
| [`03-HCL-language`](./03-HCL-language)         | The same infrastructure as a native Pulumi HCL program                           |
| [`04-pulumi-modules`](./04-pulumi-modules)     | Python Pulumi program consuming the module via generated typed SDK               |
| [`upload-module`](./upload-module)             | Go helper for publishing `01-terraform-module` to the Pulumi private registry    |

## Which approach should you use?

|                    | Terraform backend (`02`)                                                                                 | HCL language (`03`)                                                                                     |
| ------------------ | -------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| **Best for**       | Lifting existing Terraform state into Pulumi Cloud with zero code changes                                | Reusing existing HCL code as a Pulumi-native program                                                    |
| **Code changes**   | None — same HCL, same CLI                                                                                | Add `Pulumi.yaml`; swap `terraform` commands for `pulumi`                                               |
| **State**          | Migrated from local/S3 to Pulumi Cloud                                                                   | Managed by Pulumi Cloud from the start                                                                  |
| **When to choose** | Team wants Pulumi Cloud visibility and policies immediately, without disrupting their Terraform workflow | Developers comfortable in HCL who want to adopt Pulumi gradually, or want to reuse existing module code |

---

## Step 1 — Stand up the infrastructure with local state

Start in `02-terraform-code`. At this point the project is plain Terraform — no Pulumi involvement yet.

```bash
cd 02-terraform-code
terraform init
terraform apply -var="project=myapp" -var="environment=development"
```

State is written to a local `terraform.tfstate` file. This is the starting point your customer is likely at today.

---

## Step 2 — Switch to Pulumi Cloud as the state backend

### 2a. Back up existing state

Do this before touching any configuration:

```bash
terraform state pull > terraform_state_backup.tfstate
```

### 2b. Authenticate with Pulumi Cloud

```bash
terraform login tf.pulumi.com
```

Your browser will open to the Pulumi Cloud access tokens page. Click **Create token**, enter a description and expiration, then copy the generated token and paste it into the terminal prompt.

### 2c. Uncomment the backend block

In `02-terraform-code/versions.tf`, uncomment the `backend "remote"` block and fill in your organization and workspace name:

```hcl
backend "remote" {
  hostname     = "tf.pulumi.com"
  organization = "HuckStream"

  workspaces {
    name = "terraform_encryptedBucket"
  }
}
```

### 2d. Migrate state

```bash
terraform init -migrate-state
```

Terraform will ask to copy existing state to the new backend — confirm yes.

### 2e. Set execution mode to local

New Pulumi workspaces default to remote execution, which can't reach local module paths. Set the workspace to run operations locally:

```bash
pulumi stack tag set terraform:execution-mode local -s HuckStream/terraform/encryptedBucket
```

Verify the migration with:

```bash
terraform plan
```

It should report no changes.

---

## Step 3 — Explore what Pulumi Cloud gives you

Open the stack in Pulumi Cloud: https://app.pulumi.com/HuckStream/terraform/encryptedBucket

### State visibility

To generate some history to show, make a real infrastructure change. In `02-terraform-code/main.tf`, set `versioning_enabled = false` on the `app_bucket` module, then apply:

```bash
terraform apply -var="project=myapp" -var="environment=development"
```

Now open the stack in Pulumi Cloud and navigate to the **Activity** tab:
https://app.pulumi.com/HuckStream/terraform/encryptedBucket

You'll see two updates — the initial apply and this change. Click into an update to show the diff: which resources were modified, what properties changed, and a full timeline. This is state history that Terraform alone can't give you — no more opaque `.tfstate` files on someone's laptop or in a shared S3 bucket.

---

## Step 4 — Run the same infrastructure as a native Pulumi HCL program

Switch to `03-HCL-language`. The `.tf` files are nearly identical to `02`, but this is now a Pulumi-native program — no backend configuration needed, and state is managed by Pulumi Cloud automatically.

```bash
cd ../03-HCL-language

# Fetch provider SDKs
pulumi install

# Deploy
pulumi up
```

### Policy enforcement

Add this stack to the **hcl-language-demo** stack group in Pulumi Cloud, then run `pulumi up`. During the preview, Pulumi automatically evaluates the attached policy pack and surfaces a critical advisory:

```
⚠️  pulumi-best-practices-aws@v1.4.1
    - [advisory] [severity: critical]  cloudtrail-enabled
      CloudTrail must be enabled with at least 1 trail(s). Found 0 trail(s).
```

This fires without any changes to the HCL — policy enforcement is applied at the org level, across every stack in the group.

### Tear down

```bash
pulumi destroy
```
