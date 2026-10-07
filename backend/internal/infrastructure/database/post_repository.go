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
	FeedKey   string `dynamodbav:"feedKey"`
	Title     string `dynamodbav:"title"`
	Content   string `dynamodbav:"content"`
	AuthorID  string `dynamodbav:"authorId"`
	CreatedAt string `dynamodbav:"createdAt"`
	UpdatedAt string `dynamodbav:"updatedAt"`
}

const FeedKey = "POST"

type PostRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewPostRepository(client *dynamodb.Client, tableName string) *PostRepository {
	return &PostRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *PostRepository) Create(ctx context.Context, p post.Post) error {
	item := postItem{
		ID:        p.ID,
		FeedKey:   FeedKey,
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

func (r *PostRepository) Update(ctx context.Context, p post.Post) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: p.ID,
			},
		},
		UpdateExpression: aws.String("SET #title = :title, #content = :content, #updatedAt = :updatedAt"),
		ExpressionAttributeNames: map[string]string{
			"#title":     "title",
			"#content":   "content",
			"#updatedAt": "updatedAt",
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

func (r *PostRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})

	return err
}

const (
	FeedKeyPost = "POST"
)

func (r *PostRepository) ListPosts(ctx context.Context, limit int32, cursor string) ([]post.Post, string, error) {
	lastKey, err := decodePostCursor(cursor)
	if err != nil {
		return nil, "", err
	}

	output, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("feedKey-createdAt-index"),
		KeyConditionExpression: aws.String("feedKey = :feedKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":feedKey": &types.AttributeValueMemberS{
				Value: FeedKeyPost,
			},
		},
		Limit:             aws.Int32(limit),
		ScanIndexForward:  aws.Bool(false),
		ExclusiveStartKey: lastKey,
	})
	if err != nil {
		return nil, "", err
	}

	posts := make([]post.Post, 0, len(output.Items))
	for _, item := range output.Items {
		var p postItem

		if err := attributevalue.UnmarshalMap(item, &p); err != nil {
			return nil, "", err
		}

		post, err := parsePostItem(p)
		if err != nil {
			return nil, "", err
		}

		posts = append(posts, post)
	}

	nextCursor, err := encodePostCursor(output.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return posts, nextCursor, nil
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
