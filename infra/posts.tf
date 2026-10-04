resource "aws_dynamodb_table" "posts" {
  name         = "posts"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "authorId"
    type = "S"
  }

  attribute {
    name = "feedKey"
    type = "S"
  }

  attribute {
    name = "createdAt"
    type = "S"
  }

  global_secondary_index {
    name            = "feedKey-createdAt-index"
    projection_type = "ALL"

    key_schema {
      attribute_name = "feedKey"
      key_type       = "HASH"
    }

    key_schema {
      attribute_name = "createdAt"
      key_type       = "RANGE"
    }
  }

  global_secondary_index {
    name            = "authorId-index"
    projection_type = "ALL"

    key_schema {
      attribute_name = "authorId"
      key_type       = "HASH"
    }
  }
}
