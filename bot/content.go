package bot

import (
	"fmt"
	"math/rand"

	"rob-bot/data"

	"github.com/bwmarrin/discordgo"
)

// MsgCycler -
type MsgCycler struct {
	messages []string
	indices  []int
	current  int
}

// TextCmd -
type TextCmd struct {
	triggers []string
	cycler   *MsgCycler
	format   string
	reaction string
}

// NewMsgCycler -
func NewMsgCycler(input []string) *MsgCycler {
	m := &MsgCycler{
		messages: input,
		indices:  nil,
		current:  0,
	}
	m.Shuffle()
	return m
}

// Shuffle -
func (m *MsgCycler) Shuffle() {
	ids := make([]int, len(m.messages))
	for i := range ids {
		ids[i] = i
	}
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	m.current = 0
	m.indices = ids
}

// Next -
func (m *MsgCycler) Next() string {
	if m.current >= len(m.indices) {
		m.Shuffle()
	}
	output := m.messages[m.indices[m.current]]
	m.current++
	return output
}

var (
	horseCycler  *MsgCycler
	reaganCycler *MsgCycler
	triviaCycler *MsgCycler
	wisdomCycler *MsgCycler
)

var textCommands = []TextCmd{}

func init() {
	horseCycler = NewMsgCycler(data.HorseFacts)
	reaganCycler = NewMsgCycler(data.ReaganFacts)
	triviaCycler = NewMsgCycler(data.Trivia)
	wisdomCycler = NewMsgCycler(data.Wisdom)

	textCommands = []TextCmd{
		{
			triggers: []string{"horse fact", "house fact", "honse?"},
			cycler:   horseCycler,
			format:   "**Horse Fact:**\n>>> %s",
			reaction: "🐴",
		},
		{
			triggers: []string{"reagan fact", "reagan facts", "reaganfacts"},
			cycler:   reaganCycler,
			format:   ">>> %s",
		},
		{
			triggers: []string{"trivia"},
			cycler:   triviaCycler,
			format:   ">>> %s",
		},
		{
			triggers: []string{"wisdom"},
			cycler:   wisdomCycler,
			format:   "**Remember:** %s",
		},
	}
}

func handleTextCmd(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	msg string,
	cmd TextCmd,
) bool {
	for _, trigger := range cmd.triggers {
		if contains(msg, trigger) {
			if cmd.reaction != "" {
				s.MessageReactionAdd(m.ChannelID, m.ID, cmd.reaction)
			}
			text := cmd.cycler.Next()
			fullMsg := fmt.Sprintf(cmd.format, text)
			s.ChannelMessageSend(m.ChannelID, fullMsg)
			return true
		}
	}
	return false
}

func handleMessages(s *discordgo.Session, m *discordgo.MessageCreate, msg string) bool {
	for _, cmd := range textCommands {
		if handleTextCmd(s, m, msg, cmd) {
			return true
		}
	}
	return false
}
