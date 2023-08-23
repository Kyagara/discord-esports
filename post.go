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
				Content: "Sent all embeds without any errors.", Flags: discordgo.MessageFlagsEphemeral,
			},
		})
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
