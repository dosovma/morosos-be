package store

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

type TemplateDynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

var _ ports.TemplateStore = (*TemplateDynamoDBStore)(nil)

func NewTemplateDynamoDBStore(ctx context.Context, tableName string) *TemplateDynamoDBStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &TemplateDynamoDBStore{
		client:    client,
		tableName: tableName,
	}
}

func (d *TemplateDynamoDBStore) TemplateGet(ctx context.Context, name string) (string, error) {
	response, err := d.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: &d.tableName,
			Key: map[string]ddbtypes.AttributeValue{
				"name": &ddbtypes.AttributeValueMemberS{Value: name},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if len(response.Item) == 0 {
		return "", nil
	}

	template := entity.Template{}
	err = attributevalue.UnmarshalMap(response.Item, &template)

	if err != nil {
		return "", fmt.Errorf("error getting item %w", err)
	}

	return template.Text, nil
}
