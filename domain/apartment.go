package domain

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/types"
)

type Apartment struct {
	store      types.ApartmentStore
	tyuaClient types.TuyaClient
}

func NewApartmentDomain(s types.ApartmentStore, t types.TuyaClient) *Apartment {
	return &Apartment{
		store:      s,
		tyuaClient: t,
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

func (a *Apartment) TurnOffDevices(ctx context.Context, id string) error {
	apartment, err := a.store.ApartmentGet(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if apartment == nil {
		return errors.New("not found")
	}

	for i, device := range apartment.Devices {
		// invoke tuya client to turnOff devices
		apartment.Devices[i].IsOn = false

		log.Printf("device name ::: %s", device.Name)

		if err = a.tyuaClient.PostDevice("vdevo174111102058365"); err != nil {
			log.Printf("failed to action device ::: %s", "vdevo174111102058365")
		}
	}

	log.Printf("apartment ::: %v", apartment)

	return a.store.ApartmentPut(ctx, *apartment)
}

func (a *Apartment) TurnOnDevices(ctx context.Context, id string) error {
	apartment, err := a.store.ApartmentGet(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if apartment == nil {
		return errors.New("not found")
	}

	for i, device := range apartment.Devices {
		// invoke tuya client to turnOff devices
		apartment.Devices[i].IsOn = true

		log.Printf("device name ::: %s", device.Name)
	}

	log.Printf("apartment ::: %v", apartment)

	return a.store.ApartmentPut(ctx, *apartment)
}
