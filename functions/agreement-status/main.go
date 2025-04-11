package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/store"
)

func main() {
	agreementStore := store.NewAgreementDynamoDBStore(context.TODO(), "agreements")
	apartmentStore := store.NewApartmentDynamoDBStore(context.TODO(), "apartments")
	agreementDomain := domain.NewAgreementDomain(agreementStore, apartmentStore)
	handler := handlers.NewAgreementHandler(agreementDomain)

	lambda.Start(handler.StatusHandler)
}
