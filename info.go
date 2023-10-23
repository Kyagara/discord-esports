package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func InfoCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var cmds []string
	for _, c := range client.commands {
		cmds = append(cmds, c.Interaction.Name)
	}

	embed := &discordgo.MessageEmbed{
		Title:     "Info",
		Color:     DISCORD_EMBED_COLOR,
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: session.State.User.AvatarURL("256")},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Commands", Value: fmt.Sprintf("%v.", strings.Join(cmds, ", "))},
			{Name: "Started", Value: fmt.Sprintf("<t:%v:R>.", started.UnixMilli()/1000)},
			{Name: "Source code", Value: "[Github](https://github.com/Kyagara/discord-esports)"},
		}}

	respondWithEmbed(interaction.Interaction, []*discordgo.MessageEmbed{embed})
}
