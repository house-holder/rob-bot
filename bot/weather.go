package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"rob-bot/weather"
)

func handleWeather(s *discordgo.Session, m *discordgo.MessageCreate, msg string) bool {
	if icao, found := strings.CutPrefix(msg, "metar "); found {
		reply := weather.CmdMETAR(icao)
		s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
		return true
	}

	if icao, found := strings.CutPrefix(msg, "taf "); found {
		reply := weather.CmdTAF(icao)
		s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
		return true
	}

	if icao, found := strings.CutPrefix(msg, "wx "); found {
		reply := weather.CmdWX(icao)
		s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
		return true
	}

	if icao, found := strings.CutPrefix(msg, "atis "); found {
		reply, code, err := weather.CmdATIS(icao)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
			return true
		}

		if len(code) > 0 {
			letter := strings.ToUpper(code)[0]
			if letter >= 'A' && letter <= 'Z' {
				emojiRune := '\U0001F1E6' + rune(letter-'A')
				s.MessageReactionAdd(m.ChannelID, m.ID, string(emojiRune))
			}
		}

		s.ChannelMessageSend(m.ChannelID, reply)
		return true
	}

	return false
}

