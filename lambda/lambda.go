package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/blueimp/passphrase"
	"github.com/blueimp/passphrase/internal/parse"
)

const defaultNumber = 4
const maxNumber = 100

func logRequest(request *events.APIGatewayProxyRequest) {
	encodedRequest, err := json.Marshal(request)
	if err != nil {
		log.Println("Error:", err)
	} else {
		log.Println("Request:", string(encodedRequest))
	}
}

// Handler is the Lambda function handler:
func Handler(request *events.APIGatewayProxyRequest) (
	events.APIGatewayProxyResponse,
	error,
) {
	logRequest(request)
	number := parse.NaturalNumber(
		request.QueryStringParameters["n"],
		defaultNumber,
		maxNumber,
	)
	pass, err := passphrase.String(number)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"cache-control":             "private",
			"content-type":              "text/plain;charset=utf-8",
			"strict-transport-security": "max-age=31536000;includeSubDomains;preload",
			"x-content-type-options":    "nosniff",
		},
		Body: pass,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
