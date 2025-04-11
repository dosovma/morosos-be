package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/types"
)

type Apartment struct {
	store types.ApartmentStore
}

func NewApartmentDomain(s types.ApartmentStore) *Apartment {
	return &Apartment{
		store: s,
	}
}

func (a *Apartment) GetApartment(ctx context.Context, id string) (*types.Apartment, error) {
	apartment, err := a.store.ApartmentGet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return apartment, nil
}

func (a *Apartment) CreateApartment(ctx context.Context, apartment types.Apartment) (string, error) {
	apartment.ID = uuid.New().String()

	if err := a.store.ApartmentPut(ctx, apartment); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return apartment.ID, nil
}
