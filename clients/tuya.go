package clients

import (
	"context"
	"fmt"
	"log"

	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/httplib"

	"github.com/dosovma/morosos-be/types"
)

type TuyaClient struct{}

var _ types.TuyaClient = (*TuyaClient)(nil)

func NewTuyaClient() *TuyaClient {
	// custom init config
	connector.InitWithOptions(
		env.WithApiHost(httplib.URL_EU),
		env.WithMsgHost(httplib.MSG_EU),
		env.WithAccessID("53jq8rhe8x8mgkvhfc3s"),
		env.WithAccessKey("ac2b32f8563e448aab0281b7df4ab92f"),
		env.WithAppName("morosos"),
		env.WithDebugMode(true),
	)

	return &TuyaClient{}
}

func (t *TuyaClient) PostDevice(id string) error {
	//body, _ := ioutil.ReadAll(t.Request.Body)

	body := `
	{
    "commands":[
      {
        "code": "switch_led_1",
        "value": false
      },
      {
        "code": "bright_value_1",
        "value": 10
      },
      {
        "code": "brightness_min_1",
        "value": 50
      },
      {
        "code": "led_type_1",
        "value": "halogen"
      },
      {
        "code": "brightness_max_1",
        "value": 80
      },
      {
        "code": "countdown_1",
        "value": 5
      }
    ]
	}
`

	resp := &types.PostDeviceCmdResponse{}
	err := connector.MakePostRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s/commands", id)),
		connector.WithPayload([]byte(body)),
		connector.WithResp(resp),
	)
	if err != nil {
		return err
	}

	log.Printf("resp ::: %v", resp)

	return nil
}
