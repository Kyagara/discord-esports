package main

import (
	"discord-esports/models"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type EsportsData struct {
	LastPostTimestamp   time.Time
	LastUpdateTimestamp time.Time
	LOLSchedule         map[string][]LOLEsportsLeagueSchedule     `json:"lolSchedule"`
	VALSchedule         map[string][]VALEsportsTournamentSchedule `json:"valSchedule"`
	m                   sync.Mutex
}

// Loads or create json file
func loadOrCreateFile(filePath string, target interface{}) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.WriteFile(filePath, []byte("{}"), 0644)
		if err != nil {
			return fmt.Errorf("error creating '%v' file: %v", filePath, err)
		}
	}

	err := loadFile(filePath, &target)
	if err != nil {
		return fmt.Errorf("error loading file '%v': %v", filePath, err)
	}
	return nil
}

func loadFile(filePath string, target interface{}) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file '%v': %v", filePath, err)
	}

	err = json.Unmarshal(file, &target)
	if err != nil {
		return fmt.Errorf("error parsing json '%v': %v", filePath, err)
	}
	return nil
}

func saveEsportsFile() {
	client.logger.Info("Saving esports data file.")
	esports.m.Lock()
	defer esports.m.Unlock()
	file, err := json.MarshalIndent(&esports, "", " ")
	if err != nil {
		client.logger.Error(fmt.Sprintf("error parsing esports data file: %v", err))
	}
	err = os.WriteFile(ESPORTS_FILE_PATH, file, 0644)
	if err != nil {
		client.logger.Error(fmt.Sprintf("Error writing esports data file: %v", err))
	}
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

		champion := models.Champion{}
		err = loadFile(fmt.Sprintf("%v/%v", NORMALIZED_CHAMPIONS_FOLDER_PATH, f.Name()), &champion)
		if err != nil {
			return fmt.Errorf("error loading champion file '%v': %v", f.Name(), err)
		}

		spellsInfo[champion.Key] = make([]SpellInfo, 0)
		spellsEmbeds[champion.Key] = make(map[string][]SpellEmbeds)

		for i, spell := range champion.Spells.Passive {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "P", i))
			spellsEmbeds[champion.Key]["P"] = append(spellsEmbeds[champion.Key]["P"], createChampionSpellEmbed(&champion, &spell))
		}

		for i, spell := range champion.Spells.Q {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "Q", i))
			spellsEmbeds[champion.Key]["Q"] = append(spellsEmbeds[champion.Key]["Q"], createChampionSpellEmbed(&champion, &spell))
		}

		for i, spell := range champion.Spells.W {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "W", i))
			spellsEmbeds[champion.Key]["W"] = append(spellsEmbeds[champion.Key]["W"], createChampionSpellEmbed(&champion, &spell))
		}

		for i, spell := range champion.Spells.E {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "E", i))
			spellsEmbeds[champion.Key]["E"] = append(spellsEmbeds[champion.Key]["E"], createChampionSpellEmbed(&champion, &spell))
		}

		for i, spell := range champion.Spells.R {
			spellsInfo[champion.Key] = append(spellsInfo[champion.Key], createSpellInfo(&spell, "R", i))
			spellsEmbeds[champion.Key]["R"] = append(spellsEmbeds[champion.Key]["R"], createChampionSpellEmbed(&champion, &spell))
		}

		championsEmbeds[champion.Key] = createChampionEmbed(&champion)
		championsNames[champion.Name] = champion.Key
	}

	return nil
}
