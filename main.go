package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
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

	// championsNames[Bel'veth] = Belveth
	championsNames map[string]string = make(map[string]string)

	// spells[championKey]
	spellsInfo map[string][]SpellInfo = make(map[string][]SpellInfo)

	// spellsEmbeds[championKey][spellKey][spellIndex]
	spellsEmbeds map[string]map[string][]SpellEmbeds = make(map[string]map[string][]SpellEmbeds)

	// championsEmbeds[championKey]
	championsEmbeds map[string]ChampionEmbeds = make(map[string]ChampionEmbeds)

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

	files, err := os.ReadDir("./champions")
	if err != nil {
		log.Fatal(err)
	}

	championsNames = make(map[string]string)

	for _, f := range files {
		file, err := os.ReadFile(fmt.Sprintf("./champions/%v", f.Name()))
		if err != nil {
			log.Fatal(err)
		}

		champion := WikiChampion{}
		err = json.Unmarshal(file, &champion)
		if err != nil {
			log.Fatal(err)
		}

		spellsInfo[champion.Key] = make([]SpellInfo, 0)
		spellsEmbeds[champion.Key] = make(map[string][]SpellEmbeds)

		for _, spell := range champion.Spells.Passive {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "P", -1))
			spellsEmbeds[champion.Key]["P"] = append(spellsEmbeds[champion.Key]["P"], createChampionSpellEmbed(&champion, &spell, "P"))
		}
		for _, spell := range champion.Spells.Q {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "Q", -1))
			spellsEmbeds[champion.Key]["Q"] = append(spellsEmbeds[champion.Key]["Q"], createChampionSpellEmbed(&champion, &spell, "Q"))
		}
		for _, spell := range champion.Spells.W {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "W", -1))
			spellsEmbeds[champion.Key]["W"] = append(spellsEmbeds[champion.Key]["W"], createChampionSpellEmbed(&champion, &spell, "W"))
		}
		for _, spell := range champion.Spells.E {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "E", -1))
			spellsEmbeds[champion.Key]["E"] = append(spellsEmbeds[champion.Key]["E"], createChampionSpellEmbed(&champion, &spell, "E"))
		}
		for _, spell := range champion.Spells.R {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "R", -1))
			spellsEmbeds[champion.Key]["R"] = append(spellsEmbeds[champion.Key]["R"], createChampionSpellEmbed(&champion, &spell, "R"))
		}

		championsEmbeds[champion.Key] = createChampionEmbed(&champion)
		championsNames[champion.Name] = champion.Key
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
