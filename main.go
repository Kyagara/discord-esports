package main

import (
	"discord-esports/models"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	CONFIG_FILE_PATH                 string = "./data/config.json"
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

func loadWikiData() error {
	files, err := os.ReadDir(CHAMPIONS_FOLDER_PATH)
	if err != nil {
		return fmt.Errorf("error reading '%v' dir: %v", CHAMPIONS_FOLDER_PATH, err)
	}

	championsNames = make(map[string]string)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		file, err := os.ReadFile(fmt.Sprintf("%v/%v", NORMALIZED_CHAMPIONS_FOLDER_PATH, f.Name()))
		if err != nil {
			return fmt.Errorf("error reading '%v/%v': %v", NORMALIZED_CHAMPIONS_FOLDER_PATH, f.Name(), err)
		}

		champion := models.Champion{}
		err = json.Unmarshal(file, &champion)
		if err != nil {
			return fmt.Errorf("error parsing json '%v': %v", f.Name(), err)
		}

		spellsInfo[champion.Key] = make([]SpellInfo, 0)
		spellsEmbeds[champion.Key] = make(map[string][]SpellEmbeds)

		for i, spell := range champion.Spells.Passive {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "P", i))
			spellsEmbeds[champion.Key]["P"] = append(spellsEmbeds[champion.Key]["P"], createChampionSpellEmbed(&champion, &spell, "P"))
		}

		for i, spell := range champion.Spells.Q {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "Q", i))
			spellsEmbeds[champion.Key]["Q"] = append(spellsEmbeds[champion.Key]["Q"], createChampionSpellEmbed(&champion, &spell, "Q"))
		}

		for i, spell := range champion.Spells.W {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "W", i))
			spellsEmbeds[champion.Key]["W"] = append(spellsEmbeds[champion.Key]["W"], createChampionSpellEmbed(&champion, &spell, "W"))
		}

		for i, spell := range champion.Spells.E {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "E", i))
			spellsEmbeds[champion.Key]["E"] = append(spellsEmbeds[champion.Key]["E"], createChampionSpellEmbed(&champion, &spell, "E"))
		}

		for i, spell := range champion.Spells.R {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "R", i))
			spellsEmbeds[champion.Key]["R"] = append(spellsEmbeds[champion.Key]["R"], createChampionSpellEmbed(&champion, &spell, "R"))
		}

		championsEmbeds[champion.Key] = createChampionEmbed(&champion)
		championsNames[champion.Name] = champion.Key
	}

	return nil
}
