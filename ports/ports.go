package ports

import (
	"context"

	"github.com/dosovma/morosos-be/domain/entity"
)

// Input

type Apartment interface {
	GetApartment(ctx context.Context, id string) (*entity.Apartment, error)
	CreateApartment(ctx context.Context, apartment entity.Apartment) (string, error)
	SwitchDevices(ctx context.Context, id string, isOn bool) error
}

type Agreement interface {
	CreateAgreement(ctx context.Context, id string, agreement entity.Agreement) (string, error)
	GetAgreement(ctx context.Context, id string) (*entity.Agreement, error)
	SignAgreement(ctx context.Context, id string) error
}

// Output

type Store interface {
	AgreementStore
	ApartmentStore
}

type AgreementStore interface {
	AgreementGet(context.Context, string) (*entity.Agreement, error)
	AgreementPut(context.Context, entity.Agreement) error
}

type ApartmentStore interface {
	ApartmentGet(context.Context, string) (*entity.Apartment, error)
	ApartmentPut(context.Context, entity.Apartment) error
}

type TuyaClient interface {
	PostDevice(string, bool) error
}

type Bus interface {
	Publish(context.Context, Event) error
}
