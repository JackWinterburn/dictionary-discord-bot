package main

import (
	"fmt"
	"os"
	"strings"
  "net/http"
  "io/ioutil"
  "encoding/json"

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

// format the command from a message to suit wikipedias criteria
func formatMessage(msg string) string {
	msg = msg[1:]
	msg = strings.ToLower(msg)
	msg = strings.Replace(msg, " ", "_", -1)
	return msg
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}
	// make sure the bot only replies in science channels
	for i := 0; i < len(channels); i++ {
		if m.ChannelID == channels[i] {
			if strings.ToLower(m.Content[0:1]) == "!" {
        resp, _ := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + m.Content)
        defer resp.Body.Close()

        jsonResponse, _ := ioutil.ReadAll(resp.Body)
        var definition []byte
        json.Unmarshal([]byte(jsonResponse), &definition)
        fmt.Printf(string(definition))

        fmt.Printf("%s\n", jsonResponse)
				_, _ = s.ChannelMessageSend(m.ChannelID, "https://en.wikipedia.org/wiki/"+m.Content)
			}
		}
	}
}

func main() {
	openConnection()
}
