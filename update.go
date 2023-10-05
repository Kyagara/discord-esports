package main

import (
	"fmt"
	"time"

	"github.com/Kyagara/equinox/clients/data_dragon"
	"github.com/bwmarrin/discordgo"
)

func UpdateCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if hasPermissions(session, interaction) {
		err := updateAllData()
		if err != nil {
			respondWithErrorEmbed(interaction.Interaction, err)
			return
		}

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Updated all data without any errors.", Flags: discordgo.MessageFlagsEphemeral,
			},
		})
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

	ddVersion, err = dd.DataDragon.Version.Latest()
	if err != nil {
		return fmt.Errorf("error updating data dragon version: %v", err)
	}

	versionUpdated = time.Now()

	champions, err = dd.DataDragon.Champion.AllChampions(ddVersion, data_dragon.EnUS)
	if err != nil {
		return fmt.Errorf("error updating champions data: %v", err)
	}

	championsNames = []string{}
	for _, c := range champions {
		championsNames = append(championsNames, c.ID)
	}

	lastUpdate = time.Now()
	return nil
}
