package clients

import (
	"log"

	tuyasdk "github.com/iot-eco-system/tuya-iot-service-sdk"
	"github.com/iot-eco-system/tuya-iot-service-sdk/model"

	"github.com/dosovma/morosos-be/ports"
)

type TuyaClient struct{}

var _ ports.TuyaClient = (*TuyaClient)(nil)

func NewTuyaClient() *TuyaClient {
	return &TuyaClient{}
}

func (*TuyaClient) PostDevice(id string, isOn bool) error {
	log.Println("tuya starts")

	client := tuyasdk.NewTuyaAPIClient(
		tuyasdk.NewTuyaAPIClientOptions{
			Host:     "https://openapi.tuyaeu.com",
			ClientID: "53jq8rhe8x8mgkvhfc3s",
			Secret:   "ac2b32f8563e448aab0281b7df4ab92f",
		},
	)

	client.Start()

	if err := client.GetToken(); err != nil {
		return err
	}

	req := &model.DeviceSendCommandRequest{
		DeviceID: id,
		Commands: buildCmd(id, isOn),
	}

	_, err := client.DeviceSendCommands(req)
	if err != nil {
		client.Stop()

		log.Printf("failed to process command ::: %s", err)

		return err
	}

	client.Stop()

	log.Println("tuya command executed")

	return nil
}

func buildCmd(deviceID string, isOn bool) []model.DeviceProperty {
	switch deviceID {
	case "vdevo174111102058365":
		return []model.DeviceProperty{
			{
				Code:  "switch_led_1",
				Value: isOn,
			},
		}
	case "bf4f1d68db7f487077qsfd":
		return []model.DeviceProperty{
			{
				Code:  "switch_1",
				Value: isOn,
			},
		}
	case "vdevo174489686258065":
		return []model.DeviceProperty{
			{
				Code:  "switch_1",
				Value: isOn,
			},
			{
				Code:  "switch_2",
				Value: isOn,
			},
			{
				Code:  "switch_3",
				Value: isOn,
			},
		}
	}

	return []model.DeviceProperty{}
}
