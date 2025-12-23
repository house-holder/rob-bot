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
	"sync"
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

	if icao, found := strings.CutPrefix(msg, "taf "); found {
		raw, err := fetchTAF(icao)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "somethin ain't right")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%s`", raw))
		return
	}

	if icao, found := strings.CutPrefix(msg, "wx "); found {
		icao = strings.ToUpper(icao)
		var metar, taf string
		var metarErr, tafErr error
		var wg sync.WaitGroup

		wg.Add(2)

		go func() {
			defer wg.Done()
			metar, metarErr = fetchMETAR(icao)
		}()

		go func() {
			defer wg.Done()
			taf, tafErr = fetchTAF(icao)
		}()

		wg.Wait()

		if metarErr != nil && tafErr != nil {
			s.ChannelMessageSend(m.ChannelID, "somethin definitely ain't right")
			return
		}

		if metarErr != nil {
			metar = "[METAR unavailable]"
		}
		if tafErr != nil {
			taf = "[TAF unavailable]"
		}

		reply := fmt.Sprintf("```\n%s\n\n%s\n```", metar, taf)
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
}

func fetchMETAR(icao string) (string, error) {
	rootURL := "https://aviationweather.gov/api/data"
	url := fmt.Sprintf("%s/metar?ids=%s&format=json", rootURL, icao)

	client := &http.Client{Timeout: 7 * time.Second}
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

func fetchTAF(icao string) (string, error) {
	rootURL := "https://aviationweather.gov/api/data"
	url := fmt.Sprintf("%s/taf?ids=%s&format=json", rootURL, icao)

	client := &http.Client{Timeout: 7 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned %d", resp.StatusCode)
	}

	var data []struct {
		RawTAF string `json:"rawTAF"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}
	if len(data) == 0 {
		return "", fmt.Errorf("no METAR for %s", icao)
	}

	raw := data[0].RawTAF
	raw = strings.ReplaceAll(raw, " FM", "\n  FM")
	raw = strings.ReplaceAll(raw, " PROB", "\n    PROB")
	raw = strings.ReplaceAll(raw, " TEMPO", "\n    TEMPO")
	raw = strings.ReplaceAll(raw, " BECMG", "\n    BECMG")
	return raw, nil
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
