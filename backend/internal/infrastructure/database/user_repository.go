package database

import (
	"blog-aws-backend/internal/domain/user"
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type userItem struct {
	ID           string `dynamodbav:"id"`
	Email        string `dynamodbav:"email"`
	PasswordHash string `dynamodbav:"passwordHash"`
	Role         string `dynamodbav:"role"`
	CreatedAt    string `dynamodbav:"createdAt"`
	UpdatedAt    string `dynamodbav:"updatedAt"`
	DeletedAt    string `dynamodbav:"deletedAt"`
}

type UserRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepository(client *dynamodb.Client, tableName string) *UserRepository {
	return &UserRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *UserRepository) Create(ctx context.Context, u user.User) error {
	item := userItem{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	// id が重複していないことを確認して挿入
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &r.tableName,
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	})

	if err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			return user.ErrorEmailAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, u user.User) error {
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (user.User, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      &r.tableName,
		Key:            stringAttribute("id", id),
		ConsistentRead: aws.Bool(true),
	})

	if err != nil {
		return user.User{}, err
	}

	if len(result.Item) == 0 {
		return user.User{}, user.ErrorNotFound
	}

	var item userItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return user.User{}, err
	}

	return parseUserItem(item)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (user.User, error) {
	// GSI (email-index) を使用してメールアドレスで検索
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &r.tableName,
		IndexName:                 aws.String("email-index"),
		KeyConditionExpression:    aws.String("email = :email"),
		ExpressionAttributeValues: stringAttribute(":email", email),
	})

	if err != nil {
		return user.User{}, err
	}

	if len(result.Items) == 0 {
		return user.User{}, user.ErrorNotFound
	}

	var item userItem
	if err := attributevalue.UnmarshalMap(result.Items[0], &item); err != nil {
		return user.User{}, err
	}

	return parseUserItem(item)
}

// userItem を user.User に変換するヘルパー関数
func parseUserItem(item userItem) (user.User, error) {
	var createdAt time.Time
	var updatedAt time.Time
	var deletedAt *time.Time
	var err error

	if createdAt, err = time.Parse(time.RFC3339, item.CreatedAt); err != nil {
		return user.User{}, err
	}

	if item.UpdatedAt != "" {
		if updatedAt, err = time.Parse(time.RFC3339, item.UpdatedAt); err != nil {
			return user.User{}, err
		}
	}

	if item.DeletedAt != "" {
		parsed, err := time.Parse(time.RFC3339, item.DeletedAt)
		if err != nil {
			return user.User{}, err
		}

		deletedAt = &parsed
	}

	return user.User{
		ID:           item.ID,
		Email:        item.Email,
		PasswordHash: item.PasswordHash,
		Role:         user.Role(item.Role),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		DeletedAt:    deletedAt,
	}, nil
}
