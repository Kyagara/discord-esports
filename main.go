package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Kyagara/equinox"
	"github.com/Kyagara/equinox/clients/ddragon"
	"github.com/bwmarrin/discordgo"
)

var (
	registerFlag = flag.Bool("register", false, "Register commands to the guild specified in the config files.")
	removeFlag   = flag.Bool("remove", false, "Remove commands to the guild specified in the config files.")

	session *discordgo.Session

	dd             *equinox.Equinox
	ddVersion      string
	champions      map[string]ddragon.ChampionData
	championsNames []string
	versionUpdated time.Time

	config     Configuration = Configuration{}
	lastUpdate time.Time
	lastPost   time.Time
	started    time.Time = time.Now()
	now        time.Time
	tomorrow   time.Time

	lolSchedule map[string][]LOLEsportsLeagueSchedule     = make(map[string][]LOLEsportsLeagueSchedule)
	valSchedule map[string][]VALEsportsTournamentSchedule = make(map[string][]VALEsportsTournamentSchedule)
)

const embedColor int = 0xff3838

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

	dd, err = equinox.NewClient("data_dragon")
	if err != nil {
		log.Fatalf("Error starting equinox client: %v", err)
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
	post := time.NewTicker(time.Duration(config.PostDataTimer * int(time.Millisecond)))
	defer post.Stop()

	update := time.NewTicker(time.Duration(config.UpdateDateTimer * int(time.Millisecond)))
	defer update.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// The first update/send is done in the ready event, there would be a race condition if it was done by the ticker
	firstUpdate := true
	firstSend := true

loop:
	for {
		select {
		case <-stop:
			log.Print("Shutting down.")
			break loop
		case <-update.C:
			if !firstUpdate {
				err := updateAllData()
				if err != nil {
					log.Printf("Error updating data: %v", err)
				}
			}

			firstUpdate = false

		case <-post.C:
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
