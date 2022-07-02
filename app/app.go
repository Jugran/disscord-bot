package app

import (
	"diss-cord/handlers"
	"diss-cord/models"
	"fmt"
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

func (a *App) SetRouters() {
	a.Router.GET("/insult/:target", handlers.FetchInsult)
	a.Router.GET("/insult", handlers.FetchInsult)
	a.Router.POST("/echo", handlers.EchoResponseHandler)
	a.Router.POST("/insult/add", handlers.AddInsult)
}

func (a *App) Serve(wg *sync.WaitGroup) {
	port := fmt.Sprintf(":%d", a.PORT)
	a.Router.Run(port)
	wg.Done()
}
