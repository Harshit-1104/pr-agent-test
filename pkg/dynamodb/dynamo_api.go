package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"time"


	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Keys struct {
	PartitionKey string `dynamodbav:"pk"`
	SortKey      string `dynamodbav:"sk"`
}

type GetItemParams struct {
	KeyCondition   Keys
	ConsistentRead bool
}

type QueryItemsParams struct {
	ProjectionKeys   []string
	FilterParamsList []ItemParams
	PartitionKey     string
	SortKey          ItemParams
	ConsistentRead   bool
}

type ItemParams struct {
	KeyName       string
	ConditionName string
	Value         interface{}
}

type DeleteItemParams struct {
	KeyCondition Keys
}

type BatchGetItemParams struct {
	KeyCondition   []Keys
	ProjectionKeys []string
	ConsistentRead bool
}

func (dynamoClient *DynamoClient) GetItem(ctx context.Context, inputRequest GetItemParams, outputItem interface{}) (interface{}, error) {
	row, err := marshalMap(inputRequest.KeyCondition)

	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName:      aws.String(dynamoClient.Table),
		Key:            row,
		ConsistentRead: &inputRequest.ConsistentRead,
	}
	result, err := dynamoClient.getItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if outputItem == nil {
		return result.Item, err
	}
	err = unmarshalMap(result.Item, outputItem)

	return outputItem, err
}

func (dynamoClient *DynamoClient) PutItem(ctx context.Context, values interface{}) error {
	row, err := marshalMap(values)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                row,
		TableName:           aws.String(dynamoClient.Table),
		ConditionExpression: aws.String("attribute_not_exists(partitionKey)"),
	}

	_, err = dynamoClient.putItem(ctx, input)

	return err
}

type UpdateItemParams struct {
	ValuesToBeUpdated     map[string]interface{}
	ValuesToBeRemoved     []string
	KeyCondition          Keys
	ValuesToBeIncremented []string
	ValuesToBeAppended    map[string]interface{}
	ConditionalParams     []ItemParams
}

func (dynamoClient *DynamoClient) UpdateItem(ctx context.Context, inputRequest *UpdateItemParams, outputItem interface{}) error {
	expr, err := makeUpdateExpression(inputRequest)

	if err != nil {
		return err
	}

	key, err := marshalMap(inputRequest.KeyCondition)

	if err != nil {
		return err
	}

	queryInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(dynamoClient.Table),
		Key:                       key,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		ReturnValues:              types.ReturnValueAllNew,
		ConditionExpression:       expr.Condition(),
	}

	res, err := dynamoClient.updateItem(ctx, queryInput)

	if err != nil {
		return err
	}

	if outputItem != nil {
		err = unmarshalMap(res.Attributes, outputItem)
	}

	return err
}

func makeUpdateExpression(inputRequest *UpdateItemParams) (expression.Expression, error) {
	update := expression.UpdateBuilder{}
	if inputRequest.ValuesToBeUpdated != nil {
		inputRequest.ValuesToBeUpdated["lastUpdated"] = time.Now().Unix()
	} else {
		inputRequest.ValuesToBeUpdated = map[string]interface{}{
			"lastUpdated": time.Now().Unix(),
		}
	}

	for key, val := range inputRequest.ValuesToBeUpdated {
		update = update.Set(
			name(key),
			value(val),
		)
	}

	for _, val := range inputRequest.ValuesToBeRemoved {
		update = update.Remove(
			name(val),
		)
	}

	for _, val := range inputRequest.ValuesToBeIncremented {
		update.Set(
			name(val),
			plus(val, 1),
		)
	}

	for key, val := range inputRequest.ValuesToBeAppended {
		update = update.Set(
			name(key),
			listAppend(key, val),
		)
	}

	conditions, err := setFilterConditionBuilder(inputRequest.ConditionalParams)
	if err != nil {
		return expression.Expression{}, err
	}

	expressionBuilder := expression.NewBuilder()
	expressionBuilder = expressionBuilder.WithUpdate(update)

	if conditions.IsSet() {
		expressionBuilder = expressionBuilder.WithCondition(conditions)
	}

	expr, err := expressionBuilder.Build()

	if err != nil {
		return expr, err
	}

	return expr, nil
}

func (dynamoClient *DynamoClient) DeleteItem(ctx context.Context, inputRequest DeleteItemParams) error {
	row, err := marshalMap(inputRequest.KeyCondition)

	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key:       row,
		TableName: aws.String(dynamoClient.Table),
	}

	_, err = dynamoClient.deleteItem(ctx, input)

	return err
}

func (dynamoClient *DynamoClient) QueryItems(
	ctx context.Context,
	queryItemsParams *QueryItemsParams,
	outputItem interface{},
) ([]map[string]types.AttributeValue, int, error) {
	projectionBuilder := setProjectionBuilder(queryItemsParams.ProjectionKeys)

	conditionBuilder, err := setConditionBuilder(queryItemsParams.PartitionKey, queryItemsParams.SortKey)
	if err != nil {
		return nil, 0, err
	}

	filterCondition, err := setFilterConditionBuilder(queryItemsParams.FilterParamsList)
	if err != nil {
		return nil, 0, err
	}

	expr := expression.NewBuilder()
	expr = expr.WithCondition(conditionBuilder)
	if len(queryItemsParams.FilterParamsList) > 0 {
		expr = expr.WithFilter(filterCondition)
	}
	if len(queryItemsParams.ProjectionKeys) > 0 {
		expr = expr.WithProjection(projectionBuilder)
	}
	build, err := expr.Build()

	if err != nil {
		return nil, 0, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(dynamoClient.Table),
		KeyConditionExpression:    build.Condition(),
		FilterExpression:          build.Filter(),
		ExpressionAttributeValues: build.Values(),
		ExpressionAttributeNames:  build.Names(),
		ProjectionExpression:      build.Projection(),
		ConsistentRead:            &queryItemsParams.ConsistentRead,
	}

	result, err := dynamoClient.queryItems(ctx, input)
	if err != nil {
		return nil, 0, err
	}

	if len(result.Items) == 0 {
		return result.Items, 0, nil
	}

	if outputItem == nil {
		return result.Items, len(result.Items), nil
	}

	err = unmarshalListOfMaps(result.Items, outputItem)
	if err != nil {
		return nil, 0, err
	}

	return result.Items, len(result.Items), err
}

