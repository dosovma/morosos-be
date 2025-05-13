package ports

import (
	"context"

	"github.com/dosovma/morosos-be/domain/entity"
)

// Input

type Apartment interface {
	GetApartment(context.Context, string) (*entity.Apartment, error)
	CreateApartment(context.Context, entity.Apartment) (string, error)
	SwitchDevices(context.Context, string, bool) error
	GetAllApartment(context.Context, *string) (entity.ApartmentRange, error)
}

type Agreement interface {
	CreateAgreement(context.Context, string, entity.Agreement) (string, error)
	GetAgreement(context.Context, string) (*entity.Agreement, error)
	SignAgreement(context.Context, string) error
	CompleteAgreement(context.Context) error
}

// Output

type AgreementStore interface {
	AgreementGetByID(context.Context, string) (*entity.Agreement, error)
	AgreementPut(context.Context, entity.Agreement) error
	AgreementGetAllByStatus(context.Context, entity.AgreementStatus) ([]entity.Agreement, error)
}

type ApartmentStore interface {
	ApartmentGet(context.Context, string) (*entity.Apartment, error)
	ApartmentPut(context.Context, entity.Apartment) error
	ApartmentGetAll(context.Context, *string) (entity.ApartmentRange, error)
}

type TemplateStore interface {
	TemplateGet(context.Context, string) (string, error)
}

type TuyaClient interface {
	PostDevice(string, bool) error
}

type Bus interface {
	PublishAgreementEvent(context.Context, entity.Agreement) error
}

type Templater interface {
	FillTemplate(string, any) (string, error)
}

type SmsSender interface {
	Send(context.Context, string, string) error
}
