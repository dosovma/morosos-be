package entity

import (
	"log"
	"time"
)

type AgreementStatus = string

const (
	Draft     AgreementStatus = "draft"
	Signed    AgreementStatus = "signed"
	Completed AgreementStatus = "completed"
)

type AgreementAction = string

const (
	Sign AgreementAction = "sign"
)

type Agreement struct {
	ID        string          `dynamodbav:"id" json:"id"`
	StartAt   string          `dynamodbav:"start_at" json:"start_at"`
	ElapsedAt string          `dynamodbav:"elapsed_at" json:"elapsed_at"`
	UpdatedAt string          `dynamodbav:"updated_at" json:"updated_at"`
	Tenant    User            `dynamodbav:"tenant" json:"tenant"`
	Apartment ApartmentData   `dynamodbav:"apartment" json:"apartment"`
	Status    AgreementStatus `dynamodbav:"status" json:"status"`
	Text      string          `dynamodbav:"text" json:"text"` // TODO move out from entity to service layer
}

type ApartmentData struct {
	ApartmentID string `dynamodbav:"apartment_id" json:"apartment_id"`
	Landlord    User   `dynamodbav:"landlord" json:"landlord"`
	Address     string `dynamodbav:"address" json:"address"`
}

type AgreementText struct {
	ElapsedAt        string
	TenantName       string
	TenantSurname    string
	ApartmentAddress string
}

func ToAgreementText(agreement Agreement) AgreementText {
	elapsedTime, err := time.Parse("2006-01-02T15:04", agreement.ElapsedAt)
	if err != nil {
		log.Printf("failed to parse date ::: %s", err)
	}

	if !elapsedTime.IsZero() {
		agreement.ElapsedAt = elapsedTime.Format("02-01-2006 15:04")
	}

	return AgreementText{
		ElapsedAt:        agreement.ElapsedAt,
		TenantName:       agreement.Tenant.Name,
		TenantSurname:    agreement.Tenant.Surname,
		ApartmentAddress: agreement.Apartment.Address,
	}
}
