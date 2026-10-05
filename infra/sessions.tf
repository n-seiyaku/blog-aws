resource "aws_dynamodb_table" "sessions" {
  name         = "blog-sessions"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "refreshTokenHash"
    type = "S"
  }

  attribute {
    name = "userId"
    type = "S"
  }

  global_secondary_index {
    name            = "refreshTokenHash-index"
    projection_type = "ALL"

    key_schema {
      attribute_name = "refreshTokenHash"
      key_type       = "HASH"
    }
  }

  global_secondary_index {
    name            = "userId-index"
    projection_type = "ALL"

    key_schema {
      attribute_name = "userId"
      key_type       = "HASH"
    }
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }
}