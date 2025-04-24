package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

type Agreement struct {
	agreementStore ports.AgreementStore
	apartStore     ports.ApartmentStore
	bus            ports.Bus
	tyuaClient     ports.TuyaClient
}

func NewAgreementDomain(agStore ports.AgreementStore, apartStore ports.ApartmentStore, bus ports.Bus, client ports.TuyaClient) *Agreement {
	return &Agreement{
		agreementStore: agStore,
		apartStore:     apartStore,
		bus:            bus,
		tyuaClient:     client,
	}
}

func (a *Agreement) CreateAgreement(ctx context.Context, apartmentID string, agreement entity.Agreement) (string, error) {
	apartment, err := a.apartStore.ApartmentGet(ctx, apartmentID)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	agreement.ID = uuid.New().String()
	agreement.Status = entity.Draft
	agreement.Apartment = *apartment
	agreement.Text = buildAgreementText(agreement)

	if err = a.agreementStore.AgreementPut(ctx, agreement); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return agreement.ID, nil
}

func (a *Agreement) GetAgreement(ctx context.Context, id string) (*entity.Agreement, error) {
	agreement, err := a.agreementStore.AgreementGet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return agreement, nil
}

func (a *Agreement) SignAgreement(ctx context.Context, id string) error {
	agreement, err := a.agreementStore.AgreementGet(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	agreement.Status = entity.Signed

	if err = a.agreementStore.AgreementPut(ctx, *agreement); err != nil {
		return err
	}

	event := ports.Event{
		Source:     "agreement",
		Detail:     agreement.Apartment.ID,
		DetailType: "sign",
		Resources:  nil,
	}

	return a.bus.Publish(ctx, event)
}

func buildAgreementText(agreement entity.Agreement) string {
	return fmt.Sprintf(
		agreementText,
		agreement.Tenant.Name,
		agreement.Tenant.Surname,
		agreement.Apartment.Address,
		agreement.ElapsedAt,
		agreement.StartAt,
	)
}

var agreementText = `
Yo, %s %s, confirmo que he leído y acepto este acuerdo adicional relacionado con mi estancia en la vivienda ubicada en:\n
%s.\n
En esta vivienda se ha implementado un sistema que permite desconectar automáticamente los suministros de electricidad y agua al finalizar el período de alquiler.\n
Este sistema tiene como finalidad:\n
- evitar el consumo innecesario tras la salida del huésped,\n
- prevenir posibles incidentes (como fugas de agua o electrodomésticos encendidos),\n
- y facilitar la preparación del alojamiento para los próximos huéspedes.\n\n
Confirmo que:\n
✔️ Acepto que los suministros serán desconectados automáticamente al finalizar el contrato — %s.\n
✔️ Yo mismo/a he activado esta funcionalidad previamente, a través de la aplicación, antes del inicio de mi estancia.\n\n
Fecha de firma: %s
	`
