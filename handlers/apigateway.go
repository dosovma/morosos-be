package handlers

import (
	"context"
	"encoding/json"
	"github.com/dosovma/morosos-be/domain"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type APIGatewayV2Handler struct {
	agreement *domain.Agreement
}

func NewAPIGatewayV2Handler(d *domain.Agreement) *APIGatewayV2Handler {
	return &APIGatewayV2Handler{
		agreement: d,
	}
}

func (l *APIGatewayV2Handler) GetHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok {
		return errResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	product, err := l.agreement.GetAgreement(ctx, id)

	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if product == nil {
		return errResponse(http.StatusNotFound, "product not found"), nil
	} else {
		return response(http.StatusOK, product), nil
	}
}

func (l *APIGatewayV2Handler) CreateHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok {
		return errResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	product, err := l.agreement.CreateAgreement(ctx, id)

	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if product == nil {
		return errResponse(http.StatusNotFound, "product not found"), nil
	} else {
		return response(http.StatusOK, product), nil
	}
}

func response(code int, object interface{}) events.APIGatewayV2HTTPResponse {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(marshalled),
		IsBase64Encoded: false,
	}
}

func errResponse(status int, body string) events.APIGatewayV2HTTPResponse {
	message := map[string]string{
		"message": body,
	}

	messageBytes, _ := json.Marshal(&message)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(messageBytes),
	}
}
