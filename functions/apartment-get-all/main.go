package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dosovma/morosos-be/clients"
	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/sms"
	"github.com/dosovma/morosos-be/store"
)

func main() {
	apartmentStore := store.NewApartmentDynamoDBStore(context.TODO(), "apartments")
	smsClient := sms.NewSMSClient(context.TODO())
	apartmentDomain := domain.NewApartmentDomain(apartmentStore, clients.NewTuyaClient(), smsClient)
	handler := handlers.NewApartmentHandler(apartmentDomain)

	lambda.Start(handler.GetAllHandler)
}
