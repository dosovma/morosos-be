package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dosovma/morosos-be/bus"
	"github.com/dosovma/morosos-be/clients"
	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/handlers"
	"github.com/dosovma/morosos-be/store"
	"github.com/dosovma/morosos-be/templater"
)

func main() {
	agreementStore := store.NewAgreementDynamoDBStore(context.TODO(), "agreements")
	apartmentStore := store.NewApartmentDynamoDBStore(context.TODO(), "apartments")
	templateStore := store.NewTemplateDynamoDBStore(context.TODO(), "templates")
	htmlTemplater := templater.NewHtmlTemplater()
	apartmentDomain := domain.NewApartmentDomain(apartmentStore, clients.NewTuyaClient())
	eventBridge := bus.NewEventBridgeBus(apartmentDomain)

	agreementDomain := domain.NewAgreementDomain(agreementStore, apartmentStore, templateStore, eventBridge, clients.NewTuyaClient(), htmlTemplater)
	handler := handlers.NewAgreementHandler(agreementDomain)

	lambda.Start(handler.StatusHandler)
}
