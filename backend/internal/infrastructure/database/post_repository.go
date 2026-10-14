package database

import (
	"blog-aws-backend/internal/domain/post"
	"context"
	"errors"
	"fmt"
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
	ImageKey  string `dynamodbav:"imageKey,omitempty"`
	AuthorID  string `dynamodbav:"authorId"`
	CreatedAt string `dynamodbav:"createdAt"`
	UpdatedAt string `dynamodbav:"updatedAt"`
}

const FeedKeyPost = "POST"

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
		FeedKey:   FeedKeyPost,
		Title:     p.Title,
		AuthorID:  p.AuthorID,
		ImageKey:  *p.ImageKey,
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

func (r *PostRepository) Update(ctx context.Context, p post.Post, expectedImageKey string) error {
	condition := "attribute_exists(#id) AND #authorId = :authorId AND #imageKey = :oldImageKey"

	if expectedImageKey == "" {
		condition = "attribute_exists(#id) AND #authorId = :authorId AND (#imageKey = :oldImageKey OR attribute_not_exists(#imageKey))"
	}

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: p.ID,
			},
		},
		UpdateExpression:    aws.String("SET #title = :title, #content = :content, #imageKey = :imageKey, #updatedAt = :updatedAt"),
		ConditionExpression: aws.String(condition),
		ExpressionAttributeNames: map[string]string{
			"#id":        "id",
			"#authorId":  "authorId",
			"#title":     "title",
			"#content":   "content",
			"#imageKey":  "imageKey",
			"#updatedAt": "updatedAt",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":title": &types.AttributeValueMemberS{
				Value: p.Title,
			},
			":content": &types.AttributeValueMemberS{
				Value: p.Content,
			},
			":imageKey": &types.AttributeValueMemberS{
				Value: *p.ImageKey,
			},
			":updatedAt": &types.AttributeValueMemberS{
				Value: p.UpdatedAt.Format(time.RFC3339),
			},
			":authorId": &types.AttributeValueMemberS{
				Value: p.AuthorID,
			},
			":oldImageKey": &types.AttributeValueMemberS{
				Value: expectedImageKey,
			},
		},
	})
	if err != nil {
		var conditionalErr *types.ConditionalCheckFailedException

		if errors.As(err, &conditionalErr) {
			return post.ErrConflict
		}

		return fmt.Errorf("update post: %w", err)
	}

	return nil
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

func (r *PostRepository) FindById(ctx context.Context, id string) (post.Post, error) {
	output, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return post.Post{}, err
	}

	if len(output.Item) == 0 {
		return post.Post{}, post.ErrNotFound
	}

	var p postItem
	if err := attributevalue.UnmarshalMap(output.Item, &p); err != nil {
		return post.Post{}, err
	}

	return parsePostItem(p)
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
		ImageKey:  &item.ImageKey,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
