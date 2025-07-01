package handlers

import (
	"context"
	"net/http"

	"github.com/danilobml/go-lambda-dynamo/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func GetOnePersonByIDHandler(c *gin.Context) {
	id := c.Param("id")

	person, err := models.GetOnePersonById(context.TODO(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch person"})
	}

	if person == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
	}

	c.JSON(http.StatusOK, person)
}

func CreatePersonHandler(c *gin.Context) {
	var person models.Person

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	person.Id = uuid.New().String()

	if err := models.CreatePerson(context.TODO(), person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save person"})
	}

	c.JSON(http.StatusCreated, person)
}

func UpdatePersonHandler(c *gin.Context) {
	id := c.Param("id")
	var person models.Person

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if person.Name == "" || person.Website == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both name and website fields must be supplied"})
		return
	}

	person.Id = id

	if err := models.UpdatePerson(context.TODO(), person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update person"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func DeletePersonHandler(c *gin.Context) {
	id := c.Param("id")

	if err := models.DeletePerson(context.TODO(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "deletion failed"})
	}

	c.JSON(http.StatusNoContent, nil)
}
