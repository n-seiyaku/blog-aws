
resource "aws_s3_bucket" "post_images" {
  bucket_prefix = "blog-aws-post-images-"

  tags = {
    Project     = "blog-aws"
    Environment = "dev"
    Purpose     = "post-images"
  }
}

# Block all public access
resource "aws_s3_bucket_public_access_block" "post_images" {
  bucket = aws_s3_bucket.post_images.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Bucket owner manages all objects
resource "aws_s3_bucket_ownership_controls" "post_images" {
  bucket = aws_s3_bucket.post_images.id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

# Server-side encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "post_images" {
  bucket = aws_s3_bucket.post_images.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Allow React frontend to upload via presigned URL
resource "aws_s3_bucket_cors_configuration" "post_images" {
  bucket = aws_s3_bucket.post_images.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "PUT", "HEAD"]
    allowed_origins = ["http://localhost:5173"]
    expose_headers  = ["ETag"]
    max_age_seconds = 300
  }
}

output "s3_post_images_bucket" {
  description = "S3 bucket name for blog post images"
  value       = aws_s3_bucket.post_images.bucket
}

output "s3_post_images_arn" {
  description = "S3 bucket ARN"
  value       = aws_s3_bucket.post_images.arn
}
