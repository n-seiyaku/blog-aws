package database

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

func stringAttribute(name, value string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		name: &types.AttributeValueMemberS{
			Value: value,
		},
	}
}
