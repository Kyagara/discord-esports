package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func PostCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if hasPermissions(session, interaction) {
		err := postAllData()
		if err != nil {
			respondWithError(interaction.Interaction, err)
			return
		}

		respondWithMessage(interaction.Interaction, "Sent all embeds without any errors.")
	}
}

func postAllData() error {
	err := sendLOLEmbed()
	if err != nil {
		return fmt.Errorf("error sending LOL embed: %v", err)
	}

	err = sendVALEmbed()
	if err != nil {
		return fmt.Errorf("error sending VAL embed: %v", err)
	}

	lastPost = time.Now()
	return nil
}
