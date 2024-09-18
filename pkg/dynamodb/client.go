package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoClient struct {
	Client *dynamodb.Client
	Table  string
}

type Config struct {
	AwsAccessToken      string
	TableName           string
	TimeoutMilliseconds int
}

func NewClient(cfg *Config) (*DynamoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient := http.NewBuildableClient().WithTimeout(time.Millisecond * time.Duration(cfg.TimeoutMilliseconds))
	dynamoCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
		config.WithHTTPClient(httpClient),
		config.WithRetryMaxAttempts(DefaultRetry),
	)

	if err != nil {
		return nil, err
	}

	if cfg.EndpointNeeded {
		dynamoCfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           cfg.Endpoint,
					SigningRegion: cfg.Region,
				}, nil
			},
		)
	}

	if cfg.ConnectToDev {
		dynamoCfg.Credentials = credentials.NewStaticCredentialsProvider(cfg.AwsAccessID, cfg.AwsAccessKey, cfg.AwsAccessToken)
	}

	if cfg.MaxRetryAttempts > 1 {
		dynamoCfg.Retryer = func() aws.Retryer {
			return retry.NewStandard(func(o *retry.StandardOptions) {
				o.MaxAttempts = cfg.MaxRetryAttempts
				o.MaxBackoff = time.Duration(cfg.MaxRetryBackoff) * time.Millisecond
			})
		}
	}

	dynamoDBClient := &DynamoClient{
		Client: dynamodb.NewFromConfig(dynamoCfg),
		Table:  cfg.TableName,
	}

	return dynamoDBClient, nil
}

func (dynamoClient *DynamoClient) getItem(ctx context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	response, err := dynamoClient.Client.GetItem(ctx, input)
	return response, err
}

func (dynamoClient *DynamoClient) queryItems(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	response, err := dynamoClient.Client.Query(ctx, input)
	return response, err
}

func (dynamoClient *DynamoClient) putItem(ctx context.Context, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	response, err := dynamoClient.Client.PutItem(ctx, input)
	return response, err
}

func (dynamoClient *DynamoClient) deleteItem(ctx context.Context, input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	response, err := dynamoClient.Client.DeleteItem(ctx, input)
	return response, err
}

func (dynamoClient *DynamoClient) batchWriteItem(
	ctx context.Context,
	input *dynamodb.BatchWriteItemInput,
) (*dynamodb.BatchWriteItemOutput, error) {
	response, err := dynamoClient.Client.BatchWriteItem(ctx, input)
	return response, err
}

func marshalMap(item interface{}) (map[string]types.AttributeValue, error) {
	attributeValue, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}
	return attributeValue, nil
}

func name(name string) expression.NameBuilder {
	return expression.Name(name)
}

func value(value interface{}) expression.ValueBuilder {
	return expression.Value(value)
}

func plus(name string, value uint64) expression.SetValueBuilder {
	return expression.Plus(expression.Name(name).IfNotExists(expression.Value(0)), expression.Value(value))
}

func listAppend(name string, val interface{}) expression.SetValueBuilder {
	return expression.ListAppend(expression.Name(name).IfNotExists(expression.Value([]interface{}{})), expression.Value(val))
}

func unmarshalListOfMaps(value []map[string]types.AttributeValue, items interface{}) error {
	err := attributevalue.UnmarshalListOfMaps(value, items)
	if err != nil {
		return err
	}
	return nil
}

func unmarshalMap(value map[string]types.AttributeValue, item interface{}) error {
	err := attributevalue.UnmarshalMap(value, item)
	if err != nil {
		return err
	}
	return nil
}

func convertKeysToAttributeValues(keys []Keys) ([]map[string]types.AttributeValue, error) {
	items := make([]map[string]types.AttributeValue, len(keys))
	for i, v := range keys {
		val, err := marshalMap(v)
		if err != nil {
			return nil, err
		}
		items[i] = val
	}
	return items, nil
}

func convertDataToWriteRequest(itemsArr []interface{}) ([]types.WriteRequest, error) {
	items := make([]types.WriteRequest, len(itemsArr))
	for i, v := range itemsArr {
		val, err := marshalMap(v)
		if err != nil {
			return nil, err
		}
		items[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: val,
			},
		}
	}

	return items, nil
}

func convertWriteRequestItemToInterface(itemsArr []types.WriteRequest) ([]interface{}, error) {
	items := make([]interface{}, len(itemsArr))
	for i, v := range itemsArr {
		err := unmarshalMap(v.PutRequest.Item, &items[i])
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (dynamoClient *DynamoClient) transactWriteItems(ctx context.Context, input *dynamodb.TransactWriteItemsInput) error {
	_, err := dynamoClient.Client.TransactWriteItems(ctx, input)
	return err
}

func (dynamoClient *DynamoClient) updateItem(ctx context.Context, input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	response, err := dynamoClient.Client.UpdateItem(ctx, input)
	return response, err
}

func (dynamoClient *DynamoClient) batchGetItem(ctx context.Context, input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	response, err := dynamoClient.Client.BatchGetItem(ctx, input)
	return response, err
}
