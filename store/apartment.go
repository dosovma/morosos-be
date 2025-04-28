package store

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
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

func (d *ApartmentDynamoDBStore) ApartmentGetAll(ctx context.Context, next *string) (entity.ApartmentRange, error) {
	apartmentRange := entity.ApartmentRange{
		Apartments: []entity.Apartment{},
	}

	input := &dynamodb.ScanInput{
		TableName: &d.tableName,
		Limit:     aws.Int32(20),
	}

	if next != nil {
		input.ExclusiveStartKey = map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: *next},
		}
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return apartmentRange, fmt.Errorf("failed to get items from DynamoDB: %w", err)
	}

	if err = attributevalue.UnmarshalListOfMaps(result.Items, &apartmentRange.Apartments); err != nil {
		return apartmentRange, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	if len(result.LastEvaluatedKey) > 0 {
		if key, ok := result.LastEvaluatedKey["id"]; ok {
			nextKey := key.(*ddbtypes.AttributeValueMemberS).Value
			apartmentRange.Next = &nextKey
		}
	}

	return apartmentRange, nil
}
