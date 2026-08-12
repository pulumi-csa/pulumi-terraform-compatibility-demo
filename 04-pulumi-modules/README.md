# Pulumi Python — Using the Published Module

A Python Pulumi program that consumes the [`encrypted-s3-bucket`](../01-terraform-module) module after it has been published to the Pulumi private registry and converted to a typed SDK.

## Prerequisites

The module must be published to Pulumi Cloud first. See the [publishing instructions](../01-terraform-module/README.md#publishing-to-pulumi) in `01-terraform-module`.

## Setup

```bash
# Install the generated module SDK alongside other dependencies
pulumi package add encrypted-s3-bucket-aws@4.0.0

# Install dependencies
pulumi install
```

## Usage

```bash
pulumi up
```

Variables are set via Pulumi config:

```bash
pulumi config set project myapp
pulumi config set environment development
```

## What this shows

The same infrastructure as `02-terraform-code` and `03-HCL-language`, but:

- Written in **Python** using native Pulumi idioms
- The S3 module is consumed as a **typed SDK** — inputs and outputs are fully typed, not raw HCL strings
- The KMS key uses the native `pulumi_aws` provider directly, while the buckets use the published module
- State, secrets, and policy enforcement are all managed by Pulumi Cloud automatically
