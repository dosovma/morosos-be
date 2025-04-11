package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/types"
)

type AgreementHandler struct {
	agreement *domain.Agreement
}

func NewAgreementHandler(d *domain.Agreement) *AgreementHandler {
	return &AgreementHandler{
		agreement: d,
	}
}

func (l *AgreementHandler) GetHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := event.QueryStringParameters["id"]
	if !ok {
		return errResponse(http.StatusBadRequest, "missing 'id' query parameter"), nil
	}

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

func (l *AgreementHandler) CreateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var agreement types.Agreement
	if err := json.Unmarshal([]byte(event.Body), &agreement); err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, err := l.agreement.CreateAgreement(ctx, agreement)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if id == "" {
		return errResponse(http.StatusNotFound, "agreement not found"), nil
	} else {
		return response(http.StatusOK, id), nil
	}
}

func (l *AgreementHandler) StatusHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//TODO check action start process sign

	if err := l.agreement.SignAgreement(ctx, "id"); err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return response(http.StatusOK, nil), nil

}
