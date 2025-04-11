package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/types"
)

type Agreement struct {
	store types.AgreementStore
}

func NewAgreementDomain(s types.AgreementStore) *Agreement {
	return &Agreement{
		store: s,
	}
}

func (a *Agreement) GetAgreement(ctx context.Context, id string) (*types.Agreement, error) {
	agreement, err := a.store.AgreementGet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return agreement, nil
}

func (a *Agreement) CreateAgreement(ctx context.Context, agreement types.Agreement) (string, error) {
	agreement.ID = uuid.New().String()

	if err := a.store.AgreementPut(ctx, agreement); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return agreement.ID, nil
}
