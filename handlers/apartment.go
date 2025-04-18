package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/types"
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
	id, ok := event.PathParameters["apartment_id"]
	if !ok {
		return errProxyResponse(http.StatusBadRequest, "missing 'id' path parameter"), nil
	}

	apartment, err := l.apartment.GetApartment(ctx, id)
	if err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if apartment == nil {
		return errProxyResponse(http.StatusNotFound, "apartment not found"), nil
	} else {
		return proxyResponse(http.StatusOK, apartment), nil
	}
}

func (l *ApartmentHandler) CreateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var apartment types.Apartment
	if err := json.Unmarshal([]byte(event.Body), &apartment); err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, err := l.apartment.CreateApartment(ctx, apartment)
	if err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if id == "" {
		return errProxyResponse(http.StatusNotFound, "apartment not found"), nil
	} else {
		return proxyResponse(http.StatusOK, id), nil
	}
}

func (l *ApartmentHandler) StatusHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var status types.Status
	if err := json.Unmarshal([]byte(event.Body), &status); err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, ok := event.PathParameters["apartment_id"]
	if !ok {
		return errProxyResponse(http.StatusBadRequest, "missing 'id' path parameter"), nil
	}

	switch status.Action {
	case types.ApartmentOff:
		if err := l.apartment.TurnOffDevices(ctx, id); err != nil {
			return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
		}
	case types.ApartmentOn:
		if err := l.apartment.TurnOnDevices(ctx, id); err != nil {
			return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
		}
	default:
		return errProxyResponse(http.StatusBadRequest, "invalid action type"), nil
	}

	return proxyResponse(http.StatusOK, nil), nil

}
