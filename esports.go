package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func EsportsCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	game := optionMap["game"].StringValue()

	if optionMap["update"] != nil {
		if !hasPermissions(session, interaction) {
			return
		}

		update := optionMap["update"].BoolValue()
		if update {
			now = time.Now()
			tomorrow = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

			var err error
			switch game {
			case "lol":
				err = updateLOLData()
			case "val":
				err = updateVALData()
			default:
				respondWithMessage(interaction.Interaction, "Game not found.")
				return
			}

			if err != nil {
				client.logger.Error(fmt.Sprintf("Error updating %v data: %v", strings.ToUpper(game), err))
				respondWithMessage(interaction.Interaction, fmt.Sprintf("An error occured updating %v data.", strings.ToUpper(game)))
				return
			}

			respondWithMessage(interaction.Interaction, "Updated data without any errors.")
			return
		}
	}

	switch game {
	case "lol":
		respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{createLOLMessageEmbed()})
	case "val":
		respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{createVALMessageEmbed()})
	default:
		respondWithMessage(interaction.Interaction, "Game not found.")
	}
}
