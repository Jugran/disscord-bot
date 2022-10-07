package app

import (
	"diss-cord/handlers"
	"diss-cord/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	PORT   uint64
	DB     *gorm.DB
}

func (a *App) Initialize(config *models.Config) {
	a.Router = gin.Default()
	a.PORT = config.Port

	models.ConnectDb()
	a.DB = models.DB

	a.Router.SetTrustedProxies(nil)

	// set db values here

	//  init discord bot instance here
}

func echoResponseHandler(c *gin.Context) {
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

func (a *App) SetRouters() {
	a.Router.POST("/echo", echoResponseHandler)

	a.Router.GET("/insult", handlers.FetchInsultHandler)
	a.Router.GET("/insult/:target", handlers.FetchInsult)
	a.Router.POST("/insult/add", handlers.AddInsultHandler)

	a.Router.PUT("/role/:name", handlers.AddRoleHandler)

	a.Router.GET("/user/all", handlers.FetchUsersHandler)
	a.Router.POST("/user/add", handlers.AddUserHandler)
	a.Router.PATCH("/user/role", handlers.AddUserRoleHandler)

}

func (a *App) Serve(wg *sync.WaitGroup) {
	port := fmt.Sprintf(":%d", a.PORT)
	a.Router.Run(port)
	wg.Done()
}
