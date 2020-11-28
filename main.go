package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// BotID will be the alias for the bot
var BotID string
var token string = os.Getenv("TOKEN")
var channels = [...]string{"773318228184137759", "773318250805198848", "773318284325027883"}

func openConnection() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	dg.AddHandler(messageHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("bot is running!")

	// keeps the bot running
	<-make(chan struct{})
	return
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}
	// make sure the bot only replies in science channels
	for i := 0; i < len(channels); i++ {
		if m.ChannelID == channels[i] {
			if m.Content == "!osmosis" {
				_, _ = s.ChannelMessageSend(m.ChannelID, "osmosis is the net movement of water against the conc. gradient")
			}
		}
	}
}

func main() {
	openConnection()
}
