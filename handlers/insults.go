package handlers

import (
	"diss-cord/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type addInsultSchema struct {
	Text     string   `json:"text" binding:"required"`
	Severity uint8    `json:"severity" binding:"required"`
	Roles    []string `json:"roles" binding:"required"`
}

func GetInsultHandler(c *gin.Context) {
	insults, count := models.FindInsults()

	fmt.Printf("Fetched %d insult entries", count)

	c.JSON(http.StatusOK, gin.H{
		"insults": insults,
		"count":   count,
	})
}

func EchoResponseHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)

	// generic map type
	var jsonData map[string]interface{}

	if err != nil {
		// Handle error
		fmt.Println("Error", err)
		c.Status(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal([]byte(body), &jsonData)

	if err != nil {
		fmt.Println("Error", err)
		c.Status(http.StatusBadRequest)
		return
	}

	fmt.Println(jsonData)

	c.JSON(http.StatusAccepted, gin.H{
		"data": jsonData,
	})
}

func AddInsult(c *gin.Context) {
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

	models.AddInsult(&insult, body.Roles)

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
