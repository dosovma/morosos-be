package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/dosovma/morosos-be/types"
)

type AgreementDynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

var _ types.AgreementStore = (*AgreementDynamoDBStore)(nil)

func NewAgreementDynamoDBStore(ctx context.Context, tableName string) *AgreementDynamoDBStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &AgreementDynamoDBStore{
		client:    client,
		tableName: tableName,
	}
}

func (d *AgreementDynamoDBStore) AgreementGet(ctx context.Context, id string) (*types.Agreement, error) {
	response, err := d.client.GetItem(
		ctx, &dynamodb.GetItemInput{
			TableName: &d.tableName,
			Key: map[string]ddbtypes.AttributeValue{
				"id": &ddbtypes.AttributeValueMemberS{Value: id},
			},
		},
	)

	log.Printf("get item sucess: %v", response)

	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if len(response.Item) == 0 {
		return nil, nil
	}

	agreement := types.Agreement{}
	err = attributevalue.UnmarshalMap(response.Item, &agreement)

	if err != nil {
		return nil, fmt.Errorf("error getting item %w", err)
	}

	return &agreement, nil
}

func (d *AgreementDynamoDBStore) AgreementPut(ctx context.Context, agreement types.Agreement) error {
	agreement.UpdatedAt = time.Now().String()

	item, err := attributevalue.MarshalMap(&agreement)
	if err != nil {
		return fmt.Errorf("unable to marshal agreement: %w", err)
	}

	_, err = d.client.PutItem(
		ctx, &dynamodb.PutItemInput{
			TableName: &d.tableName,
			Item:      item,
		},
	)

	if err != nil {
		return fmt.Errorf("cannot put item: %w", err)
	}

	return nil
}
