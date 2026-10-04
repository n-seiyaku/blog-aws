package database

import (
	"blog-aws-backend/internal/domain/post"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type postItem struct {
	ID        string `dynamodbav:"id"`
	Title     string `dynamodbav:"title"`
	AuthorID  string `dynamodbav:"authorId"`
	Content   string `dynamodbav:"content"`
	CreatedAt string `dynamodbav:"createdAt"`
	UpdatedAt string `dynamodbav:"updatedAt"`
}

type PostRepository struct {
	client    *dynamodb.Client
	tableName string
}

const LIMIT_POST = 10

func NewPostRepository(client *dynamodb.Client, tableName string) *PostRepository {
	return &PostRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *PostRepository) Create(ctx context.Context, p *post.Post) error {
	item := postItem{
		ID:        p.ID,
		Title:     p.Title,
		AuthorID:  p.AuthorID,
		Content:   p.Content,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      av,
		// idが重複していないことを確認して挿入
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	})

	return err
}

func (r *PostRepository) Update(ctx context.Context, p *post.Post) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: p.ID,
			},
		},
		UpdateExpression: aws.String("SET #title = :title, #content = :content, #updateAt = :updateAt"),
		ExpressionAttributeNames: map[string]string{
			"#title":    "title",
			"#content":  "content",
			"#updateAt": "updateAt",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":title": &types.AttributeValueMemberS{
				Value: p.Title,
			},
			":content": &types.AttributeValueMemberS{
				Value: p.Content,
			},
			":updatedAt": &types.AttributeValueMemberS{
				Value: p.UpdatedAt.Format(time.RFC3339),
			},
		},
    ConditionExpression: aws.String("attribute_exists(id)"),
	})
	return err
}

func (r *PostRepository) Delete(ctx context.Context, p *post.Post) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
    TableName: &r.tableName,
    Key: map[string]types.AttributeValue{
      "id": &types.AttributeValueMemberS{
        Value: p.ID,
      },
    },
  })
  
  return err
}

func (r *PostRepository) GetAll(ctx context.Context) ([]post.Post, error) {
	items, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:            &r.tableName,
		Limit:                aws.Int32(LIMIT_POST),
		ProjectionExpression: aws.String("id, title, authorId, content, createdAt, updatedAt"),
	})
	if err != nil {
		return nil, err
	}

  var posts []post.Post
  for _, item := range items.Items {
    var p postItem
    if err := attributevalue.UnmarshalMap(item, &p); err != nil {
      return nil, err
    }

    post, err := parsePostItem(p)
    if err != nil {
      return nil, err
    }
    posts = append(posts, post)
  }

  return posts, nil
}

func parsePostItem(item postItem) (post.Post, error) {
  createdAt, err := time.Parse(time.RFC3339, item.CreatedAt)
  if err != nil {
    return post.Post{}, err
  }

  updatedAt, err := time.Parse(time.RFC3339, item.UpdatedAt)
  if err != nil {
    return post.Post{}, err
  }

  return post.Post{
    ID:        item.ID,
    Title:     item.Title,
    AuthorID:  item.AuthorID,
    Content:   item.Content,
    CreatedAt: createdAt,
    UpdatedAt: updatedAt,
  }, nil
}
