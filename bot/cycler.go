package bot

import (
	"math/rand"

	"rob-bot/data"
)

// MsgCycler -
type MsgCycler struct {
	messages []string
	indices  []int
	current  int
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

func init() {
	horseCycler = NewMsgCycler(data.HorseFacts)
	reaganCycler = NewMsgCycler(data.ReaganFacts)
	triviaCycler = NewMsgCycler(data.Trivia)
	wisdomCycler = NewMsgCycler(data.Wisdom)
}
