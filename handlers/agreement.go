package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/types"
)

type APIGatewayV2Handler struct {
	agreement *domain.Agreement
}

func NewAPIGatewayV2Handler(d *domain.Agreement) *APIGatewayV2Handler {
	return &APIGatewayV2Handler{
		agreement: d,
	}
}

func (l *APIGatewayV2Handler) GetHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println(event)

	id, ok := event.QueryStringParameters["id"]
	if !ok {
		return errResponse(http.StatusBadRequest, "missing 'id' query parameter"), nil
	}

	log.Printf("reveiced id: %s", id)

	agreement, err := l.agreement.GetAgreement(ctx, id)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if agreement == nil {
		return errResponse(http.StatusNotFound, "agreement not found"), nil
	} else {
		return response(http.StatusOK, agreement), nil
	}
}

func (l *APIGatewayV2Handler) CreateHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("body ::: %s", event.Body)

	var agreement types.Agreement
	if err := json.Unmarshal([]byte(event.Body), &agreement); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)

		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, err := l.agreement.CreateAgreement(ctx, agreement)
	log.Printf("id ::: %s", id)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if id == "" {
		return errResponse(http.StatusNotFound, "product not found"), nil
	} else {
		return response(http.StatusOK, id), nil
	}
}
