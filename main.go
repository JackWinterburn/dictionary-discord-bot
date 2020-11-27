package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// BotID will be the alias for the bot
var BotID string

const token string = "NzgyMDEyMzg1NDU1MjQzMjc1.X8F_yQ.apG48ylBTX5gmJycAYJrwoSOcFc"

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
	if m.Content == "!osmosis" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "osmosis is the net movement of water against the conc. gradient")
	}
	fmt.Println(m.Content)
}

func main() {
	openConnection()
}
