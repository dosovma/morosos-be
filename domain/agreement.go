package domain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

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
	elapsedTime, err := time.Parse("2006-01-02T15:04", agreement.ElapsedAt)
	if err != nil {
		log.Printf("failed to parse date ::: %s", err)
	}

	if !elapsedTime.IsZero() {
		agreement.ElapsedAt = elapsedTime.Format("02-01-2006 15:04")
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
	agreement, err := a.agreementStore.AgreementGet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return agreement, nil
}

func (a *Agreement) SignAgreement(ctx context.Context, id string) error {
	agreement, err := a.agreementStore.AgreementGet(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	agreement.Status = entity.Signed

	if err = a.agreementStore.AgreementPut(ctx, *agreement); err != nil {
		return err
	}

	event := ports.Event{
		Source:     "agreement",
		Detail:     agreement.Apartment.ApartmentID,
		DetailType: "sign",
		Resources:  nil,
	}

	return a.bus.Publish(ctx, event)
}

func (a *Agreement) buildAgreementText(ctx context.Context, agreement entity.Agreement) (string, error) {
	template, err := a.templateStore.TemplateGet(ctx, entity.AgreementTemplateName)
	if err != nil {
		return "", err
	}

	agreementText, err := a.templater.FillTemplate(ctx, template, ports.ToAgreementText(agreement))
	if err != nil {
		return "", err
	}

	return agreementText, nil
}
