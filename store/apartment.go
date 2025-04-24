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

type ApartmentDynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

var _ ports.ApartmentStore = (*ApartmentDynamoDBStore)(nil)

func NewApartmentDynamoDBStore(ctx context.Context, tableName string) *ApartmentDynamoDBStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &ApartmentDynamoDBStore{
		client:    client,
		tableName: tableName,
	}
}

func (d *ApartmentDynamoDBStore) ApartmentGet(ctx context.Context, id string) (*entity.Apartment, error) {
	response, err := d.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: &d.tableName,
			Key: map[string]ddbtypes.AttributeValue{
				"id": &ddbtypes.AttributeValueMemberS{Value: id},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if len(response.Item) == 0 {
		return nil, nil
	}

	apartment := entity.Apartment{}
	err = attributevalue.UnmarshalMap(response.Item, &apartment)

	if err != nil {
		return nil, fmt.Errorf("error getting item %w", err)
	}

	return &apartment, nil
}

func (d *ApartmentDynamoDBStore) ApartmentPut(ctx context.Context, apartment entity.Apartment) error {
	item, err := attributevalue.MarshalMap(&apartment)
	if err != nil {
		return fmt.Errorf("unable to marshal apartment: %w", err)
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
