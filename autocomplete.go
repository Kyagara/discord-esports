package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func autoCompleteChampionName(session *discordgo.Session, interaction *discordgo.InteractionCreate, userInput string) {
	filteredNames := make(map[string]string)
	for id, name := range championsNames {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(userInput)) {
			filteredNames[name] = id
		}

		if len(filteredNames) == 20 {
			break
		}
	}

	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for id, name := range filteredNames {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{Name: name, Value: id})
	}

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})

	if err != nil {
		client.logger.Error(fmt.Sprintf("Error sending champion autocomplete: %v", err))
	}
}

func autoCompleteSpell(session *discordgo.Session, interaction *discordgo.InteractionCreate, championName string) {
	spells := spellsInfo[championName]

	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for _, spell := range spells {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{Name: spell.FullName, Value: fmt.Sprintf("%v,%v", spell.Key, spell.Index)})
	}

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})

	if err != nil {
		client.logger.Error(fmt.Sprintf("Error sending spell autocomplete: %v", err))
	}
}
