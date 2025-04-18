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
	id, ok := event.PathParameters["agreement_id"]
	if !ok {
		return errProxyResponse(http.StatusBadRequest, "missing 'id' path parameter"), nil
	}

	agreement, err := l.agreement.GetAgreement(ctx, id)
	if err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if agreement == nil {
		return errProxyResponse(http.StatusNotFound, "agreement not found"), nil
	} else {
		return proxyResponse(http.StatusOK, agreement), nil
	}
}

func (l *AgreementHandler) CreateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var agreement types.Agreement
	if err := json.Unmarshal([]byte(event.Body), &agreement); err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, err := l.agreement.CreateAgreement(ctx, agreement)
	if err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if id == "" {
		return errProxyResponse(http.StatusNotFound, "agreement not found"), nil
	} else {
		return proxyResponse(http.StatusOK, id), nil
	}
}

func (l *AgreementHandler) StatusHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var status types.Status
	if err := json.Unmarshal([]byte(event.Body), &status); err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, ok := event.PathParameters["agreement_id"]
	if !ok {
		return errProxyResponse(http.StatusBadRequest, "missing 'id' path parameter"), nil
	}

	switch status.Action {
	case types.Sign:
		if err := l.agreement.SignAgreement(ctx, id); err != nil {
			return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
		}
	default:
		return errProxyResponse(http.StatusBadRequest, "invalid action type"), nil
	}

	return proxyResponse(http.StatusOK, nil), nil

}
