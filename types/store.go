package types

import (
	"context"
)

type AgreementStore interface {
	AgreementGet(context.Context, string) (*Agreement, error)
	AgreementPut(context.Context, Agreement) error
}

type ApartmentStore interface {
	ApartmentGet(context.Context, string) (*Apartment, error)
	ApartmentPut(context.Context, Apartment) error
}
