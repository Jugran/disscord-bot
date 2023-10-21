package main

// https://discord.com/api/oauth2/authorize?client_id=981172114579140699&permissions=274894883840&scope=bot

import (
	"diss-cord/app"
	"diss-cord/bot"
	"diss-cord/models"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)

	quit := make(chan bool)
	config := models.NewConfig()

	app := app.App{}
	app.Initialize(&config)
	app.SetRouters()
	wg.Add(1)
	go app.Serve(&wg, quit)

	bot := bot.NewBot(config.Token)
	wg.Add(1)
	go bot.Start(&wg, quit)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	close(quit)
}
