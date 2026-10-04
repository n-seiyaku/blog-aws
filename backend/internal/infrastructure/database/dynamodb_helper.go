package database

import (
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type postCursor struct {
	ID        string `json:"id"`
	FeedKey   string `json:"feedKey"`
	CreatedAt string `json:"createdAt"`
}

func encodePostCursor(
	key map[string]types.AttributeValue,
) (string, error) {
	id, ok := key["id"].(*types.AttributeValueMemberS)
	if !ok {
		return "", nil
	}

	feedKey, ok := key["feedKey"].(*types.AttributeValueMemberS)
	if !ok {
		return "", nil
	}

	createdAt, ok := key["createdAt"].(*types.AttributeValueMemberS)
	if !ok {
		return "", nil
	}

	cursor := postCursor{
		ID:        id.Value,
		FeedKey:   feedKey.Value,
		CreatedAt: createdAt.Value,
	}

	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodePostCursor(
	cursor string,
) (map[string]types.AttributeValue, error) {
	if cursor == "" {
		return nil, nil
	}

	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var value postCursor

	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: value.ID,
		},
		"feedKey": &types.AttributeValueMemberS{
			Value: value.FeedKey,
		},
		"createdAt": &types.AttributeValueMemberS{
			Value: value.CreatedAt,
		},
	}, nil
}
