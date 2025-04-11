package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/types"
)

type Agreement struct {
	store types.Store
}

func NewAgreementDomain(s types.Store) *Agreement {
	return &Agreement{
		store: s,
	}
}

func (a *Agreement) GetAgreement(ctx context.Context, id string) (*types.Agreement, error) {
	agreement, err := a.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return agreement, nil
}

func (a *Agreement) CreateAgreement(ctx context.Context, agreement types.Agreement) (string, error) {
	agreement.Id = uuid.New().String()

	err := a.store.Put(ctx, agreement)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return agreement.Id, nil
}
