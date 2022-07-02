package main

// https://discord.com/api/oauth2/authorize?client_id=981172114579140699&permissions=274894883840&scope=bot

import (
	"diss-cord/app"
	"diss-cord/models"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	config := models.NewConfig()

	app := app.App{}
	app.Initialize(&config)
	app.SetRouters()
	go app.Serve(&wg)

	// bot := bot.NewBot(config.Token)
	// go bot.Start(&wg)

	wg.Wait()

}
