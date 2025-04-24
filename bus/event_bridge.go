package bus

import (
	"context"
	"fmt"

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

func (b *EventBridgeBus) Publish(ctx context.Context, event ports.Event) error {
	/*
		Source:     "agreement",
		Detail:     agreement.Apartment.ID,
		DetailType: "sign",
		Resources:  nil,
	*/

	switch event.DetailType {
	case "sign":
		if err := b.apartment.SwitchDevices(ctx, event.Detail, true); err != nil {
			return err
		}

		return scheduleEvent(ctx)
	default:
		return fmt.Errorf("unknown event type ::: %s", event.DetailType)
	}
}

func scheduleEvent(ctx context.Context /*agreement entity.Agreement*/) error {
	// TODO use event bridge
	return nil
}
