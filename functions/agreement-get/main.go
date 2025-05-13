package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dosovma/morosos-be/templater"

	"github.com/dosovma/morosos-be/bus"
	"github.com/dosovma/morosos-be/clients"
	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/sms"
	"github.com/dosovma/morosos-be/store"
)

func main() {
	agreementStore := store.NewAgreementDynamoDBStore(context.TODO(), "agreements")
	apartmentStore := store.NewApartmentDynamoDBStore(context.TODO(), "apartments")
	smsClient := sms.NewSMSClient(context.TODO())
	apartmentDomain := domain.NewApartmentDomain(apartmentStore, clients.NewTuyaClient(), smsClient)
	eventBridge := bus.NewEventBridgeBus(apartmentDomain)
	templateStore := store.NewTemplateDynamoDBStore(context.TODO(), "templates")
	htmlTemplater := templater.NewHtmlTemplater()
	agreementDomain := domain.NewAgreementDomain(agreementStore, apartmentStore, templateStore, eventBridge, clients.NewTuyaClient(), htmlTemplater)
	handler := handlers.NewAgreementHandler(agreementDomain)

	lambda.Start(handler.GetHandler)
}
