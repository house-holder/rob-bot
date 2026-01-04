package bot

import (
	"fmt"
	"strings"

	"rob-bot/weather"

	"github.com/bwmarrin/discordgo"
)

// InteractionCreate -
func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "" {
		return
	}

	data := i.ApplicationCommandData()

	switch data.Name {
	case "horsefact":
		text := horseCycler.Next()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("**Horse Fact:**\n>>> %s", text),
			},
		})
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "🐴",
		})

	case "reaganfact":
		text := reaganCycler.Next()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(">>> %s", text),
			},
		})

	case "trivia":
		text := triviaCycler.Next()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(">>> %s", text),
			},
		})

	case "wisdom":
		text := wisdomCycler.Next()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("**Remember:** %s", text),
			},
		})

	case "metar":
		icao := data.Options[0].StringValue()
		reply := weather.CmdMETAR(icao)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: reply,
			},
		})

	case "taf":
		icao := data.Options[0].StringValue()
		reply := weather.CmdTAF(icao)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: reply,
			},
		})

	case "wx":
		icao := data.Options[0].StringValue()
		reply := weather.CmdWX(icao)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: reply,
			},
		})

	case "atis":
		icao := data.Options[0].StringValue()
		reply, code, err := weather.CmdATIS(icao)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: reply,
			},
		})

		// Add reaction for ATIS code if available
		if len(code) > 0 {
			letter := strings.ToUpper(code)[0]
			if letter >= 'A' && letter <= 'Z' {
				emojiRune := '\U0001F1E6' + rune(letter-'A')
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: string(emojiRune),
				})
			}
		}
	case "go":
		icao := data.Options[0].StringValue()
		minimal, err := weather.CmdATISLetter(icao)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: minimal,
			},
		})
	}
}

// MsgCreate -
func MsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	msg := strings.ToLower(m.Content)

	handleReactions(s, m, msg)
	handleMisc(s, m, msg)
}

func contains(s, substr string) bool { // case-insensitive compare
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
