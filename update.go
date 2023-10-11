package main

import (
	"fmt"
	"time"

	"github.com/Kyagara/equinox/clients/ddragon"
	"github.com/bwmarrin/discordgo"
)

func UpdateCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if hasPermissions(session, interaction) {
		err := updateAllData()
		if err != nil {
			respondWithError(interaction.Interaction, err)
			return
		}

		respondWithMessage(interaction.Interaction, "Updated all data without any errors.")
	}
}

func updateAllData() error {
	now = time.Now()
	tomorrow = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	err := updateLOLData()
	if err != nil {
		return fmt.Errorf("error updating LOL data: %v", err)
	}

	err = updateVALData()
	if err != nil {
		return fmt.Errorf("error updating VAL data: %v", err)
	}

	ddVersion, err = client.equinox.DDragon.Version.Latest()
	if err != nil {
		return fmt.Errorf("error updating data dragon version: %v", err)
	}

	ddVersionUpdated = time.Now()

	champions, err = client.equinox.DDragon.Champion.AllChampions(ddVersion, ddragon.EnUS)
	if err != nil {
		return fmt.Errorf("error updating champions data: %v", err)
	}

	championsNames = make(map[string]string)
	for _, c := range champions {
		championsNames[c.Name] = c.ID
	}

	lastUpdate = time.Now()
	return nil
}
