package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TransactWriteInput struct {
	Put    interface{}
	Update UpdateItemParams
	Delete DeleteItemParams
}

func (dynamoClient *DynamoClient) TransactWriteItems(ctx context.Context, inputRequest []*TransactWriteInput) error {
	var transactionInputRequest []types.TransactWriteItem
	for _, req := range inputRequest {
		if req.Put != nil {
			row, err := marshalMap(req.Put)

			if err != nil {
				return err
			}
			putItemRequest := types.TransactWriteItem{
				Put: &types.Put{
					TableName:           aws.String(dynamoClient.Table),
					Item:                row,
					ConditionExpression: aws.String("attribute_not_exists(partitionKey) and attribute_not_exists(sortKey)"),
				},
			}
			transactionInputRequest = append(transactionInputRequest, putItemRequest)
		}

		if len(req.Update.ValuesToBeUpdated) != 0 ||
			len(req.Update.ValuesToBeRemoved) != 0 ||
			len(req.Update.ValuesToBeIncremented) != 0 {
			expr, err := makeUpdateExpression(&req.Update)

			if err != nil {
				return err
			}

			key, err := marshalMap(req.Update.KeyCondition)
			if err != nil {
				return err
			}

			updateItemRequest := types.TransactWriteItem{
				Update: &types.Update{
					TableName:                 aws.String(dynamoClient.Table),
					Key:                       key,
					UpdateExpression:          expr.Update(),
					ExpressionAttributeValues: expr.Values(),
					ExpressionAttributeNames:  expr.Names(),
					ConditionExpression:       expr.Condition(),
				},
			}

			transactionInputRequest = append(transactionInputRequest, updateItemRequest)
		}

		if len(req.Delete.KeyCondition.PartitionKey) != 0 && len(req.Delete.KeyCondition.SortKey) != 0 {
			keyDelete, err := marshalMap(req.Delete.KeyCondition)

			if err != nil {
				return err
			}

			deleteItemRequest := types.TransactWriteItem{
				Delete: &types.Delete{
					TableName: aws.String(dynamoClient.Table),
					Key:       keyDelete,
				},
			}
			transactionInputRequest = append(transactionInputRequest, deleteItemRequest)
		}
	}

	finalRequest := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactionInputRequest,
	}

	return dynamoClient.transactWriteItems(ctx, finalRequest)
}
