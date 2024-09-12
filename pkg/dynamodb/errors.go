package dynamodb

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func IsConditionalCheckFailedError(err error) bool {
	var dynamoErr *types.ConditionalCheckFailedException
	if errors.As(err, &dynamoErr) {
		return true
	}

	return false
}
