package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func response(code int, object interface{}) events.APIGatewayProxyResponse {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		MultiValueHeaders: nil,
		Body:              string(marshalled),
		IsBase64Encoded:   false,
	}
}

func errResponse(status int, body string) events.APIGatewayProxyResponse {
	message := map[string]string{
		"message": body,
	}

	messageBytes, _ := json.Marshal(&message)

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(messageBytes),
	}
}
