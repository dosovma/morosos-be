package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

var _ ports.Apartment = (*Apartment)(nil)

type Apartment struct {
	store      ports.ApartmentStore
	tyuaClient ports.TuyaClient
	smsClient  ports.SmsSender
}

func NewApartmentDomain(s ports.ApartmentStore, t ports.TuyaClient, smsClient ports.SmsSender) *Apartment {
	return &Apartment{
		store:      s,
		tyuaClient: t,
		smsClient:  smsClient,
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
		if err = a.switchDevice(ctx, device, isOn); err != nil {
			return err
		}

		apartment.Devices[i].IsOn = isOn

		log.Println("device status changed")
	}

	return a.store.ApartmentPut(ctx, *apartment)
}

func (a *Apartment) GetAllApartment(ctx context.Context, next *string) (entity.ApartmentRange, error) {
	if next != nil && strings.TrimSpace(*next) == "" {
		next = nil
	}

	apartmentRange, err := a.store.ApartmentGetAll(ctx, next)
	if err != nil {
		return apartmentRange, fmt.Errorf("%w", err)
	}

	return apartmentRange, nil
}

func (a *Apartment) switchDevice(ctx context.Context, device entity.Device, isOn bool) error {
	if device.PhoneNumber != nil && *device.PhoneNumber != "" {
		message := entity.ApartmentOff
		if isOn {
			message = entity.ApartmentOn
		}

		if err := a.smsClient.Send(ctx, *device.PhoneNumber, message); err != nil {
			log.Printf("failed to switch device by sms ::: %s", device.ID)

			return err
		}

		return nil
	}

	if err := a.tyuaClient.PostDevice(device.ID, isOn); err != nil {
		log.Printf("failed to switch device ::: %s", device.ID)

		return err
	}

	return nil
}
