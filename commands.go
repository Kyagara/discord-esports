package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func interactionsEvent(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
		handler(session, interaction)
	}
}

func registerCommands() error {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return fmt.Errorf("invalid bot parameters: %v", err)
	}

	err = session.Open()
	if err != nil {
		return fmt.Errorf("error opening session: %v", err)
	}

	defer session.Close()

	log.Print("Registering commands.")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, cmd := range commands {
		command, err := session.ApplicationCommandCreate(session.State.User.ID, config.GuildID, cmd)
		if err != nil {
			return fmt.Errorf("error registering '%v' command: %v", cmd.Name, err)
		}

		registeredCommands[i] = command
		log.Printf("Registered '%v' command.", command.Name)
	}

	log.Print("Commands registered.")
	return nil
}

func removeCommands() error {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return fmt.Errorf("invalid bot parameters: %v", err)
	}

	err = session.Open()
	if err != nil {
		return fmt.Errorf("error opening session: %v", err)
	}

	defer session.Close()

	log.Print("Removing commands")

	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, config.GuildID)
	if err != nil {
		return fmt.Errorf("error fetching registered commands: %v", err)
	}

	for _, cmd := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, config.GuildID, cmd.ID)
		if err != nil {
			return fmt.Errorf("error deleting '%v' command: %v", cmd.Name, err)
		}

		log.Printf("Deleted '%v' command.", cmd.Name)
	}

	log.Print("Commands removed.")
	return nil
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
	}

	commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		"lol":    LOLEsportsCommand,
		"val":    VALEsportsCommand,
		"update": UpdateCommand,
		"info":   InfoCommand,
		"post":   PostCommand,
	}
)
