package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/store"
)

func main() {
	dynamodb := store.NewAgreementDynamoDBStore(context.TODO(), "agreements")
	agreementDomain := domain.NewAgreementDomain(dynamodb)
	handler := handlers.NewAgreementHandler(agreementDomain)

	lambda.Start(handler.CreateHandler)
}
