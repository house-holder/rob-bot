package bot

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"rob-bot/data"
)

var tenor = "https://tenor.com/view"
var nod = "/owyeah-gif-5873858916975845615"
var avTrigger = "lay it on me"
var triggers = []string{"reagan", "ronald", "nancy"}

func handleMisc(s *discordgo.Session, m *discordgo.MessageCreate, msg string) bool {
	if contains(msg, ", do you agree?") ||
		contains(msg, "who agrees?") ||
		contains(msg, "someone agree") {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s", tenor, nod))
		return true
	}

	for _, t := range triggers {
		if contains(msg, t) {
			if t == "nancy" {
				s.ChannelMessageSend(m.ChannelID, "hey! fuck the reagans")
				return true
			}
			s.ChannelMessageSend(m.ChannelID, "oh yeah, fuck ronald reagan")
			return true
		}
	}

	if contains(msg, avTrigger) {
		list := rand.Intn(10)
		if list < 5 {
			mem := data.MemoryItems[rand.Intn(len(data.MemoryItems))]
			fullMsg := fmt.Sprintf("**Remember:** %s", mem)
			s.ChannelMessageSend(m.ChannelID, fullMsg)
			return true
		} else {
			trivia := data.TransportTrivia[rand.Intn(len(data.TransportTrivia))]
			fullMsg := fmt.Sprintf(">>> %s", trivia)
			s.ChannelMessageSend(m.ChannelID, fullMsg)
			return true
		}
	}

	return false
}

