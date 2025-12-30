package bot

import "github.com/bwmarrin/discordgo"

func handleReactions(s *discordgo.Session, m *discordgo.MessageCreate, msg string) bool {
	if contains(msg, "horse") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🐴")
	}
	if contains(msg, "dub") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇼")
	}
	return false
}

