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
	GetAllApartment(ctx context.Context, next *string) (entity.ApartmentRange, error)
}

type Agreement interface {
	CreateAgreement(ctx context.Context, id string, agreement entity.Agreement) (string, error)
	GetAgreement(ctx context.Context, id string) (*entity.Agreement, error)
	SignAgreement(ctx context.Context, id string) error
}

// Output

type AgreementStore interface {
	AgreementGet(context.Context, string) (*entity.Agreement, error)
	AgreementPut(context.Context, entity.Agreement) error
}

type ApartmentStore interface {
	ApartmentGet(context.Context, string) (*entity.Apartment, error)
	ApartmentPut(context.Context, entity.Apartment) error
	ApartmentGetAll(ctx context.Context, next *string) (entity.ApartmentRange, error)
}

type TemplateStore interface {
	TemplateGet(ctx context.Context, name string) (string, error)
}

type TuyaClient interface {
	PostDevice(string, bool) error
}

type Bus interface {
	Publish(context.Context, Event) error
}

type Templater interface {
	FillTemplate(context.Context, string, any) (string, error)
}
