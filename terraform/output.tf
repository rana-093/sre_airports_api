output "bucket_id" {
  description = "ID of S3 bucket"
  value       = aws_s3_bucket.airport_test_bucket.id
}

output "bucket_arn" {
  description = "ARN of S3 bucket"
  value       = aws_s3_bucket.airport_test_bucket.arn
}
