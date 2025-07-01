package handlers

import (
	"context"
	"net/http"

	"github.com/danilobml/go-lambda-dynamo/models"
	"github.com/gin-gonic/gin"
)

func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hi you have hit the test route",
	})
}

func GetAllPeopleHandler(c *gin.Context) {
	people, err := models.GetAllPeople(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch people"})
		return
	}

	c.JSON(http.StatusOK, people)
}
