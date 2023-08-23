package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	registerFlag = flag.Bool("register", false, "Register commands to the guild specified in the config files.")
	removeFlag   = flag.Bool("remove", false, "Remove commands to the guild specified in the config files.")

	session      *discordgo.Session
	config       BotConfiguration = BotConfiguration{}
	lastUpdate   time.Time
	lastPost     time.Time
	tomorrowDate string = ""

	lolSchedule map[string][]LOLEsportsLeagueSchedule     = make(map[string][]LOLEsportsLeagueSchedule)
	valSchedule map[string][]VALEsportsTournamentSchedule = make(map[string][]VALEsportsTournamentSchedule)
)

func init() { flag.Parse() }

func init() {
	if *registerFlag && *removeFlag {
		log.Fatal("Use only one of these commands at a time.")
	}

	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if *registerFlag {
		err := registerCommands()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *removeFlag {
		err := removeCommands()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	session.AddHandler(interactionsEvent)
	session.AddHandler(readyEvent)

	session.Identify.Intents = discordgo.IntentMessageContent | discordgo.IntentsGuildMessages

	err := session.Open()
	if err != nil {
		log.Fatalf("Error opening session: %v", err)
	}

	defer session.Close()

	ticker()
}

func ticker() {
	sender := time.NewTicker(time.Duration(config.SendToChannelTimeout * int(time.Millisecond)))
	defer sender.Stop()

	updater := time.NewTicker(time.Duration(config.DataUpdateTimeout * int(time.Millisecond)))
	defer updater.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	firstUpdate := true
	firstSend := true

loop:
	for {
		select {
		case <-stop:
			log.Print("Shutting down.")
			break loop
		case <-updater.C:
			if !firstUpdate {
				err := updateAllData()
				if err != nil {
					log.Printf("Error updating data: %v", err)
				}
			}

			firstUpdate = false

		case <-sender.C:
			if !firstSend {
				err := postAllData()
				if err != nil {
					log.Printf("Error sending data: %v", err)
				}
			}

			firstSend = false
		}
	}
}
