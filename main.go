package main

// https://discord.com/api/oauth2/authorize?client_id=981172114579140699&permissions=274894883840&scope=bot

import (
	"diss-cord/app"
	"diss-cord/bot"
	"diss-cord/models/config"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	config := config.NewConfig()

	app := app.App{}
	app.Initialize(&config)
	app.SetRouters()
	go app.Serve(&wg)

	bot := bot.NewBot(config.Token)
	go bot.Start(&wg)

	wg.Wait()

}