func setFilterConditionBuilder(filterParamsList []ItemParams) (expression.ConditionBuilder, error) {
	var filterCondition expression.ConditionBuilder
	if len(filterParamsList) > 0 {
		keyName := filterParamsList[0].KeyName
		filterValue := filterParamsList[0].Value
		switch filterParamsList[0].ConditionName {
		case Equal:
			filterCondition = name(keyName).Equal(value(filterValue))
		case BeginsWith:
			filterCondition = name(keyName).BeginsWith(fmt.Sprintf("%v", filterValue))
		case GreaterThanEqual:
			filterCondition = name(keyName).GreaterThanEqual(value(filterValue))
		default:
			return expression.ConditionBuilder{}, errors.New(InvalidOperationErrorMessage)
		}
		if len(filterParamsList) > 1 {
			for _, val := range filterParamsList[1:] {
				switch val.ConditionName {
				case Equal:
					filterCondition = filterCondition.And(name(val.KeyName).Equal(value(val.Value)))
				case BeginsWith:
					filterCondition = filterCondition.And(name(val.KeyName).BeginsWith(fmt.Sprintf("%v", val.Value)))
				case GreaterThanEqual:
					filterCondition = filterCondition.And(name(val.KeyName).GreaterThanEqual(value(filterValue)))
				default:
					return expression.ConditionBuilder{}, errors.New(InvalidOperationErrorMessage)
				}
			}
		}
	}
	return filterCondition, nil
}

func setConditionBuilder(partitionKey string, sortKey ItemParams) (expression.ConditionBuilder, error) {
	var conditionBuilder = name(partitionKeyName).Equal(value(partitionKey))

	if len(sortKey.KeyName) == 0 {
		return conditionBuilder, nil
	}
	switch sortKey.ConditionName {
	case Equal:
		conditionBuilder = conditionBuilder.And(name(sortKeyName).Equal(value(sortKey.Value)))
	case BeginsWith:
		conditionBuilder = conditionBuilder.And(name(sortKeyName).BeginsWith(fmt.Sprintf("%v", sortKey.Value)))
	default:
		return expression.ConditionBuilder{}, errors.New(InvalidOperationErrorMessage)
	}

	return conditionBuilder, nil
}

func setProjectionBuilder(projectionKeys []string) expression.ProjectionBuilder {
	var projectionBuilder expression.ProjectionBuilder
	if len(projectionKeys) > 0 {
		projectionBuilder = expression.NamesList(name(projectionKeys[0]))
		if len(projectionKeys) > 1 {
			for _, val := range projectionKeys[1:] {
				projectionBuilder = projectionBuilder.AddNames(name(val))
			}
		}
	}
	return projectionBuilder
}

func (dynamoClient *DynamoClient) BatchGetItem(ctx context.Context, inputRequest BatchGetItemParams, outputItem interface{}) (interface{}, error) {
	rows, err := convertKeysToAttributeValues(inputRequest.KeyCondition)

	if err != nil {
		return nil, err
	}

	var projections *string
	var names map[string]string

	if len(inputRequest.ProjectionKeys) > 0 {
		projectionBuilder := setProjectionBuilder(inputRequest.ProjectionKeys)
		expr := expression.NewBuilder()
		expr = expr.WithProjection(projectionBuilder)
		build, buildErr := expr.Build()
		if buildErr != nil {
			return nil, buildErr
		}

		projections = build.Projection()
		names = build.Names()
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			dynamoClient.Table: {
				Keys:                     rows,
				ProjectionExpression:     projections,
				ExpressionAttributeNames: names,
				ConsistentRead:           &inputRequest.ConsistentRead,
			},
		},
	}

	result, err := dynamoClient.batchGetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.UnprocessedKeys) > 0 {
		return nil, fmt.Errorf("unable to process keys for batch get item, unprocessed keys: %+v", result.UnprocessedKeys)
	}

	err = attributevalue.UnmarshalListOfMaps(result.Responses[dynamoClient.Table], outputItem)
	if err != nil {
		return nil, err
	}
	return result.Responses[dynamoClient.Table], err
}

func (dynamoClient *DynamoClient) BatchPutItem(ctx context.Context, inputItem []interface{}) ([]interface{}, error) {
	rows, err := convertDataToWriteRequest(inputItem)

	if err != nil {
		return nil, err
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			dynamoClient.Table: rows,
		},
	}

	result, err := dynamoClient.batchWriteItem(ctx, input)

	if err != nil {
		return nil, err
	}

	if len(result.UnprocessedItems) > 0 {
		unprocessedRows, err := convertWriteRequestItemToInterface(input.RequestItems[dynamoClient.Table])
		if err != nil {
			return nil, err
		}
		return unprocessedRows, nil
	}

	return nil, nil
}
