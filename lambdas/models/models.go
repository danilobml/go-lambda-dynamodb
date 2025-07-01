package models

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var Dynamo *dynamodb.Client

func InitDynamo() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	Dynamo = dynamodb.NewFromConfig(cfg)
}

type Person struct {
	Id      string `dynamodbav:"id" json:"id"`
	Name    string `dynamodbav:"name" json:"name"`
	Website string `dynamodbav:"website" json:"website"`
}

func GetAllPeople(ctx context.Context) ([]Person, error) {
	out, err := Dynamo.Scan(ctx, &dynamodb.ScanInput{
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
