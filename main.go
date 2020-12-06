package main

import (
	f "fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// BotID will be the alias for the bot
var BotID string
var token string = os.Getenv("TOKEN")
var channels = [...]string{"773318228184137759", "773318250805198848", "773318284325027883"}

func openConnection() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		f.Println(err.Error())
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		f.Println(err.Error())
	}

	BotID = u.ID

	dg.AddHandler(messageHandler)

	err = dg.Open()
	if err != nil {
		f.Println(err.Error())
		return
	}

	f.Println("bot is running!")

	// keeps the bot running
	<-make(chan struct{})
	return
}

func in(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// parses a json string to get definitions from dictionary api
func parseJSONForDefinition(JSONString string) []string {
	definitions := strings.Split(JSONString, "\"definition\":")
	definitions = definitions[1:]

	for def := range definitions {
		definition := ""
		str := definitions[def]
		for idx := 0; idx < len(str); idx++ {
			currLetter := string(str[idx+1])
			if currLetter == "\"" {
				definitions[def] = definition
				break
			} else {
				definition += currLetter
			}
		}
	}
	return definitions
}

func convertDefinitionsToString(msg string, definitions ...string) string {
	str := "Definitions for " + msg + ":\n\n"
	for idx, def := range definitions {
		if idx == 0 {
			str += ">>> " + "• " + def
		} else {
			str += "\n\n ━━━━━━━━━━━━━━━━━━\n" + "• " + def
		}
	}
	return str
}

func getDefenitionsFromJSONString(JSONString string, msg string) string {
	// the api will return this JSON string if no defenitions are found
	errmsg := `{"title":"No Definitions Found","message":"Sorry pal, we couldn't find definitions for the word you were looking for.","resolution":"You can try the search again at later time or head to the web instead."}`

	if JSONString == errmsg {
		return "I couldn't find a definition for you..."
	}

	parsedDefinitions := parseJSONForDefinition(JSONString)
	definitions := convertDefinitionsToString(msg, parsedDefinitions...)

	return definitions
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}
	// make sure the bot only replies in science channels
	for i := 0; i < len(channels); i++ {
		if m.ChannelID == channels[i] {
			if strings.ToLower(m.Content[0:2]) == "d!" {
				m.Content = m.Content[2:]
				resp, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + m.Content)
				if err != nil {
					f.Println(err)
				}
				defer resp.Body.Close()

				// get the actual JSON from the api
				JSONResponse, _ := ioutil.ReadAll(resp.Body)

				m.Content = "**" + m.Content + "**"
				definitions := getDefenitionsFromJSONString(string(JSONResponse), m.Content)
				_, _ = s.ChannelMessageSend(m.ChannelID, definitions)
			}
		}
	}
}

func main() {
	openConnection()
}
