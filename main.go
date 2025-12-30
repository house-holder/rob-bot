package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"rob-bot/bot"
)

func main() {
	_ = godotenv.Load()
	dg, err := discordgo.New("Bot " + os.Getenv("ROB_BOT"))
	if err != nil {
		log.Fatal("Create Discord session fail: ", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildMessageReactions
	dg.AddHandler(bot.MsgCreate)

	dg.Open()
	log.Println("rob-bot is running, ctrl+c to stop")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	<-sc
	log.Println("stopping rob-bot")
	dg.Close()
}
