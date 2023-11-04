package bot

import (
	"diss-cord/models"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Session     *discordgo.Session
	DefaultRole string
}

func NewBot(Token string, DefaultRole string) *DiscordBot {

	bot := DiscordBot{}

	session, err := discordgo.New("Bot " + Token)

	if err != nil {
		log.Fatal("error creating Discord session")
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages
	session.AddHandler(bot.messageCreate)

	bot.Session = session
	bot.DefaultRole = DefaultRole

	return &bot
}

func (b *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// fmt.Println("Message received ->", m.Content, "By user", m.Author.Username, "Discord ID", m.Author.ID)

	user, err := models.CheckDiscordUser(m.Author.ID)

	if err != nil {
		// user not found, add new user
		user.Name = &m.Author.Username
		roles := &[]string{b.DefaultRole}

		user.AddNewUser(roles)
	}

	// fetch insults for user
	insult, err := user.GetInsult()

	if err != nil {
		fmt.Println("Cannot fetch user insult:", err)
		return
	}

	messageReply := insult.Text

	// _, err := s.ChannelMessageSend(m.ChannelID, "Hello to you too sir!")
	msg, err := s.ChannelMessageSendReply(m.ChannelID, messageReply, m.Reference())
	if err != nil {
		return
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
