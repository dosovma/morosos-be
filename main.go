package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(json.RawMessage) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("got: %s", event),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
