package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var tenor = "https://tenor.com/view"
var nod = "/owyeah-gif-5873858916975845615"

var avTrigger = "lay it on me"
var factTriggers = []string{
	"reagan fact", "reagan facts", "reaganfacts",
}
var triggers = []string{"reagan", "ronald", "nancy"}

func msgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // detect self-messages and ignore
		return
	}
	msg := strings.ToLower(m.Content)

	if strings.Contains(msg, ", do you agree?") ||
		strings.Contains(msg, "who agrees?") ||
		strings.Contains(msg, "someone agree") {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s", tenor, nod))
	}
	if strings.Contains(msg, "horse") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🐴")
	}
	if strings.Contains(msg, "dub") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇼")
	}

	// if m.Author.ID == "822006025229959168" {
	// 	s.MessageReactionAdd(m.ChannelID, m.ID, "🐴")
	// }

	// if (strings.Contains(msg, "6") || strings.Contains(msg, "six")) &&
	// 	(strings.Contains(msg, "7") || strings.Contains(msg, "seven")) {
	// 	s.MessageReactionAdd(m.ChannelID, m.ID, "6️⃣")
	// 	s.MessageReactionAdd(m.ChannelID, m.ID, "7️⃣")
	// }

	if icao, found := strings.CutPrefix(msg, "metar "); found {
		reply := cmdMETAR(icao)
		s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
		return
	}

	if icao, found := strings.CutPrefix(msg, "taf "); found {
		reply := cmdTAF(icao)
		s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
		return
	}

	if icao, found := strings.CutPrefix(msg, "wx "); found {
		reply := cmdWX(icao)
		s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
		return
	}

	if strings.Contains(msg, "horse fact") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🐴")
		fact := horseFacts[rand.Intn(len(horseFacts))]
		fullMsg := fmt.Sprintf("**Horse Fact:**\n>>> %s", fact)
		s.ChannelMessageSend(m.ChannelID, fullMsg)
		return
	}

	for _, ft := range factTriggers {
		if strings.Contains(msg, ft) {
			fact := reaganFacts[rand.Intn(len(reaganFacts))]
			fullMsg := fmt.Sprintf(">>> %s", fact)
			s.ChannelMessageSend(m.ChannelID, fullMsg)
			return
		}
	}

	for _, t := range triggers {
		if strings.Contains(msg, t) {
			if t == "nancy" {
				s.ChannelMessageSend(m.ChannelID, "hey! fuck the reagans")
				return
			}
			s.ChannelMessageSend(m.ChannelID, "oh yeah, fuck ronald reagan")
			return
		}
	}

	if strings.Contains(msg, avTrigger) {
		list := rand.Intn(10)
		if list < 5 {
			mem := memoryItems[rand.Intn(len(memoryItems))]
			fullMsg := fmt.Sprintf("**Remember:** %s", mem)
			s.ChannelMessageSend(m.ChannelID, fullMsg)
			return
		} else {
			trivia := transportTrivia[rand.Intn(len(transportTrivia))]
			fullMsg := fmt.Sprintf(">>> %s", trivia)
			s.ChannelMessageSend(m.ChannelID, fullMsg)
			return
		}
	}
}

func main() {
	_ = godotenv.Load() // failsafe for dev/prod
	dg, err := discordgo.New("Bot " + os.Getenv("ROB_BOT"))
	if err != nil {
		log.Fatal("Create Discord session fail: ", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildMessageReactions // NOTE: assumed needed this
	dg.AddHandler(msgCreate)

	dg.Open()
	log.Println("rob-bot is running, ctrl+c to stop")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	<-sc
	log.Println("stopping rob-bot")
	dg.Close()
}
