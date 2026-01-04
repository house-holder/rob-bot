package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"rob-bot/bot"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func registerCommands(s *discordgo.Session) {
	commands := []discordgo.ApplicationCommand{
		{
			Name:        "horsefact",
			Description: "Get a random horse fact",
		},
		{
			Name:        "reaganfact",
			Description: "Get a fact about Ronald Reagan",
		},
		{
			Name:        "trivia",
			Description: "Get aviation trivia",
		},
		{
			Name:        "wisdom",
			Description: "Get aviation aphorisms",
		},
		{
			Name:        "metar",
			Description: "Fetch METAR for an airport",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "icao",
					Description: "ICAO airport code (KDVT, KSGU)",
					Required:    true,
				},
			},
		},
		{
			Name:        "taf",
			Description: "Fetch TAF for an airport",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "icao",
					Description: "ICAO airport code (KFNT, KMCO)",
					Required:    true,
				},
			},
		},
		{
			Name:        "wx",
			Description: "Fetch all available weather info for a given airport",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "icao",
					Description: "ICAO airport code (KCGI, KDEN)",
					Required:    true,
				},
			},
		},
		{
			Name:        "atis",
			Description: "Fetch the D-ATIS for an airport, if available",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "icao",
					Description: "ICAO airport code (KSTL, KIAD)",
					Required:    true,
				},
			},
		},
		{
			Name:        "go",
			Description: "D-ATIS for people in a hurry (if available, incl. age/timestamp)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "icao",
					Description: "ICAO airport code (KSFO, KDEN)",
					Required:    true,
				},
			},
		},
	}

	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &cmd)
		if err != nil {
			log.Printf("Cannot create command %v: %v", cmd.Name, err)
		}
	}

}

func main() {
	_ = godotenv.Load()
	dg, err := discordgo.New("Bot " + os.Getenv("ROB_BOT"))
	if err != nil {
		log.Fatal("Create Discord session fail: ", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentsGuilds

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		registerCommands(s)
	})

	dg.AddHandler(bot.MsgCreate)
	dg.AddHandler(bot.InteractionCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal("Failed to open connection: ", err)
	}
	log.Println("rob-bot is running, ctrl+c to stop")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	<-sc
	log.Println("stopping rob-bot")
	dg.Close()
}
