package main

import (
	"fmt"
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

	var update bool
	if optionMap["update"] != nil {
		update = optionMap["update"].BoolValue()
		if update && hasPermissions(session, interaction) {
			now = time.Now()
			tomorrow = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

			switch game {
			case "LOL":
				err := updateLOLData()
				if err != nil {
					respondWithError(interaction.Interaction, fmt.Errorf("error updating LOL data: %v", err))
					return
				}
			case "VAL":
				err := updateVALData()
				if err != nil {
					respondWithError(interaction.Interaction, fmt.Errorf("error updating VAL data: %v", err))
					return
				}
			}

			respondWithMessage(interaction.Interaction, "Updated data without any errors.")
			return
		}
	}

	switch game {
	case "LOL":
		respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{createLOLMessageEmbed()})
	case "VAL":
		respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{createVALMessageEmbed()})
	}
}
