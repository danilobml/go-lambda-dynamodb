package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type App struct {
	Id string
}

type Person struct {
	Id      string `dynamodbav:"id" json:"id"`
	Name    string `dynamodbav:"name" json:"name"`
	Website string `dynamodbav:"website" json:"website"`
}

func newApp(id string) *App {
	return &App{
		Id: id,
	}
}

var dynamo *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}

	dynamo = dynamodb.NewFromConfig(cfg)
}

func GetAllPeople(ctx context.Context) ([]Person, error) {
	out, err := dynamo.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String("people"),
	})
	if err != nil {
		return nil, err
	}

	var people []Person
	err = attributevalue.UnmarshalListOfMaps(out.Items, &people)
	if err != nil {
		return nil, err
	}

	return people, nil
}

var INTERNAL_SERVER_ERROR_RESPONSE = events.APIGatewayProxyResponse{
	StatusCode: http.StatusInternalServerError,
	Headers:    map[string]string{"Content-Type": "application/json"},
	Body:       `{"error": "internal server error"}`,
}

func GenerateResponse(responseJson []byte) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                     "applicatuon/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "Content-Type",
			"Access-Control-Allow-Methods":     "OPTIONS, POST, GET, PUT, DELETE",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(responseJson),
	}
}

func (app *App) Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "GET" && request.Path == "test" {
		responseBody := map[string]string{
			"message": "Hi you have hit the test route",
		}

		body, err := json.Marshal(responseBody)
		if err != nil {
			return INTERNAL_SERVER_ERROR_RESPONSE, nil
		}

		response := GenerateResponse(body)
		return response, nil
	}

	if request.HTTPMethod == "GET" && request.Path == "/people" {
		people, err := GetAllPeople(context.TODO())
		if err != nil {
			return INTERNAL_SERVER_ERROR_RESPONSE, nil
		}

		body, err := json.Marshal(people)
		if err != nil {
			return INTERNAL_SERVER_ERROR_RESPONSE, nil
		}

		response := GenerateResponse(body)
		return response, nil
	}

	return INTERNAL_SERVER_ERROR_RESPONSE, nil
}

func main() {
	id := "id"

	app := newApp(id)

	lambda.Start(app.Handler)
}
