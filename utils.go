package main

import (
	"fmt"
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
	if len(client.config.ModRoles) == 0 {
		return true
	}

	for _, memberRole := range interaction.Member.Roles {
		for _, modRole := range client.config.ModRoles {
			if memberRole == modRole {
				return true
			}
		}
	}

	respondWithMessage(interaction.Interaction, "You don't have permission to run this command.")
	return false
}

func respondWithEmbed(interaction *discordgo.Interaction, embed []*discordgo.MessageEmbed) {
	err := client.session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embed,
		},
	})

	if err != nil {
		client.logger.Error(fmt.Sprintf("Error responding with embed: %v", err))
	}
}

func respondWithMessage(interaction *discordgo.Interaction, message string) {
	err := client.session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		client.logger.Error(fmt.Sprintf("Error responding with message: %v", err))
	}
}

func contains(ids []string, str string) bool {
	for _, id := range ids {
		if id == str {
			return true
		}
	}

	return false
}
