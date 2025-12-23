package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
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
		raw, err := fetchMETAR(icao)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "somethin ain't right")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%s`", raw))
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
}

func fetchMETAR(icao string) (string, error) {
	url := fmt.Sprintf("https://aviationweather.gov/api/data/metar?ids=%s&format=json", icao)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned %d", resp.StatusCode)
	}

	var data []struct {
		RawOb string `json:"rawOb"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}
	if len(data) == 0 {
		return "", fmt.Errorf("no METAR for %s", icao)
	}

	return data[0].RawOb, nil
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
