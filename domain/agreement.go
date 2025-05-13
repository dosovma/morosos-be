package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

var _ ports.Agreement = (*Agreement)(nil)

type Agreement struct {
	agreementStore ports.AgreementStore
	apartmentStore ports.ApartmentStore
	templateStore  ports.TemplateStore
	bus            ports.Bus
	tyuaClient     ports.TuyaClient
	templater      ports.Templater
}

func NewAgreementDomain(agStore ports.AgreementStore, apartStore ports.ApartmentStore, tmplStore ports.TemplateStore, bus ports.Bus, client ports.TuyaClient, tmpl ports.Templater) *Agreement {
	return &Agreement{
		agreementStore: agStore,
		apartmentStore: apartStore,
		templateStore:  tmplStore,
		bus:            bus,
		tyuaClient:     client,
		templater:      tmpl,
	}
}

func (a *Agreement) CreateAgreement(ctx context.Context, apartmentID string, agreement entity.Agreement) (string, error) {
	apartment, err := a.apartmentStore.ApartmentGet(ctx, apartmentID)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	agreement.ID = uuid.New().String()
	agreement.Status = entity.Draft
	agreement.Apartment = entity.ApartmentData{
		ApartmentID: apartmentID,
		Landlord:    apartment.Landlord,
		Address:     apartment.Address,
	}

	agreement.Text, err = a.buildAgreementText(ctx, agreement)
	if err != nil {
		return "", err
	}

	if err = a.agreementStore.AgreementPut(ctx, agreement); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return agreement.ID, nil
}

func (a *Agreement) GetAgreement(ctx context.Context, id string) (*entity.Agreement, error) {
	agreement, err := a.agreementStore.AgreementGetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return agreement, nil
}

func (a *Agreement) SignAgreement(ctx context.Context, id string) error {
	agreement, err := a.agreementStore.AgreementGetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	agreement.Status = entity.Signed

	if err = a.agreementStore.AgreementPut(ctx, *agreement); err != nil {
		return err
	}

	return a.bus.PublishAgreementEvent(ctx, *agreement)
}

func (a *Agreement) buildAgreementText(ctx context.Context, agreement entity.Agreement) (string, error) {
	template, err := a.templateStore.TemplateGet(ctx, entity.AgreementTemplateName)
	if err != nil {
		return "", err
	}

	agreementText, err := a.templater.FillTemplate(template, entity.ToAgreementText(agreement))
	if err != nil {
		return "", err
	}

	return agreementText, nil
}

func (a *Agreement) CompleteAgreement(ctx context.Context) error {
	agreements, err := a.agreementStore.AgreementGetAllByStatus(ctx, entity.Signed)
	if err != nil {
		return err
	}

	for _, agreement := range agreements {
		agreement.Status = entity.Completed

		if err = a.agreementStore.AgreementPut(ctx, agreement); err != nil {
			return err
		}

		if err = a.bus.PublishAgreementEvent(ctx, agreement); err != nil {
			return err
		}
	}

	return nil
}
