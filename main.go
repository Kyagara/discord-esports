package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Kyagara/equinox/clients/ddragon"
)

const embedColor int = 0xff3838

// Flags
var (
	registerFlag = flag.Bool("register", false, "Register commands to the guild specified in the config files.")
	removeFlag   = flag.Bool("remove", false, "Remove commands to the guild specified in the config files.")
)

var (
	client  *Client   = &Client{}
	started time.Time = time.Now()

	ddVersion        string
	ddVersionUpdated time.Time

	champions      map[string]ddragon.ChampionData
	championsNames []string

	lastUpdate time.Time
	lastPost   time.Time
	now        time.Time
	tomorrow   time.Time

	lolSchedule map[string][]LOLEsportsLeagueSchedule     = make(map[string][]LOLEsportsLeagueSchedule)
	valSchedule map[string][]VALEsportsTournamentSchedule = make(map[string][]VALEsportsTournamentSchedule)
)

func init() { flag.Parse() }

func main() {
	if *registerFlag && *removeFlag {
		log.Fatal("Use only one of these commands at a time.")
	}

	client, err := newClient()
	if err != nil {
		log.Fatal(fmt.Errorf("error creating client: %v", err))
	}

	client.connect()
	defer client.session.Close()

	if *registerFlag {
		err := client.registerCommands()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *removeFlag {
		err := client.removeCommands()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	client.mainLoop()
}
