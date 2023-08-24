package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func UpdateCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if hasPermissions(session, interaction) {
		err := updateAllData()
		if err != nil {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(), Flags: discordgo.MessageFlagsEphemeral,
				},
			})
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
	now := time.Now()
	tomorrow = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	err := updateLOLData()
	if err != nil {
		return fmt.Errorf("error updating LOL data: %v", err)
	}

	err = updateVALData()
	if err != nil {
		return fmt.Errorf("error updating VAL data: %v", err)
	}

	lastUpdate = time.Now()
	return nil
}
