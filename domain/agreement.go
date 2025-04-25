package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

type templateName = string

const (
	agreementTemplateName templateName = "agreement"
)

type Agreement struct {
	agreementStore ports.AgreementStore
	apartStore     ports.ApartmentStore
	templateStore  ports.TemplateStore
	bus            ports.Bus
	tyuaClient     ports.TuyaClient
	templater      ports.Templater
}

func NewAgreementDomain(agStore ports.AgreementStore, apartStore ports.ApartmentStore, tmplStore ports.TemplateStore, bus ports.Bus, client ports.TuyaClient, tmpl ports.Templater) *Agreement {
	return &Agreement{
		agreementStore: agStore,
		apartStore:     apartStore,
		templateStore:  tmplStore,
		bus:            bus,
		tyuaClient:     client,
		templater:      tmpl,
	}
}

func (a *Agreement) CreateAgreement(ctx context.Context, apartmentID string, agreement entity.Agreement) (string, error) {
	apartment, err := a.apartStore.ApartmentGet(ctx, apartmentID)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	agreement.ID = uuid.New().String()
	agreement.Status = entity.Draft
	agreement.Apartment = *apartment

	if err = a.agreementStore.AgreementPut(ctx, agreement); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return agreement.ID, nil
}

func (a *Agreement) GetAgreement(ctx context.Context, id string) (*entity.Agreement, error) {
	agreement, err := a.agreementStore.AgreementGet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	agreement.Text, err = a.buildAgreementText(ctx, *agreement)
	if err != nil {
		return nil, err
	}

	return agreement, nil
}

func (a *Agreement) SignAgreement(ctx context.Context, id string) error {
	agreement, err := a.agreementStore.AgreementGet(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	agreement.Status = entity.Signed
	agreement.Text, err = a.buildAgreementText(ctx, *agreement)
	if err != nil {
		return err
	}

	if err = a.agreementStore.AgreementPut(ctx, *agreement); err != nil {
		return err
	}

	event := ports.Event{
		Source:     "agreement",
		Detail:     agreement.Apartment.ID,
		DetailType: "sign",
		Resources:  nil,
	}

	return a.bus.Publish(ctx, event)
}

func (a *Agreement) buildAgreementText(ctx context.Context, agreement entity.Agreement) (string, error) {
	template, err := a.templateStore.TemplateGet(ctx, agreementTemplateName)
	if err != nil {
		return "", err
	}

	text, err := a.templater.FillTemplate(ctx, template, ports.ToAgreementText(agreement, agreement.Apartment))
	if err != nil {
		return "", err
	}

	return text, nil
}
