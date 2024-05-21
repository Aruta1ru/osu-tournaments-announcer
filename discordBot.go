package main

import (
	"discord-go/events"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var Token string

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	discordToken, exists := os.LookupEnv("DISCORD_TOKEN")
	if exists {
		flag.StringVar(&Token, "t", discordToken, "Bot token")
		flag.Parse()
	}
}

func main() {
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	s.AddHandler(events.Ready)
	s.AddHandler(events.MessageCreate)

	err = s.Open()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Bot is running.\nPress CTRL+C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	s.Close()

}
