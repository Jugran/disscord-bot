package bot

import (
	"diss-cord/models"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Session *discordgo.Session
	Insults *models.InsultData
}

func NewBot(Token string) *DiscordBot {

	bot := DiscordBot{}

	session, err := discordgo.New("Bot " + Token)

	if err != nil {
		log.Fatal("error creating Discord session")
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages
	session.AddHandler(bot.messageCreate)
	bot.Session = session

	insults := models.NewInsults()

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

func (b *DiscordBot) Start(wg *sync.WaitGroup, quit <-chan bool) {
	defer wg.Done()
	err := b.Session.Open()
	// Cleanly close down the Discord session.
	defer b.Session.Close()

	if err != nil {
		log.Fatal("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	<-quit
	fmt.Println("Shutting down bot...")
}
