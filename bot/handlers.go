package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	msg := strings.ToLower(m.Content)

	handleReactions(s, m, msg)

	if handleWeather(s, m, msg) {
		return
	}

	if handleFacts(s, m, msg) {
		return
	}

	handleMisc(s, m, msg)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

