package handlers

import (
	"diss-cord/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addInsultSchema struct {
	Text     string   `json:"text" binding:"required"`
	Severity uint8    `json:"severity" binding:"required"`
	Roles    []string `json:"roles" binding:"required"`
}

type targetSchema struct {
	ID *uint `uri:"target"` // user id of the target user
}

func FetchInsultHandler(c *gin.Context) {
	insults, count := models.FindInsultsAction()

	fmt.Printf("Fetched %d insult entries", count)

	c.JSON(http.StatusOK, gin.H{
		"insults": insults,
		"count":   count,
	})
}

func AddInsultHandler(c *gin.Context) {
	var body addInsultSchema
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid Input",
		})
		return
	}

	insult := models.Insult{
		Text:     body.Text,
		Severity: body.Severity,
	}

	models.AddInsultAction(&insult, &body.Roles)

	c.JSON(http.StatusCreated, &insult)
}

func FetchInsult(c *gin.Context) {
	var target targetSchema

	if err := c.ShouldBindUri(&target); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid Target",
		})
		return
	}

	user_id := target.ID

	result, insult := models.GetInsultForUser(user_id)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"insults": insult.Text, "result": result.RowsAffected})
}
