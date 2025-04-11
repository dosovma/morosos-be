package domain

import (
	"context"
	"fmt"

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

	return agreement, nil
}

func (a *Agreement) CreateAgreement(ctx context.Context, agreement types.Agreement) (string, error) {
	agreement.ID = uuid.New().String()

	if err := a.agreementStore.AgreementPut(ctx, agreement); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return agreement.ID, nil
}

func (a *Agreement) SignAgreement(ctx context.Context, id string) error {
	//TODO fetch apartment devices and turn it on
	return nil
}
