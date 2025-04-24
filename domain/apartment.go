package domain

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

type Apartment struct {
	store      ports.ApartmentStore
	tyuaClient ports.TuyaClient
}

func NewApartmentDomain(s ports.ApartmentStore, t ports.TuyaClient) *Apartment {
	return &Apartment{
		store:      s,
		tyuaClient: t,
	}
}

func (a *Apartment) CreateApartment(ctx context.Context, apartment entity.Apartment) (string, error) {
	apartment.ID = uuid.New().String()

	if err := a.store.ApartmentPut(ctx, apartment); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return apartment.ID, nil
}

func (a *Apartment) GetApartment(ctx context.Context, id string) (*entity.Apartment, error) {
	apartment, err := a.store.ApartmentGet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return apartment, nil
}

func (a *Apartment) SwitchDevices(ctx context.Context, id string, isOn bool) error {
	apartment, err := a.store.ApartmentGet(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if apartment == nil {
		return errors.New("not found")
	}

	for i, device := range apartment.Devices {
		if err = a.tyuaClient.PostDevice(device.ID, isOn); err != nil {
			log.Printf("failed to turn on device ::: %s", device.ID)

			return err
		}

		apartment.Devices[i].IsOn = isOn
	}

	return a.store.ApartmentPut(ctx, *apartment)
}
