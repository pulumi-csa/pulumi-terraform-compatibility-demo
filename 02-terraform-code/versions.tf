terraform {
  required_version = ">= 1.3.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 6.0"
    }
  }

  # To use Pulumi as the state backend, uncomment this block and run:
  #   terraform login tf.pulumi.com
  #   terraform init -migrate-state
  #
  # backend "remote" {
  #   hostname     = "tf.pulumi.com"
  #   organization = "elisabeth-demo" #Your organization name

  #   workspaces {
  #     name = "terraform_encryptedBucket" #In the format projectName_stackName
  #   }
  # }
}

provider "aws" {
  region = var.aws_region
}
