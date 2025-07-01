package models

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func GetOnePersonById(ctx context.Context, id string) (*Person, error) {
	out, err := Dynamo.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("people"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, nil
	}

	var person Person
	err = attributevalue.UnmarshalMap(out.Item, &person)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

func CreatePerson(ctx context.Context, p Person) error {
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}

	_, err = Dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("people"),
		Item:      item,
	})

	return err
}

func UpdatePerson(ctx context.Context, p Person) error {
	_, err := Dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("people"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: p.Id},
		},

		UpdateExpression: aws.String("SET #name = :name, #website = :website"),
		ExpressionAttributeNames: map[string]string{
			"#name":    "name",
			"#website": "website",
		},

		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":    &types.AttributeValueMemberS{Value: p.Name},
			":website": &types.AttributeValueMemberS{Value: p.Website},
		},
	})

	return err
}

func DeletePerson(ctx context.Context, id string) error {
	_, err := Dynamo.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("people"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	return err
}
