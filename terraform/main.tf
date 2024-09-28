provider "aws" {
  region = var.aws_region

resource "aws_s3_bucket" "airport_test_bucket" {
  bucket = "airport_test_bucket"
  acl    = var.acl

  tags = {
    Name        = "airport_test_bucket"
    Environment = var.environment
  }
}
