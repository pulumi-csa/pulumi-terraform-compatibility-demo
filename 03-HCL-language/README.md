# Pulumi HCL Language — Encrypted S3 Buckets

The same infrastructure as [`02-terraform-code`](../02-terraform-code), written as a native Pulumi program in HCL. The `.tf` files are nearly identical — the only difference is that Pulumi takes over state management, the CLI, and the deployment lifecycle.

## What's different from Terraform

|                 | Terraform           | Pulumi HCL              |
| --------------- | ------------------- | ----------------------- |
| Project file    | none                | `Pulumi.yaml`           |
| State backend   | configured in HCL   | managed by Pulumi Cloud |
| `backend` block | required            | ignored                 |
| Deploy          | `terraform apply`   | `pulumi up`             |
| Destroy         | `terraform destroy` | `pulumi destroy`        |
| Preview         | `terraform plan`    | `pulumi preview`        |
| Outputs         | `terraform output`  | `pulumi stack output`   |

## Usage

```bash
# Fetch provider SDKs (run once, or after changing providers)
pulumi install

# Deploy
pulumi up

# View outputs
pulumi stack output

# Tear down
pulumi destroy
```

## Seeing policies in action

Add this stack to the **hcl-language-demo** stack group in Pulumi Cloud. Once added, run `pulumi up` and you'll see the policy pack kick in during the preview:

```
Policies:
  ⚠️  pulumi-best-practices-aws@v1.4.1
      - [advisory] [severity: critical]  cloudtrail-enabled  (pulumi:pulumi:Stack: hcl-encrypted-s3-demo)
        CloudTrail must be enabled with at least 1 trail(s). Found 0 trail(s).
```

This fires because the stack has no `aws_cloudtrail` resource. It's a good example of how Pulumi Cloud enforces best practices automatically across any stack in the group — without changing a line of HCL.
