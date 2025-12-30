package bot

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"rob-bot/data"
)

var factTriggers = []string{
	"reagan fact", "reagan facts", "reaganfacts",
}

func handleFacts(s *discordgo.Session, m *discordgo.MessageCreate, msg string) bool {
	if contains(msg, "horse fact") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🐴")
		fact := data.HorseFacts[rand.Intn(len(data.HorseFacts))]
		fullMsg := fmt.Sprintf("**Horse Fact:**\n>>> %s", fact)
		s.ChannelMessageSend(m.ChannelID, fullMsg)
		return true
	}

	for _, ft := range factTriggers {
		if contains(msg, ft) {
			fact := data.ReaganFacts[rand.Intn(len(data.ReaganFacts))]
			fullMsg := fmt.Sprintf(">>> %s", fact)
			s.ChannelMessageSend(m.ChannelID, fullMsg)
			return true
		}
	}

	return false
}

