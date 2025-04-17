package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/store"
)

func main() {
	dynamodb := store.NewApartmentDynamoDBStore(context.TODO(), "apartments")
	apartmentDomain := domain.NewApartmentDomain(dynamodb, nil)
	handler := handlers.NewApartmentHandler(apartmentDomain)

	lambda.Start(handler.CreateHandler)
}
