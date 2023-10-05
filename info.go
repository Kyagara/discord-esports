package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func InfoCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	cmds := make([]string, len(commands))
	for i, cmd := range commands {
		cmds[i] = cmd.Name
	}

	embed := &discordgo.MessageEmbed{Title: "Info", Color: embedColor, Fields: []*discordgo.MessageEmbedField{
		{Name: "Last update", Value: fmt.Sprintf("<t:%v:R>.\nTimer: %v hours", lastUpdate.UnixMilli()/1000, config.UpdateDateTimer/1000/60/60), Inline: true},
		{Name: "Last post", Value: fmt.Sprintf("<t:%v:R>.\nTimer: %v hours.", lastPost.UnixMilli()/1000, config.PostDataTimer/1000/60/60), Inline: true},
		{Name: "Started", Value: fmt.Sprintf("<t:%v:R>", started.UnixMilli()/1000), Inline: true},
		{Name: "Commands", Value: strings.Join(cmds, ", ")},
		{Name: "Source code", Value: "[Github](https://github.com/Kyagara/discord-esports)"}}}

	respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{embed})
}
