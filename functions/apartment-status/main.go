package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/store"
)

func main() {
	apartmentStore := store.NewApartmentDynamoDBStore(context.TODO(), "apartments")
	agreementDomain := domain.NewApartmentDomain(apartmentStore)
	handler := handlers.NewApartmentHandler(agreementDomain)

	lambda.Start(handler.StatusHandler)
}
