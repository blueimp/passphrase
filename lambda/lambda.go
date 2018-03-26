package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/blueimp/passphrase"
)

const defaultNumber = 4
const maxNumber = 100

func number(request *events.APIGatewayProxyRequest) int {
	parameter := request.QueryStringParameters["n"]
	if parameter == "" {
		return defaultNumber
	}
	number, err := strconv.Atoi(parameter)
	if err != nil || number > maxNumber {
		return maxNumber
	}
	return number
}

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
	pass, err := passphrase.Passphrase(number(request))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"content-type": "text/plain; charset=utf-8"},
		Body:       pass,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
