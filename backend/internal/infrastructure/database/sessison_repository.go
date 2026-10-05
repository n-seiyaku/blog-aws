package database

import (
	"blog-aws-backend/internal/domain/session"
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type sessionItem struct {
	ID               string `dynamodbav:"id"`
	UserID           string `dynamodbav:"userId"`
	RefreshTokenHash string `dynamodbav:"refreshTokenHash"`
	CreatedAt        string `dynamodbav:"createdAt"`
	ExpiresAt        string `dynamodbav:"expiresAt"`
	TTL              int64  `dynamodbav:"ttl"`
}

type SessionRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewSessionRepository(client *dynamodb.Client, tableName string) *SessionRepository {
	return &SessionRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *SessionRepository) Create(ctx context.Context, s session.Session) error {
	item := sessionItem{
		ID:               s.ID,
		UserID:           s.UserID,
		RefreshTokenHash: s.RefreshTokenHash,
		CreatedAt:        s.CreatedAt.Format(time.RFC3339),
		ExpiresAt:        s.ExpiresAt.Format(time.RFC3339),
		TTL:              s.ExpiresAt.Unix(),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})

	return err
}

func (r *SessionRepository) FindByRefreshTokenHash(ctx context.Context, refreshToken string) (session.Session, error) {
	output, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              aws.String("refreshTokenHash-index"),
		KeyConditionExpression: aws.String("refreshTokenHash = :hash"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hash": &types.AttributeValueMemberS{
				Value: refreshToken,
			},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return session.Session{}, err
	}

	if len(output.Items) == 0 {
		return session.Session{}, session.ErrorNotFound
	}

	var item sessionItem
	if err := attributevalue.UnmarshalMap(output.Items[0], &item); err != nil {
		return session.Session{}, err
	}

	createdAt, err := time.Parse(time.RFC3339, item.CreatedAt)
	if err != nil {
		return session.Session{}, err
	}

	expiresAt, err := time.Parse(time.RFC3339, item.ExpiresAt)
	if err != nil {
		return session.Session{}, err
	}

	return session.Session{
		ID:               item.ID,
		UserID:           item.UserID,
		RefreshTokenHash: item.RefreshTokenHash,
		CreatedAt:        createdAt,
		ExpiresAt:        expiresAt,
	}, nil
}

func (r *SessionRepository) Rotate(ctx context.Context, sessionID string, refreshToken string, expiresAt time.Time) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: sessionID,
			},
		},
		UpdateExpression: aws.String(
			"SET refreshTokenHash = :hash, expiresAt = :expiresAt, #ttl = :ttl",
		),
		ExpressionAttributeNames: map[string]string{
			"#ttl": "ttl",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hash": &types.AttributeValueMemberS{
				Value: refreshToken,
			},
			":expiresAt": &types.AttributeValueMemberS{
				Value: expiresAt.Format(time.RFC3339),
			},
			":ttl": &types.AttributeValueMemberN{
				Value: strconv.FormatInt(expiresAt.Unix(), 10),
			},
		},
	})

	return err
}

func (r *SessionRepository) Revoke(ctx context.Context, sessionID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: sessionID,
			},
		},
	})

	return err
}
