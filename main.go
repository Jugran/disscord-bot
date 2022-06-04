package main

// https://discord.com/api/oauth2/authorize?client_id=981172114579140699&permissions=274894883840&scope=bot

import (
	"diss-cord/bot"
	"diss-cord/models/config"
)

func main() {

	// app := app.App{}
	// config := config.NewConfig()
	// app.Initialize(&config)
	// app.SetRouters()
	// app.Serve()

	config := config.NewConfig()

	bot := bot.NewBot(config.Token)
	bot.Start()

}
