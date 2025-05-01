package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/dosovma/morosos-be/domain"
	"github.com/dosovma/morosos-be/domain/entity"
)

type AgreementHandler struct {
	agreement *domain.Agreement
}

func NewAgreementHandler(d *domain.Agreement) *AgreementHandler {
	return &AgreementHandler{
		agreement: d,
	}
}

func (l *AgreementHandler) CreateHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apartmentID, ok := event.PathParameters["apartment_id"]
	if !ok {
		return errProxyResponse(http.StatusBadRequest, "missing 'id' path parameter"), nil
	}

	var agreement entity.Agreement
	if err := json.Unmarshal([]byte(event.Body), &agreement); err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if err := validate(agreement); err != nil {
		return errProxyResponse(http.StatusBadRequest, fmt.Sprintf("validation error: %s", err)), nil
	}

	id, err := l.agreement.CreateAgreement(ctx, apartmentID, agreement)
	if err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	if id == "" {
		return errProxyResponse(http.StatusNotFound, "agreement not found"), nil
	} else {
		return proxyResponse(
			http.StatusOK, createAgreementResp{
				AgreementID: id,
			},
		), nil
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

func (l *AgreementHandler) StatusHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var status statusAgreementReq
	if err := json.Unmarshal([]byte(event.Body), &status); err != nil {
		return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	id, ok := event.PathParameters["agreement_id"]
	if !ok {
		return errProxyResponse(http.StatusBadRequest, "missing 'id' path parameter"), nil
	}

	switch status.Action {
	case entity.Sign:
		if err := l.agreement.SignAgreement(ctx, id); err != nil {
			return errProxyResponse(http.StatusInternalServerError, err.Error()), nil
		}
	default:
		return errProxyResponse(http.StatusBadRequest, "invalid action type"), nil
	}

	return proxyResponse(http.StatusOK, statusAgreementResp{Success: true}), nil
}

func (l *AgreementHandler) CompleteHandler(ctx context.Context) error {
	return l.agreement.CompleteAgreement(ctx)
}

func validate(agreement entity.Agreement) error {
	// TODO
	return nil
}
