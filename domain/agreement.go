package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/dosovma/morosos-be/types"
)

var (
	ErrJsonUnmarshal     = errors.New("failed to parse product from request body")
	ErrProductIdMismatch = errors.New("product ID in path does not match product ID in body")
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
	product, err := a.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return product, nil
}

func (a *Agreement) CreateAgreement(ctx context.Context, id string) (*types.Agreement, error) {
	product, err := a.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return product, nil
}
