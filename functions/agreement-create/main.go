package main

import (
	"context"
	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/store"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	tableName, ok := os.LookupEnv("TABLE")
	if !ok {
		panic("Need TABLE environment variable")
	}

	dynamodb := store.NewDynamoDBStore(context.TODO(), tableName)
	domain := domain.NewAgreementDomain(dynamodb)
	handler := handlers.NewAPIGatewayV2Handler(domain)
	lambda.Start(handler.CreateHandler)
}
