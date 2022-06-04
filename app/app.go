package app

import (
	"diss-cord/handlers"
	"diss-cord/models/config"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	PORT   uint64
}

func (a *App) Initialize(config *config.Config) {
	a.Router = gin.Default()
	a.PORT = config.Port

	a.Router.SetTrustedProxies(nil)

	// set db values here

	//  init discord bot instance here
}

func (a *App) SetRouters() {
	a.Router.GET("/insult", handlers.GetInsultHandler)
	a.Router.POST("/echo", handlers.EchoResponseHandler)
}

func (a *App) Serve(wg *sync.WaitGroup) {
	port := fmt.Sprintf(":%d", a.PORT)
	a.Router.Run(port)
	wg.Done()
}
