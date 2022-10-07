package handlers

import (
	"diss-cord/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type addInsultSchema struct {
	Text     string   `json:"text" binding:"required"`
	Severity uint8    `json:"severity" binding:"required"`
	Roles    []string `json:"roles" binding:"required"`
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
	var insult models.Insult
	target := c.Param("target")

	var result *gorm.DB

	if len(target) != 0 {
		result = models.DB.Order("random()").Limit(1).Where("target = ?", target).Find(&insult)
	} else {
		result = models.DB.Order("random()").Limit(1).Find(&insult)
	}

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"insult": insult, "target": target})
}
