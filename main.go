package main

import (
	"flag"
	"fmt"
	"time"
)

const (
	CONFIG_FILE_PATH                 string = "./data/config.json"
	ESPORTS_FILE_PATH                string = "./data/esports.json"
	CHAMPIONS_FOLDER_PATH            string = "./data/champions"
	NORMALIZED_CHAMPIONS_FOLDER_PATH string = "./data/champions/normalized"
	// NORMALIZED_ITEMS_FOLDER_PATH string = "./data/items/normalized"

	// Experimental
	// LOL_NEWS_FILE_PATH               string = "./data/lol_news.json"
	// VAL_NEWS_FILE_PATH               string = "./data/val_news.json"

	DISCORD_EMBED_COLOR int = 0xff3838
)

// Flags
var (
	registerFlag = flag.Bool("register", false, "Register commands to the guild specified in the config files.")
)

var (
	client  *Client   = &Client{}
	started time.Time = time.Now()

	commandButtonsID []string = []string{"modifiers", "notes", "spell"}

	// championsNames[Bel'veth] = Belveth
	championsNames map[string]string = make(map[string]string)

	// spells[championKey]
	spellsInfo map[string][]SpellInfo = make(map[string][]SpellInfo)

	// spellsEmbeds[championKey][spellKey][spellIndex]
	spellsEmbeds map[string]map[string][]SpellEmbeds = make(map[string]map[string][]SpellEmbeds)

	// championsEmbeds[championKey]
	championsEmbeds map[string]ChampionEmbeds = make(map[string]ChampionEmbeds)

	now      time.Time
	tomorrow time.Time

	esports EsportsData = EsportsData{VALSchedule: make(map[string][]VALEsportsTournamentSchedule), LOLSchedule: make(map[string][]LOLEsportsLeagueSchedule)}
)

func init() { flag.Parse() }

func main() {
	client, err := newClient()
	if err != nil {
		client.logger.Fatal(fmt.Sprintf("error creating client: %v", err))
	}

	if client.config.Commands.Champion || client.config.Commands.Spell {
		err = loadWikiData()
		if err != nil {
			client.logger.Fatal(fmt.Sprintf("error loading wiki data: %v", err))
		}
	}

	if client.config.Commands.Esports {
		err = loadOrCreateFile(ESPORTS_FILE_PATH, &esports)
		if err != nil {
			client.logger.Fatal(fmt.Sprintf("error loading esports data: %v", err))
		}
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

	client.mainLoop()
}
