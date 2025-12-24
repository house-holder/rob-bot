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
	"time"
)

var factTriggers = []string{
	"reagan facts", "reaganfacts", "reaganfax", "rgnfax", "fax", "fact", "facts",
}
var triggers = []string{"reagan", "ronald", "nancy"}

func msgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // detect self-messages and ignore
		return
	}
	msg := strings.ToLower(m.Content)

	if icao, found := strings.CutPrefix(msg, "metar "); found {
		reply := cmdMETAR(icao)
		s.ChannelMessageSend(m.ChannelID, reply)
		return
	}

	if icao, found := strings.CutPrefix(msg, "taf "); found {
		reply := cmdTAF(icao)
		s.ChannelMessageSend(m.ChannelID, reply)
		return
	}

	if icao, found := strings.CutPrefix(msg, "wx "); found {
		reply := cmdWX(icao)
		s.ChannelMessageSend(m.ChannelID, reply)
		return
	}

	for _, ft := range factTriggers {
		if strings.Contains(msg, ft) {
			fact := facts[rand.Intn(len(facts))]
			fullMessage := fmt.Sprintf(">>> %s", fact)
			s.ChannelMessageSend(m.ChannelID, fullMessage)
			return
		}
	}

	for _, t := range triggers {
		if strings.Contains(msg, t) {
			if t == "nancy" {
				s.ChannelMessageSend(m.ChannelID, "fuck the reagans")
				return
			}
			s.ChannelMessageSend(m.ChannelID, "fuck ronald reagan")
			return
		}
	}

	t := time.Now()
	if strings.Contains(msg, "christmas eve") && t.Month() == time.December && t.Day() == 24 {
		s.ChannelMessageSend(m.ChannelID, "it's christmas eve, *not* christmas steve...")
		return
	}

	if strings.Contains(msg, "ts ") || strings.HasPrefix(msg, "ts ") {
		s.ChannelMessageSend(m.ChannelID, "ts pmo fr fr")
		return
	}
}

func main() {
	_ = godotenv.Load() // failsafe for dev/prod
	dg, err := discordgo.New("Bot " + os.Getenv("ROB_BOT"))
	if err != nil {
		log.Fatal("Create Discord session fail: ", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentMessageContent
	dg.AddHandler(msgCreate)

	dg.Open()
	log.Println("rob-bot is running, ctrl+c to stop")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	<-sc
	log.Println("stopping rob-bot")
	dg.Close()
}
