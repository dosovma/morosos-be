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
}

func NewAgreementDomain(agStore types.AgreementStore, apartStore types.ApartmentStore) *Agreement {
	return &Agreement{
		agreementStore: agStore,
		apartStore:     apartStore,
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

	for _, device := range apartment.Devices {
		// TODO
		// invoke tuya client to turnOn devices
		// use agreement elapsed_at to set up devices

		log.Printf("device name: %s", device.Name)
	}

	agreement.Status = types.Signed
	if err = a.agreementStore.AgreementPut(ctx, *agreement); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
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
