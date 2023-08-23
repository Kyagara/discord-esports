package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func InfoCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	fields := []*discordgo.MessageEmbedField{
		{Name: "Last update", Value: lastUpdate.Format(time.RFC822), Inline: true},
		{Name: "Last post", Value: lastPost.Format(time.RFC822), Inline: true},
		{Name: "Source code", Value: "https://github.com/Kyagara/discord-esports"}}

	embed := &discordgo.MessageEmbed{Title: "Info", Color: 0x9b311a, Fields: fields}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
