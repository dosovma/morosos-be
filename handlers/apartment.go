package handlers

import (
	"context"
	"encoding/json"
	"github.com/dosovma/morosos-be/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/dosovma/morosos-be/domain"
)

type ApartmentHandler struct {
	apartment *domain.Apartment
}

func NewApartmentHandler(d *domain.Apartment) *ApartmentHandler {
	return &ApartmentHandler{
		apartment: d,
	}
}

func (l *ApartmentHandler) GetHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := event.QueryStringParameters["id"]
	if !ok {
		return errResponse(http.StatusBadRequest, "missing 'id' query parameter"), nil
	}

	//TODO change to 02361642-acf1-41db-a894-015b8db70e4d to get prepared apartment from db
	apartment, err := l.apartment.GetApartment(ctx, id)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if apartment == nil {
		return errResponse(http.StatusNotFound, "apartment not found"), nil
	} else {
		return response(http.StatusOK, apartment), nil
	}
}

func (l *ApartmentHandler) CreateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var apartment types.Apartment
	if err := json.Unmarshal([]byte(event.Body), &apartment); err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, err := l.apartment.CreateApartment(ctx, apartment)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if id == "" {
		return errResponse(http.StatusNotFound, "apartment not found"), nil
	} else {
		return response(http.StatusOK, id), nil
	}
}
