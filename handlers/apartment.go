package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/domain/entity"
)

type ApartmentHandler struct {
	apartment *domain.Apartment
}

func NewApartmentHandler(d *domain.Apartment) *ApartmentHandler {
	return &ApartmentHandler{
		apartment: d,
	}
}

func (l *ApartmentHandler) CreateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var apartment entity.Apartment
	if err := json.Unmarshal([]byte(event.Body), &apartment); err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, err := l.apartment.CreateApartment(ctx, apartment)
	if err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if id == "" {
		return errProxyResponse(http.StatusNotFound, "apartment not found"), nil
	}

	return proxyResponse(
		http.StatusCreated, createApartmentResp{
			ApartmentID: id,
		},
	), nil
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
	}

	return proxyResponse(http.StatusOK, apartment), nil
}

func (l *ApartmentHandler) StatusHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var status statusApartmentReq
	if err := json.Unmarshal([]byte(event.Body), &status); err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, ok := event.PathParameters["apartment_id"]
	if !ok {
		return errProxyResponse(http.StatusBadRequest, "missing 'id' path parameter"), nil
	}

	switch status.Action {
	case entity.ApartmentOff:
		if err := l.apartment.SwitchDevices(ctx, id, false); err != nil {
			return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
		}
	case entity.ApartmentOn:
		if err := l.apartment.SwitchDevices(ctx, id, true); err != nil {
			return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
		}
	default:
		return errProxyResponse(http.StatusBadRequest, "invalid action type"), nil
	}

	return proxyResponse(http.StatusOK, statusApartmentResp{Success: true}), nil
}

func (l *ApartmentHandler) GetAllHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	next := event.QueryStringParameters["next"]

	apartmentRange, err := l.apartment.GetAllApartment(ctx, &next)
	if err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return proxyResponse(http.StatusOK, apartmentRange), nil
}
