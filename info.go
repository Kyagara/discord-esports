package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func InfoCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	cmds := make([]string, len(commands))
	for i, cmd := range commands {
		cmds[i] = cmd.Name
	}

	nextUpdate := lastUpdate.Add(time.Duration(config.UpdateDateTimer) * time.Millisecond)
	nextPost := lastUpdate.Add(time.Duration(config.PostDataTimer) * time.Millisecond)

	embed := &discordgo.MessageEmbed{Title: "Info", Color: embedColor, Fields: []*discordgo.MessageEmbedField{
		{Name: "Last update", Value: fmt.Sprintf("<t:%v:R>.\nNext: <t:%v:R>.", lastUpdate.UnixMilli()/1000, nextUpdate.UnixMilli()/1000), Inline: true},
		{Name: "Last post", Value: fmt.Sprintf("<t:%v:R>.\nNext: <t:%v:R>.", lastPost.UnixMilli()/1000, nextPost.UnixMilli()/1000), Inline: true},
		{Name: "Started", Value: fmt.Sprintf("<t:%v:R>", started.UnixMilli()/1000), Inline: true},
		{Name: "Commands", Value: strings.Join(cmds, ", ")},
		{Name: "Source code", Value: "[Github](https://github.com/Kyagara/discord-esports)"}}}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
