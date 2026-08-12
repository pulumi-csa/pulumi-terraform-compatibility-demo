variable "project" {
  description = "Project name, used as a prefix for resource names"
  type        = string
  default     = "myapp"
}

variable "environment" {
  description = "Deployment environment (e.g. development, staging, production)"
  type        = string
  default     = "development"
}

variable "aws_region" {
  description = "AWS region to deploy resources into"
  type        = string
  default     = "us-east-1"
}
