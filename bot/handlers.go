package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MsgCreate -
func MsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	msg := strings.ToLower(m.Content)

	if handleWeather(s, m, msg) {
		return
	}

	if handleMessages(s, m, msg) {
		return
	}

	handleReactions(s, m, msg)
	handleMisc(s, m, msg)
}

func contains(s, substr string) bool { // case-insensitive compare
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
