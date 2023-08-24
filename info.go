package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func InfoCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var cmds []string
	for _, cmd := range commands {
		cmds = append(cmds, cmd.Name)
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "Last update", Value: fmt.Sprintf("<t:%v:R>", lastUpdate.UnixMilli()/1000), Inline: true},
		{Name: "Last post", Value: fmt.Sprintf("<t:%v:R>", lastPost.UnixMilli()/1000), Inline: true},
		{Name: "Commands", Value: strings.Join(cmds, ", ")},
		{Name: "Source code", Value: "https://github.com/Kyagara/discord-esports"}}

	embed := &discordgo.MessageEmbed{Title: "Info", Color: embedColor, Fields: fields}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
