package domain

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/types"
)

type Agreement struct {
	agreementStore types.AgreementStore
	apartStore     types.ApartmentStore
	tyuaClient     types.TuyaClient
}

func NewAgreementDomain(agStore types.AgreementStore, apartStore types.ApartmentStore, client types.TuyaClient) *Agreement {
	return &Agreement{
		agreementStore: agStore,
		apartStore:     apartStore,
		tyuaClient:     client,
	}
}

func (a *Agreement) GetAgreement(ctx context.Context, id string) (*types.Agreement, error) {
	agreement, err := a.agreementStore.AgreementGet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	agreement.Text = agreementText(*agreement)

	return agreement, nil
}

func (a *Agreement) CreateAgreement(ctx context.Context, agreement types.Agreement) (string, error) {
	agreement.ID = uuid.New().String()
	agreement.Status = types.Draft

	if err := a.agreementStore.AgreementPut(ctx, agreement); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return agreement.ID, nil
}

func (a *Agreement) SignAgreement(ctx context.Context, id string) error {
	agreement, err := a.agreementStore.AgreementGet(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	apartment, err := a.apartStore.ApartmentGet(ctx, agreement.ApartmentID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for i, device := range apartment.Devices {
		if err = a.tyuaClient.PostDevice(device.ID, true); err != nil {
			log.Printf("failed to action device ::: %s", device.ID)

			return err
		}

		apartment.Devices[i].IsOn = true
	}

	if err = a.apartStore.ApartmentPut(ctx, *apartment); err != nil {
		log.Printf("failed to store apartmnet ::: %s", apartment.ID)

		return err
	}

	agreement.Status = types.Signed

	return a.agreementStore.AgreementPut(ctx, *agreement)
}

func agreementText(agreement types.Agreement) string {
	return fmt.Sprintf(
		"Я, имя %s фамилия %s, подписываю это соглашение, находясь в здравом уме. Дата старта соглашения: %s. Дата окончания соглашения: %s",
		agreement.Tenant.Name,
		agreement.Tenant.Surname,
		agreement.StartAt,
		agreement.ElapsedAt,
	)
}
