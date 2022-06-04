package handlers

import (
	"diss-cord/models/insults"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInsultHandler(c *gin.Context) {
	insultsData := insults.NewInsults()
	insult := insultsData.GetInsult()

	c.JSON(http.StatusOK, gin.H{
		"insult": insult,
	})
}

func EchoResponseHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		// Handle error
		fmt.Println("Error", err)
		c.Status(http.StatusBadRequest)
		return
	}
	jsonData := string(body)

	fmt.Println(jsonData)

	c.JSON(http.StatusAccepted, gin.H{
		"data": jsonData,
	})
}
