package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BaseDynamoClient interface {
	GetItem(ctx context.Context, inputRequest GetItemParams, outputItem interface{}) (interface{}, error)
	PutItem(ctx context.Context, values interface{}) error
	UpdateItem(ctx context.Context, inputRequest *UpdateItemParams, outputItem interface{}) error
	DeleteItem(ctx context.Context, inputRequest DeleteItemParams) error
	QueryItems(ctx context.Context, queryItemsParams *QueryItemsParams, outputItem interface{}) ([]map[string]types.AttributeValue, int, error)
	BatchGetItem(ctx context.Context, inputRequest BatchGetItemParams, outputItem interface{}) (interface{}, error)
	BatchPutItem(ctx context.Context, inputItem []interface{}) ([]interface{}, error)
	TransactWriteItems(ctx context.Context, inputRequest []*TransactWriteInput) error
}
