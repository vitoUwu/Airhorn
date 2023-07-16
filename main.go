package main

import (
	"gobot/commands/airhorn"
	"gobot/commands/ping"
	"gobot/listeners/interactionCreate"
	"gobot/listeners/ready"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	token, exists := os.LookupEnv("DISCORD_TOKEN")
	if !exists {
		log.Fatal("missing DISCORD_TOKEN env")
	}

	applicationId, exists := os.LookupEnv("APPLICATION_ID")
	if !exists {
		log.Fatal("missing APPLICATION_ID env")
	}

	err = airhorn.LoadSound()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := disgolf.New(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates
	bot.Router.Register(ping.Data)
	bot.Router.Register(airhorn.Data)
	bot.AddHandler(bot.Router.HandleInteraction)
	bot.AddHandler(interactionCreate.Handler)
	bot.AddHandler(ready.Handler)

	err = bot.Router.Sync(bot.Session, applicationId, "")
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer bot.Close()

	stchan := make(chan os.Signal, 1)
	signal.Notify(stchan, syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)

end:
	for {
		select {
		case <-stchan:
			break end
		default:
		}
		time.Sleep(time.Second)
	}
}
