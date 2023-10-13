package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commandButtonsID []string = []string{"modifiers", "notes", "spell", "skins", "spells"}
)

func readyEvent(session *discordgo.Session, ready *discordgo.Ready) {
	client.logger.Info(fmt.Sprintf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator))

	if *removeFlag || *registerFlag {
		return
	}

	err := updateAllData()
	if err != nil {
		client.logger.Error(fmt.Sprintf("Error updating data: %v", err))
	}

	err = postAllData()
	if err != nil {
		client.logger.Error(fmt.Sprintf("Error sending data: %v", err))
	}
}

func interactionsEvent(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if *removeFlag || *registerFlag {
		return
	}

	if interaction.Type == discordgo.InteractionMessageComponent {
		id := strings.Split(interaction.MessageComponentData().CustomID, "_")

		if !contains(commandButtonsID, id[0]) {
			respondWithError(interaction.Interaction, fmt.Errorf("message component id '%s' not found", id[0]))
			return
		}

		champion, ok := championsEmbeds[id[1]]
		if !ok {
			respondWithError(interaction.Interaction, fmt.Errorf("champion key '%s' not found", id[1]))
			return
		}

		switch id[0] {
		case "skins":
			respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{&champion.Skins})
			return
		case "spells":
			respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{&champion.Spells})
			return
		}

		// id[1] = champion key
		championSpells, ok := spellsEmbeds[id[1]]
		if !ok {
			respondWithError(interaction.Interaction, fmt.Errorf("champion key '%s' not found", id[1]))
			return
		}

		// id[2] = spell key
		spells, ok := championSpells[id[2]]
		if !ok {
			respondWithError(interaction.Interaction, fmt.Errorf("spell key '%s' not found", id[2]))
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
