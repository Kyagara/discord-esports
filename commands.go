package main

import (
	"github.com/bwmarrin/discordgo"
)

func interactionsEvent(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if *removeFlag || *registerFlag {
		return
	}

	if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
		handler(session, interaction)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "lol",
			Description:              "Send a list of upcoming League of Legends games.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia uma lista de próximos jogos de League of Legends."},
		},
		{
			Name:                     "val",
			Description:              "Send a list of upcoming Valorant games.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia uma lista de próximos jogos de Valorant."},
		},
		{
			Name:                     "update",
			Description:              "Force all data to update.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Força todos os dados a serem atualizados."},
		},
		{
			Name:                     "info",
			Description:              "Send information about the bot.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia informação sobre o bot."},
		},
		{
			Name:                     "post",
			Description:              "Force all data to be sent again.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Força todos os dados a serem enviados."},
		},
		{
			Name:                     "champion",
			Description:              "Get League of Legends champion stats.",
			DescriptionLocalizations: &map[discordgo.Locale]string{discordgo.PortugueseBR: "Envia informações de um champion de League of Legends."},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:                     "champion",
					Required:                 true,
					Autocomplete:             true,
					Type:                     discordgo.ApplicationCommandOptionString,
					Description:              "Champion name.",
					DescriptionLocalizations: map[discordgo.Locale]string{discordgo.PortugueseBR: "Nome do champion."},
				},
			},
		},
	}

	commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		"lol":      LOLEsportsCommand,
		"val":      VALEsportsCommand,
		"update":   UpdateCommand,
		"info":     InfoCommand,
		"post":     PostCommand,
		"champion": ChampionCommand,
	}
)
