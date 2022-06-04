package bot

import (
	"diss-cord/models/insults"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Session *discordgo.Session
	Insults *insults.InsultData
}

func NewBot(Token string) *DiscordBot {

	bot := DiscordBot{}

	session, err := discordgo.New("Bot " + Token)

	if err != nil {
		log.Fatal("error creating Discord session")
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages
	bot.Session = session
	session.AddHandler(bot.messageCreate)

	insults := insults.NewInsults()

	bot.Insults = &insults
	bot.Insults.LoadInsults()

	return &bot
}

func (b *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println("Message received ->", m.Content)

	messageReply := b.Insults.GetInsult()

	// _, err := s.ChannelMessageSend(m.ChannelID, "Hello to you too sir!")
	msg, err := s.ChannelMessageSendReply(m.ChannelID, messageReply, m.Reference())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Sent Message", msg.ID)
}

func (b *DiscordBot) Start(wg *sync.WaitGroup) {
	err := b.Session.Open()

	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	b.Session.Close()
	wg.Done()
}
