package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/store"
)

func main() {
	dynamodb := store.NewDynamoDBStore(context.TODO(), "agreements")
	agreementDomain := domain.NewAgreementDomain(dynamodb)
	handler := handlers.NewAPIGatewayV2Handler(agreementDomain)

	lambda.Start(handler.CreateHandler)
}
