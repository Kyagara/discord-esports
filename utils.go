package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func respondWithError(interaction *discordgo.Interaction, err error) {
	client.logger.Error(fmt.Sprintf("Error executing command: %v", err))

	err = client.session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Error executing command.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		client.logger.Error(fmt.Sprintf("Error responding with error: %v", err))
	}
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

func getDefaultIntString(value string) string {
	if value != "" {
		return cases.Title(language.English, cases.NoLower).String(strings.ToLower(value))
	} else {
		return "0"
	}
}

func getDefaultString(value string) string {
	if value != "" {
		return cases.Title(language.English, cases.NoLower).String(strings.ToLower(value))
	} else {
		return "None"
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
