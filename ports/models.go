package ports

import "github.com/dosovma/morosos-be/domain/entity"

type Event struct {
	Source     string
	Detail     string
	DetailType string
	Resources  []string
}

type AgreementText struct {
	ElapsedAt        string
	TenantName       string
	TenantSurname    string
	ApartmentAddress string
}

func ToAgreementText(agreement entity.Agreement, apartment entity.Apartment) AgreementText {
	return AgreementText{
		ElapsedAt:        agreement.ElapsedAt,
		TenantName:       agreement.Tenant.Name,
		TenantSurname:    agreement.Tenant.Surname,
		ApartmentAddress: apartment.Address,
	}
}
