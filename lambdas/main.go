package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/danilobml/go-lambda-dynamo/models"
	"github.com/danilobml/go-lambda-dynamo/router"
)

var ginLambda *ginadapter.GinLambda

func init() {
	models.InitDynamo()

	r := router.SetupRouter()
	ginLambda = ginadapter.New(r)
}

func main() {
	lambda.Start(ginLambda.ProxyWithContext)
}
