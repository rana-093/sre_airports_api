variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-south-1"
}

variable "acl" {
  description = "ACL"
  type        = string
  default     = "public"
}

variable "environment" {
  description = "Environment tag"
  type        = string
  default     = "Dev"
}
