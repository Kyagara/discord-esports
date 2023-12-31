package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func readyEvent(session *discordgo.Session, ready *discordgo.Ready) {
	client.logger.Info(fmt.Sprintf("Logged in as: %v", session.State.User.Username))

	if *registerFlag {
		return
	}
}

func interactionsEvent(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if *registerFlag {
		return
	}

	if interaction.Type == discordgo.InteractionMessageComponent {
		id := strings.Split(interaction.MessageComponentData().CustomID, "_")

		if !contains(commandButtonsID, id[0]) {
			msg := fmt.Sprintf("Message component id '%s' not found.", id[0])
			client.logger.Error(msg)
			respondWithMessage(interaction.Interaction, msg)
			return
		}

		// id[1] = champion key
		championSpells, ok := spellsEmbeds[id[1]]
		if !ok {
			msg := fmt.Sprintf("Champion key '%s' not found.", id[1])
			client.logger.Error(msg)
			respondWithMessage(interaction.Interaction, msg)
			return
		}

		// id[2] = spell key
		spells, ok := championSpells[id[2]]
		if !ok {
			msg := fmt.Sprintf("Spell key '%s' not found.", id[2])
			client.logger.Error(msg)
			respondWithMessage(interaction.Interaction, msg)
			return
		}

		// id[3] = spell index
		spellIndex, err := strconv.Atoi(id[3])
		if err != nil {
			client.logger.Error(fmt.Sprintf("error converting spell index to int: %v", err))
		}

		switch id[0] {
		case "modifiers":
			respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{&spells[spellIndex].Modifiers})
		case "notes":
			respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{&spells[spellIndex].Notes})
		case "spell":
			respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{&spells[spellIndex].General})
		}

		return
	}

	if command, ok := client.commands[interaction.ApplicationCommandData().Name]; ok {
		command.Handler(session, interaction)
	}
}
