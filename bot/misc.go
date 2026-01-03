package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var tenor = "https://tenor.com/view"
var nod = "/owyeah-gif-5873858916975845615"
var cheeks = "/nsfw-gif-12427338058739620942"
var triggers = []string{"reagan", "ronald", "nancy"}

func handleMisc(s *discordgo.Session, m *discordgo.MessageCreate, msg string) bool {
	if contains(msg, ", do you agree?") ||
		contains(msg, "who agrees?") ||
		contains(msg, "someone agree") ||
		contains(msg, "think, bob") ||
		contains(msg, "think, rob?") {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s", tenor, nod))
		return true
	}

	if contains(msg, "cheek") {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s", tenor, cheeks))
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

	return false
}
