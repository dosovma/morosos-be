package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

type AgreementDynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

var _ ports.AgreementStore = (*AgreementDynamoDBStore)(nil)

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

func (d *AgreementDynamoDBStore) AgreementGetByID(ctx context.Context, id string) (*entity.Agreement, error) {
	response, err := d.client.GetItem(
		ctx, &dynamodb.GetItemInput{
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

	agreement := entity.Agreement{}
	err = attributevalue.UnmarshalMap(response.Item, &agreement)

	if err != nil {
		return nil, fmt.Errorf("error getting item %w", err)
	}

	return &agreement, nil
}

func (d *AgreementDynamoDBStore) AgreementPut(ctx context.Context, agreement entity.Agreement) error {
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

func (d *AgreementDynamoDBStore) AgreementGetAllByStatus(ctx context.Context, status entity.AgreementStatus) ([]entity.Agreement, error) {
	const (
		limit    = 20
		timeZone = 2
	)

	agreementRange := make([]entity.Agreement, 0, limit)

	log.Printf("time to filter %v", time.Now().Add(timeZone*time.Hour).Format("2006-01-02T15:04"))

	input := &dynamodb.ScanInput{
		TableName: &d.tableName,
		Limit:     aws.Int32(limit),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":statusValue":    &ddbtypes.AttributeValueMemberS{Value: status},
			":elapsedAtValue": &ddbtypes.AttributeValueMemberS{Value: time.Now().Add(timeZone * time.Hour).Format("2006-01-02T15:04")},
		},
		ExpressionAttributeNames: map[string]string{
			"#status":     "status",
			"#elapsed_at": "elapsed_at",
		},
		FilterExpression: func(s string) *string { return &s }("#status = :statusValue and #elapsed_at <= :elapsedAtValue"),
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items from DynamoDB: %w", err)
	}

	if err = attributevalue.UnmarshalListOfMaps(result.Items, &agreementRange); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return agreementRange, nil
}
