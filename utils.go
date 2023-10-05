package main

import (
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func newRequest(endpoint string) (*http.Request, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "discord-esports-bot - https://github.com/Kyagara/discord-esports")
	return req, nil
}

func hasPermissions(session *discordgo.Session, interaction *discordgo.InteractionCreate) bool {
	if len(config.ModRoles) == 0 {
		return true
	}

	for _, memberRole := range interaction.Member.Roles {
		for _, modRole := range config.ModRoles {
			if memberRole == modRole {
				return true
			}
		}
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You don't have permission to run this command.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return false
}

func respondWithErrorEmbed(interaction *discordgo.Interaction, err error) {
	log.Printf("Error executing command: %v", err)

	session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Error executing command.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
