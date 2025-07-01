package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/danilobml/go-lambda-dynamo/handlers"
	"github.com/danilobml/go-lambda-dynamo/models"
)

var ginLambda *ginadapter.GinLambda

func init() {
	models.InitDynamo()

	router := gin.Default()

	router.GET("/test", handlers.TestHandler)
	router.GET("/people", handlers.GetAllPeopleHandler)

	ginLambda = ginadapter.New(router)
}

func main() {
	lambda.Start(ginLambda.ProxyWithContext)
}
