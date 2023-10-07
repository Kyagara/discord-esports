package main

import (
	"flag"
	"fmt"
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
		client.logger.Fatal("Use only one of these commands at a time.")
	}

	client, err := newClient()
	if err != nil {
		client.logger.Fatal(fmt.Sprintf("error creating client: %v", err))
	}

	err = client.connect()
	if err != nil {
		client.logger.Fatal(fmt.Sprintf("error connecting to Discord session: %v", err))
	}

	defer client.session.Close()

	if *registerFlag {
		err := client.registerCommands()
		if err != nil {
			client.logger.Fatal(fmt.Sprintf("error registering commands: %v", err))
		}
		return
	}

	if *removeFlag {
		err := client.removeCommands()
		if err != nil {
			client.logger.Fatal(fmt.Sprintf("error removing commands: %v", err))
		}
		return
	}

	client.mainLoop()
}
