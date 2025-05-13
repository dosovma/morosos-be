package bus

import (
	"context"
	"fmt"

	"github.com/dosovma/morosos-be/domain/entity"
	"github.com/dosovma/morosos-be/ports"
)

type EventBridgeBus struct {
	apartment ports.Apartment
}

var _ ports.Bus = (*EventBridgeBus)(nil)

func NewEventBridgeBus(apartment ports.Apartment) *EventBridgeBus {
	return &EventBridgeBus{
		apartment: apartment,
	}
}

func (b *EventBridgeBus) PublishAgreementEvent(ctx context.Context, agreement entity.Agreement) error {
	switch agreement.Status {
	case entity.Signed:
		event := Event{
			Source:     Agreement,
			Detail:     agreement.Apartment.ApartmentID,
			DetailType: entity.Signed,
			Resources:  nil,
		}

		return b.apartment.SwitchDevices(ctx, event.Detail, true)
	case entity.Completed:
		event := Event{
			Source:     Agreement,
			Detail:     agreement.Apartment.ApartmentID,
			DetailType: entity.Completed,
			Resources:  nil,
		}

		return b.apartment.SwitchDevices(ctx, event.Detail, false)
	default:
		return fmt.Errorf("unknown agreement status ::: %s", agreement.Status)
	}
}
