package sms

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"

	"github.com/dosovma/morosos-be/ports"
)

const (
	AwsSenderID = "arn:aws:sms-voice:eu-north-1:539247474956:sender-id/5ELEMENTO/ES"
)

type AwsPinPointClient struct {
	client *pinpointsmsvoicev2.Client
}

var _ ports.SmsSender = (*AwsPinPointClient)(nil)

func NewSMSClient(ctx context.Context) *AwsPinPointClient {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config ::: %s", err)
	}

	return &AwsPinPointClient{
		client: pinpointsmsvoicev2.NewFromConfig(cfg),
	}
}

func (c *AwsPinPointClient) Send(ctx context.Context, phoneNumber string, message string) error {
	params := &pinpointsmsvoicev2.SendTextMessageInput{
		DestinationPhoneNumber: &phoneNumber,
		MessageBody:            &message,
		OriginationIdentity:    func(s string) *string { return &s }(AwsSenderID),
	}

	if _, err := c.client.SendTextMessage(ctx, params); err != nil {
		log.Printf("failed to send sms ::: %s", err)

		return err
	}

	return nil
}
